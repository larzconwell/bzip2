package bzip2

import (
	"compress/flate"
	"fmt"
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

// These constants are copied from the flate package, so that
// code does not also have to import "compress/flate".
const (
	BestSpeed       = flate.BestSpeed
	BestCompression = flate.BestCompression
)

// Writer is an io.WriteCloser. Writes to a Writer are
// compressed and written to an underlying io.Writer.
type Writer struct {
	bw          *bits.Writer
	block       *block
	crc         uint32
	wroteHeader bool
	closed      bool
	err         error
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
// The levels range from 1 (BestSpeed) to 9 (BestCompression);
// higher levels typically run slower but compress more.
//
// If level is in the range [1, 9] then the error returned will
// be nil. Otherwise the error returned will be non-nil.
func NewWriterLevel(w io.Writer, level int) (*Writer, error) {
	if level < BestSpeed || level > BestCompression {
		return nil, fmt.Errorf("bzip2: invalid compression level: %d", level)
	}

	return &Writer{
		bw:    bits.NewWriter(w),
		block: newBlock(level * baseBlockSize),
	}, nil
}

// Write writes a compressed form of p to the underlying
// io.Writer. The compressed bytes are not necessarily
// flushed until the Writer is closed.
func (w *Writer) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	var n int

	// Handle writing the file header.
	if !w.wroteHeader {
		w.err = w.writeHeader()
		if w.err != nil {
			return 0, w.err
		}

		w.wroteHeader = true
	}

	n, w.err = w.write(p)
	return n, w.err
}

// writeHeader writes the file header.
func (w *Writer) writeHeader() error {
	w.bw.WriteBits(16, fileMagic)
	w.bw.WriteBits(8, 'h')
	w.bw.WriteBits(8, uint64('0'+w.block.size/baseBlockSize))

	return w.bw.Err()
}

// write handles the writing of block data and writing
// completed blocks to the underlying io.Writer.
func (w *Writer) write(p []byte) (int, error) {
	n, err := w.block.Write(p)
	if err != errBlockSizeReached {
		return n, err
	}

	// Write the complete block, left over
	// bytes being written to a new block.
	err = w.writeBlock()
	if err != nil {
		return n, err
	}

	if n != len(p) {
		var nn int
		nn, err = w.write(p[n:])
		n += nn
	}

	return n, err
}

// writeBlock writes the current block to the
// underlying io.Writer and updates the files crc.
func (w *Writer) writeBlock() error {
	err := w.block.WriteBlock(w.bw)
	if err != nil {
		return err
	}

	w.crc = ((w.crc << 1) | (w.crc >> 31)) ^ w.block.crc
	w.block = newBlock(w.block.size)
	return nil
}

// Flush flushes any pending compressed data
// to the underlying writer.
func (w *Writer) Flush() error {
	if w.err != nil {
		return w.err
	}
	if w.closed || w.block.Len() == 0 {
		return nil
	}

	// Handle writing the file header.
	if !w.wroteHeader {
		w.err = w.writeHeader()
		if w.err != nil {
			return w.err
		}

		w.wroteHeader = true
	}

	w.err = w.writeBlock()
	return w.err
}

// Reset discards the state of Writer and makes it equivalent
// to the result of NewWriter or NewWriterLevel, but writing
// to dst instead.
func (w *Writer) Reset(dst io.Writer) {
	w.bw = bits.NewWriter(dst)
	w.block = newBlock(w.block.size)
	w.crc = 0
	w.wroteHeader = false
	w.closed = false
	w.err = nil
}

// Close closes the Writer, flushing any unwritten data to the
// underlying io.Writer, but does not close the underlying io.Writer.
func (w *Writer) Close() error {
	if w.err != nil {
		return w.err
	}
	if w.closed {
		return nil
	}
	w.closed = true

	// Handle writing the file header.
	if !w.wroteHeader {
		w.err = w.writeHeader()
		if w.err != nil {
			return w.err
		}

		w.wroteHeader = true
	}

	// Flush the current block.
	if w.block.Len() != 0 {
		w.err = w.writeBlock()
		if w.err != nil {
			return w.err
		}
	}

	w.bw.WriteBits(48, finalMagic)
	w.bw.WriteBits(32, uint64(w.crc))
	bufferedBits := w.bw.Buffered()
	if bufferedBits != 0 {
		w.bw.WriteBits(8-bufferedBits, 0)
	}

	w.err = w.bw.Err()
	return w.err
}
