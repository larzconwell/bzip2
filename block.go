package bzip2

import (
	"bytes"
	"errors"
	"hash/crc32"
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
func (b *block) Len() int {
	return b.buf.Len()
}

// Write writes p to the blocks buffer. If writing p exceeds the blocks size
// only the bytes that can fit will be written and errBlockSizeReached is
// returned.
func (b *block) Write(p []byte) (int, error) {
	exceedsSize := false

	// Do the initial RLE step on demand since RLE can actually make p grow.
	// This ensures that the block size doesn't end up more than b.size
	// after RLE.
	data := rlEncode(p)
	if b.buf.Len()+len(data) > b.size {
		exceedsSize = true
		data = data[:b.size-b.buf.Len()]
	}

	n, err := b.buf.Write(data)
	p = p[:rlIndexOf(data, n-1)+1]
	if err == nil {
		b.crc = crc32.Update(b.crc, crc32.IEEETable, p)

		if exceedsSize || b.buf.Len() == b.size {
			err = errBlockSizeReached
		}
	}

	return len(p), err
}

// WriteBlock compresses the contented buffered and writes a block to the bit
// writer given. The number of bits written is returned.
func (b *block) WriteBlock(bw *bitWriter) (int, error) {
	data := b.buf.Bytes()

	// Get the symbols used in data. Int is used here to simplify code
	// to generate the bitmap.
	symbols := make([]int, 256)
	for _, b := range data {
		symbols[int(b)] = 1
	}

	// BWT step.
	bwt := make([]byte, len(data))
	bwtidx := bwTransform(bwt, data)
	data = bwt

	// MTF step.
	mtfTransform(symbols, data, data)

	// Write the block header.
	bw.WriteBits(48, blockMagic)
	bw.WriteBits(32, uint64(b.crc))
	bw.WriteBits(1, 0)
	bitsWrote := 81

	// Write the BWT index.
	bw.WriteBits(24, uint64(bwtidx))
	bitsWrote += 24

	// Write the sparse bit array for used symbols.
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
	bw.WriteBits(16, uint64(symbolRangeUsedBitmap))
	bitsWrote += 16
	for _, symRange := range symbolRanges {
		if symRange > 0 {
			bw.WriteBits(16, uint64(symRange))
			bitsWrote += 16
		}
	}

	return bitsWrote, bw.Err()
}
