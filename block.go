package bzip2

import (
	"bytes"
	"errors"
	"hash/crc32"
	"math"
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
	exceedsSize := false

	// Do the initial RLE step on demand since RLE can actually make p grow.
	// This ensures that the buffer doesn't end up more than b.size
	// afterwards.
	encoded := rlEncode(p)
	if b.buf.Len()+len(encoded) > b.size {
		exceedsSize = true
		encoded = encoded[:b.size-b.buf.Len()]
	}

	n, err := b.buf.Write(encoded)
	p = p[:rlIndexOf(n-1, encoded)+1]
	if err == nil {
		b.crc = crc32.Update(b.crc, crc32.IEEETable, p)

		if exceedsSize || b.buf.Len() == b.size {
			err = errBlockSizeReached
		}
	}

	return len(p), err
}

// WriteBlock compresses the content buffered and writes a block to the bit
// writer given. The number of bits written is returned.
func (b *block) WriteBlock(bw *bitWriter) (int, error) {
	data := b.buf.Bytes()
	symbols, reducedSymbols := symbolSet(data)

	// Generate the bitmaps for the used symbols.
	symbolRangeUsedBitmap := 0
	symbolRanges := make([]int, 16)
	for i, symRange := range symbolRanges {
		// Toggle the bits for the 16 symbols in the range.
		for j := 0; j < 16; j++ {
			symRange = (symRange << 1) + symbols[16*i+j]
		}
		symbolRanges[i] = symRange

		// Toggle the bit for the range in the bitmap.
		rangePresent := 0
		if symRange > 0 {
			rangePresent = 1
		}
		symbolRangeUsedBitmap = (symbolRangeUsedBitmap << 1) + rangePresent
	}

	// BWT step.
	bwt := make([]byte, len(data))
	bwtidx := bwTransform(bwt, data)

	// MTF step.
	mtf := bwt
	mtfTransform(reducedSymbols, mtf, bwt)

	// RLE2 step.
	rle := rl2Encode(reducedSymbols, mtf)
	freq := symbolFrequencies(reducedSymbols, rle)

	// Setup the huffman trees required to encode rle.
	huffmanTrees, treeIndexes := generateHuffmanTrees(freq, rle)

	// Get the MTF encoded huffman tree indexes.
	treeIndexesSymbols := make([]byte, len(huffmanTrees))
	for i := range huffmanTrees {
		treeIndexesSymbols[i] = byte(i)
	}
	treeIndexesBytes := make([]byte, len(treeIndexes))
	for i, idx := range treeIndexes {
		treeIndexesBytes[i] = byte(idx)
	}
	mtfTransform(treeIndexesSymbols, treeIndexesBytes, treeIndexesBytes)

	// Write the block header.
	bw.WriteBits(48, blockMagic)
	bw.WriteBits(32, uint64(b.crc))
	bw.WriteBits(1, 0)
	bitsWrote := 81

	// Write the BWT index.
	bw.WriteBits(24, uint64(bwtidx))
	bitsWrote += 24

	// Write the symbol bitmaps.
	bw.WriteBits(16, uint64(symbolRangeUsedBitmap))
	bitsWrote += 16
	for _, symRange := range symbolRanges {
		if symRange > 0 {
			bw.WriteBits(16, uint64(symRange))
			bitsWrote += 16
		}
	}

	// Write the huffman tree numbers.
	bw.WriteBits(3, uint64(len(huffmanTrees)))
	bw.WriteBits(15, uint64(len(treeIndexes)))
	bitsWrote += 18

	// Write the huffman tree indexes in unary encoding.
	for _, idx := range treeIndexesBytes {
		for i := byte(0); i < idx; i++ {
			bw.WriteBits(1, 1)
		}

		bw.WriteBits(1, 0)
		bitsWrote += int(idx) + 1
	}

	// Write the delta encoded code-lengths for the huffman trees codes.
	for _, tree := range huffmanTrees {
		// Get the smallest code-length in the huffman tree.
		length := 0
		for i, code := range tree.Codes {
			if i == 0 || code.Len() < length {
				length = code.Len()
			}
		}
		bw.WriteBits(5, uint64(length))
		bitsWrote += 5

		// Write the code-lengths as modifications to the base length.
		for _, code := range tree.Codes {
			delta := int(math.Abs(float64(length - code.Len())))

			// 2 is increment, 3 is decrement.
			op := uint64(2)
			if length > code.Len() {
				op = 3
			}
			length = code.Len()

			for i := 0; i < delta; i++ {
				bw.WriteBits(2, op)
			}

			bw.WriteBits(1, 0)
			bitsWrote += delta + 1
		}
	}

	// Write the encoded contents, using the huffman trees generated
	// switching them out every 50 symbols.
	decoded := 0
	treeIndex := 0
	huffmanTree := huffmanTrees[treeIndexes[treeIndex]]
	for _, b := range rle {
		if decoded == 50 {
			decoded = 0
			treeIndex++
			huffmanTree = huffmanTrees[treeIndexes[treeIndex]]
		}
		code := huffmanTree.Codes[b]

		bw.WriteBits(uint(code.Len()), code.Bits)
		bitsWrote += code.Len()
		decoded++
	}

	return bitsWrote, bw.Err()
}

// symbolSet gets the symbol set for a slice of bytes. Int is used instead
// of bool to simplify the symbol bitmap code. The reduced list of symbols
// used is also returned.
func symbolSet(data []byte) ([256]int, []byte) {
	var symbols [256]int
	reduced := make([]byte, 0, 256)

	for _, b := range data {
		symbols[b] = 1
	}

	for i, present := range symbols {
		if present > 0 {
			reduced = append(reduced, byte(i))
		}
	}

	return symbols, reduced
}
