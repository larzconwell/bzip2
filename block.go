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
	// Get the bytes we can write in this block.
	exceedsSize := false
	if b.buf.Len()+len(p) > b.size {
		exceedsSize = true
		p = p[:b.size-b.buf.Len()]
	}

	n, err := b.buf.Write(p)
	if err == nil {
		b.crc = crc32.Update(b.crc, crc32.IEEETable, p)

		if exceedsSize || b.buf.Len() == b.size {
			err = errBlockSizeReached
		}
	}

	return n, err
}

// WriteBlock compresses the contented buffered and writes a block to the bit
// writer given. The number of bits written is returned.
func (b *block) WriteBlock(bw *bitWriter) (int, error) {
	bw.WriteBits(48, blockMagic)
	bw.WriteBits(32, uint64(b.crc))
	bw.WriteBits(1, 0)
	bitsWrote := 81

	// Initial RLE step.
	data := encodeRLE(b.buf.Bytes())
	// BWT step.
	ptr := bwTransform(data, data)
	bw.WriteBits(24, uint64(ptr))

	return bitsWrote, bw.Err()
}
