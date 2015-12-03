package bzip2

import (
	"errors"
	"io"

	"github.com/larzconwell/bzip2/internal/bits"
)

const (
	// fileMagic is the bzip2 magic number, BZ.
	fileMagic = 0x425a
	// finalMagic signifies the end of the block data.
	finalMagic = 0x177245385090
	// baseBlockSize is the base for block sizes.
	baseBlockSize = 100000
)

var (
	// ErrWriteAfterClose occurs when a Write
	// occurs after a Writer is closed.
	ErrWriteAfterClose = errors.New("bzip2: write after close")
)

// Writer is an io.WriteCloser. Writes to a Writer are
// compressed and written to the underlying io.Writer.
type Writer struct {
	bw          *bits.Writer
	block       *block
	crc         uint32
	wroteHeader bool
	closed      bool
}

// NewWriter returns a new Writer. Writes to the returned
// writer are compressed and written to w.
//
// It is the caller's responsibility to call Close on the
// Writer when done. Writes may be buffered and not
// flushed until Close.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		bw:    bits.NewWriter(w),
		block: newBlock(6 * baseBlockSize),
	}
}

// NewWriterLevel is like NewWriter but specifies the
// compression level.
//
// The level must be between 1 and 9, any other values are
// set to the closest valid level.
func NewWriterLevel(w io.Writer, level int) *Writer {
	if level < 1 {
		level = 1
	} else if level > 9 {
		level = 9
	}

	return &Writer{
		bw:    bits.NewWriter(w),
		block: newBlock(level * baseBlockSize),
	}
}

// Write writes a compressed form of p to the underlying
// io.Writer. The compressed bytes are not necessarily
// flushed until the Writer is closed.
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

// write handles the writing of block data and writing
// completed blocks to the underlying io.Writer.
func (b *Writer) write(p []byte) (int, error) {
	n, err := b.block.Write(p)
	if err != errBlockSizeReached {
		return n, err
	}

	// Write the complete block, left over
	// bytes being written to a new block.
	err = b.writeBlock()
	if err != nil {
		return n, err
	}

	if n != len(p) {
		var nn int
		nn, err = b.write(p[n:])
		n += nn
	}

	return n, err
}

// writeBlock writes the current block to the
// underlying io.Writer and updates the files crc.
func (b *Writer) writeBlock() error {
	err := b.block.WriteBlock(b.bw)
	if err != nil {
		return err
	}

	b.crc = ((b.crc << 1) | (b.crc >> 31)) ^ b.block.crc
	b.block = newBlock(b.block.size)
	return nil
}

// Flush flushes any pending compressed data
// to the underlying io.Writer.
func (b *Writer) Flush() error {
	if b.closed || b.block.Len() == 0 {
		return nil
	}

	return b.writeBlock()
}

// Reset discards the Writers state and makes it equivalent to the result
// of its original state from NewWriter or NewWriterLevel, but writing to
// w instead. This permits reusing a Writer rather than allocating a new one.
func (b *Writer) Reset(w io.Writer) {
	b.bw = bits.NewWriter(w)
	b.block = newBlock(b.block.size)
	b.crc = 0
	b.wroteHeader = false
	b.closed = false
}

// Close closes the Writer, flushing any unwritten data to the
// underlying io.Writer, but does not close the underlying io.Writer.
func (b *Writer) Close() error {
	if b.closed {
		return nil
	}
	b.closed = true

	// Flush the current block.
	if b.block.Len() != 0 {
		err := b.writeBlock()
		if err != nil {
			return err
		}
	}

	b.bw.WriteBits(48, finalMagic)
	b.bw.WriteBits(32, uint64(b.crc))
	bufferedBits := b.bw.Buffered()
	if bufferedBits != 0 {
		b.bw.WriteBits(8-bufferedBits, 0)
	}

	return b.bw.Err()
}
