package bzip2

import (
	"errors"
	"io"
)

const (
	fileMagic     = 0x425a         // "BZ"
	finalMagic    = 0x177245385090 // sqrt(pi)
	baseBlockSize = 100000         // Base bytes for the block size.
)

var (
	// ErrWriteAfterClose occurs when a write occurs after a Writer is closed.
	ErrWriteAfterClose = errors.New("bzip2: write after close")
)

// Writer is an io.WriteCloser. Writes to a Writer are compressed and written
// to the underlying writer.
type Writer struct {
	bw          *bitWriter
	block       *block
	crc         uint32
	totalWrote  int // In bits.
	wroteHeader bool
	closed      bool
}

// NewWriter returns a new Writer. Writes to the returned writer are compressed
// and written to w.
//
// It is the caller's responsibility to call Close on the WriteCloser when done.
// Writes may be buffered and not flushed until Close.
func NewWriter(w io.Writer) *Writer {
	return &Writer{bw: newBitWriter(w), block: newBlock(6 * baseBlockSize)}
}

// NewWriterLevel is like NewWriter except a specific blockSize is given to
// control the level of compression.
//
// The blockSize must be between 1 and 9, any other values are set to the
// closest valid blockSize.
func NewWriterLevel(w io.Writer, blockSize int) *Writer {
	if blockSize < 1 {
		blockSize = 1
	} else if blockSize > 9 {
		blockSize = 9
	}

	return &Writer{bw: newBitWriter(w), block: newBlock(blockSize * baseBlockSize)}
}

// Write writes a compressed form of p to the underlying writer. The writes may
// be buffered until a Close or Flush.
func (b *Writer) Write(p []byte) (int, error) {
	if b.closed {
		return 0, ErrWriteAfterClose
	}

	if !b.wroteHeader {
		err := b.writeHeader()
		if err != nil {
			return 0, err
		}

		b.wroteHeader = true
	}

	return b.write(p)
}

// writeHeader writes the file header.
func (b *Writer) writeHeader() error {
	b.bw.WriteBits(16, fileMagic)
	b.bw.WriteBits(8, 'h')
	b.bw.WriteBits(8, uint64('0'+b.block.size/baseBlockSize))

	return b.bw.Err()
}

// write handles the writing of block data and writing completed blocks to
// underlying writer.
func (b *Writer) write(p []byte) (int, error) {
	n, err := b.block.Write(p)
	if err == errBlockSizeReached {
		err := b.writeBlock()
		if err != nil {
			return n, err
		}

		if n != len(p) {
			var nn int
			nn, err = b.write(p[n:])
			n += nn
		}
	}

	return n, err
}

// writeBlock writes the current block to the underlying writer, and updates
// the files crc and total bits wrote.
func (b *Writer) writeBlock() error {
	n, err := b.block.WriteBlock(b.bw)
	if err != nil {
		return err
	}

	b.crc = ((b.crc << 1) | (b.crc >> 31)) ^ b.block.crc
	b.totalWrote += n
	b.block = newBlock(b.block.size)
	return nil
}

// Flush flushes any pending compressed data to the underlying writer.
func (b *Writer) Flush() error {
	if b.block.Len() == 0 {
		return nil
	}

	return b.writeBlock()
}

// Reset resets the Writers state and makes it equivalent to the result of
// NewWriter or NewWriterLevel, writing to w.
func (b *Writer) Reset(w io.Writer) {
	b.bw = newBitWriter(w)
	b.block = newBlock(b.block.size)
	b.crc = 0
	b.totalWrote = 0
	b.wroteHeader = false
	b.closed = false
}

// Close closes the Writer, flushing any unwritten data to the underlying writer.
func (b *Writer) Close() error {
	if b.closed {
		return nil
	}
	b.closed = true

	err := b.Flush()
	if err != nil {
		return err
	}

	b.bw.WriteBits(48, finalMagic)
	b.bw.WriteBits(32, uint64(b.crc))
	padding := b.totalWrote % 8
	if padding != 0 {
		b.bw.WriteBits(uint(padding), 0)
	}

	return b.bw.Err()
}
