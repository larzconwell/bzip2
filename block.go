package bzip2

import (
	"bytes"
	"errors"
	"math"

	"github.com/larzconwell/bzip2/internal/bits"
	"github.com/larzconwell/bzip2/internal/bwt"
	"github.com/larzconwell/bzip2/internal/crc32"
	"github.com/larzconwell/bzip2/internal/huffman"
	"github.com/larzconwell/bzip2/internal/mtf"
	"github.com/larzconwell/bzip2/internal/rle"
	"github.com/larzconwell/bzip2/internal/rle2"
	"github.com/larzconwell/bzip2/internal/symbols"
)

const (
	blockMagic = 0x314159265359 // pi
)

var (
	errBlockSizeReached = errors.New("bzip2: Block size reached")
)

// block handles the compression of data for a single block of a given size.
type block struct {
	buf  *bytes.Buffer
	size int
	crc  uint32
}

// newBlock creates a compression block for data up to the given size.
func newBlock(size int) *block {
	return &block{buf: bytes.NewBuffer(make([]byte, 0, size)), size: size}
}

// Len returns the number of bytes written to the block.
func (b block) Len() int {
	return b.buf.Len()
}

// Write writes p to the blocks buffer. If writing p exceeds the blocks size
// only the bytes that can fit will be written and errBlockSizeReached is
// returned.
func (b *block) Write(p []byte) (int, error) {
	blockSizeReached := false

	// Do the initial RLE step on demand since RLE can actually make p grow.
	// This ensures that the buffer doesn't end up more than b.size
	// afterwards.
	encoded := rle.Encode(p)
	if b.buf.Len()+len(encoded) > b.size {
		encoded = encoded[:b.size-b.buf.Len()]

		// If the last 4 bytes are a repeat then a repeat
		// byte was stripped off leaving invalid RLE data.
		// To combat this just strip an extra character
		// off letting it get written to another block.
		if len(encoded) >= 4 {
			idx := len(encoded) - 4
			b := encoded[idx]
			data := []byte{b, b, b, b}

			if string(encoded[idx:]) == string(data) {
				encoded = encoded[:idx+3]
				blockSizeReached = true
			}
		}
	}

	n, err := b.buf.Write(encoded)
	p = p[:rle.IndexOf(n-1, encoded)+1]
	if err == nil {
		b.crc = crc32.Update(b.crc, p)

		if b.buf.Len() == b.size || (n == len(encoded) && blockSizeReached) {
			err = errBlockSizeReached
		}
	}

	return len(p), err
}

// WriteBlock compresses the content buffered and writes a block to the bit
// writer given.
func (b *block) WriteBlock(bw *bits.Writer) error {
	data := b.buf.Bytes()
	syms, reducedSyms := symbols.Get(data)

	// BWT step.
	bwtData := make([]byte, len(data))
	bwtidx := bwt.Transform(bwtData, data)

	// MTF step.
	mtfData := bwtData
	mtf.Transform(reducedSyms, mtfData, bwtData)

	// RLE2 step.
	rleData := rle2.Encode(reducedSyms, mtfData)
	freqs := rle2.GetFrequencies(reducedSyms, rleData)

	// Setup the huffman trees required to encode rle.
	trees, selections := huffman.GenerateTrees(freqs, rleData)

	// Get the MTF encoded huffman tree selections.
	treeSelectionSymbols := make(symbols.ReducedSet, len(trees))
	for i := range trees {
		treeSelectionSymbols[i] = byte(i)
	}
	treeSelectionBytes := make([]byte, len(selections))
	for i, selection := range selections {
		treeSelectionBytes[i] = byte(selection)
	}
	mtf.Transform(treeSelectionSymbols, treeSelectionBytes, treeSelectionBytes)

	// Write the block header.
	bw.WriteBits(48, blockMagic)
	bw.WriteBits(32, uint64(b.crc))
	bw.WriteBits(1, 0)

	// Write the contents that build the decoding steps.
	bw.WriteBits(24, uint64(bwtidx))
	b.writeSymbolBitmaps(bw, syms)
	bw.WriteBits(3, uint64(len(trees)))
	bw.WriteBits(15, uint64(len(selections)))
	b.writeTreeSelections(bw, treeSelectionBytes)
	b.writeTreeCodes(bw, trees)

	// Write the encoded contents, using the huffman trees generated
	// switching them out every 50 symbols.
	encoded := 0
	idx := 0
	tree := trees[selections[idx]]
	for _, b := range rleData {
		if encoded == huffman.TreeSelectionLimit {
			encoded = 0
			idx++
			tree = trees[selections[idx]]
		}
		code := tree.Codes[b]

		bw.WriteBits(uint(code.Len), code.Bits)
		encoded++
	}

	return bw.Err()
}

// writeSymbolBitmaps writes the bitmaps for the used symbols.
func (b *block) writeSymbolBitmaps(bw *bits.Writer, syms symbols.Set) {
	rangesUsed := 0
	ranges := make([]int, 16)

	for i, r := range ranges {
		// Toggle the bits for the 16 symbols in the range.
		for j := 0; j < 16; j++ {
			r = (r << 1) | syms[16*i+j]
		}
		ranges[i] = r

		// Toggle the bit for the range in the bitmap.
		present := 0
		if r > 0 {
			present = 1
		}
		rangesUsed = (rangesUsed << 1) | present
	}

	bw.WriteBits(16, uint64(rangesUsed))
	for _, r := range ranges {
		if r > 0 {
			bw.WriteBits(16, uint64(r))
		}
	}
}

// writeTreeSelections writes the huffman tree selections in unary encoding.
func (b *block) writeTreeSelections(bw *bits.Writer, selections []byte) {
	for _, selection := range selections {
		for i := byte(0); i < selection; i++ {
			bw.WriteBits(1, 1)
		}

		bw.WriteBits(1, 0)
	}
}

// writeTreeCodes writes the delta encoded code-lengths for
// the huffman trees codes.
func (b *block) writeTreeCodes(bw *bits.Writer, trees []*huffman.Tree) {
	for _, tree := range trees {
		// Get the smallest code-length in the huffman tree.
		codelen := 0
		for i, code := range tree.Codes {
			if i == 0 || code.Len < codelen {
				codelen = code.Len
			}
		}
		bw.WriteBits(5, uint64(codelen))

		// Write the code-lengths as modifications to the current length.
		for _, code := range tree.Codes {
			delta := int(math.Abs(float64(codelen - code.Len)))

			// 2 is increment, 3 is decrement.
			op := uint64(2)
			if codelen > code.Len {
				op = 3
			}
			codelen = code.Len

			for i := 0; i < delta; i++ {
				bw.WriteBits(2, op)
			}

			bw.WriteBits(1, 0)
		}
	}
}
