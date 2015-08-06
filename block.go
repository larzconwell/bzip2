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
	bw.WriteBits(48, blockMagic)
	bw.WriteBits(32, uint64(b.crc))
	bw.WriteBits(1, 0)
	bitsWrote := 81

	// BWT step.
	data := b.buf.Bytes()
	bwt := make([]byte, len(data))
	ptr := bwTransform(bwt, data)
	bw.WriteBits(24, uint64(ptr))

	return bitsWrote, bw.Err()
}
