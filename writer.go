package bzip2

import (
	"compress/flate"
	"io"

	"github.com/larzconwell/buffbits"
)

// These constants are copied from the flate package, so that
// code does not also have to import "compress/flate".
const (
	BestSpeed          = flate.BestSpeed
	BestCompression    = flate.BestCompression
	DefaultCompression = flate.DefaultCompression
)

const (
	defaultLevel  = 6
	baseBlockSize = 100_000
)

// Writer is an io.WriteCloser that compresses and writes data to an underlying io.Writer.
// If an error occurs while writing to a Writer, no more data will be written and all
// subsequent calls with return an error. The caller should call Close when done to
// flush any pending data.
type Writer struct {
	bw          *buffbits.Writer
	checksum    uint32
	blockSize   int
	wroteHeader bool
	err         error
}

// NewWriter creates a Writer that compresses and writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		bw:        buffbits.NewWriter(w),
		blockSize: defaultLevel * baseBlockSize,
	}
}

// NewWriterLevel is like NewWriter but uses a user supplied compression
// level instead of DefaultCompression.
//
// Valid compression levels are DefaultCompression or any integer between
// BestSpeed and BestCompression inclusively.
func NewWriterLevel(w io.Writer, level int) (*Writer, error) {
	if level != DefaultCompression && level < BestSpeed || level > BestCompression {
		return nil, ErrInvalidCompressionLevel
	}

	if level == DefaultCompression {
		level = defaultLevel
	}

	return &Writer{
		bw:        buffbits.NewWriter(w),
		blockSize: level * baseBlockSize,
	}, nil
}

// Err returns the first error that was encountered by the Writer.
func (w *Writer) Err() error {
	return w.err
}

// Reset discards any state and switches writing to the provided writer.
func (w *Writer) Reset(writer io.Writer) {
	w.bw.Reset(writer)
	w.checksum = 0
	w.wroteHeader = false
	w.err = nil
}

// Close closes the Writer by flushing any unwritten data to the
// underlying io.Writer, but does not close the underlying io.Writer.
func (w *Writer) Close() error {
	if w.err != nil {
		return w.err
	}

	w.err = w.writeHeader()
	if w.err != nil {
		return w.err
	}

	w.err = w.writeFooter()
	if w.err != nil {
		return w.err
	}

	w.err = w.bw.Flush()
	if w.err != nil {
		return w.err
	}

	w.err = ErrClosed
	return nil
}

// Write writes the compressed form of data to the writer,
// which may be buffered until the caller closes the writer.
func (w *Writer) Write(data []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	w.err = w.writeHeader()
	if w.err != nil {
		return 0, w.err
	}

	return len(data), nil
}

func (w *Writer) writeHeader() error {
	if w.wroteHeader {
		return nil
	}
	w.wroteHeader = true

	// Header is written in ascii.
	w.bw.Write(beginStreamMagic, 16)
	w.bw.Write('h', 8)
	w.bw.Write(uint64('0'+w.blockSize/baseBlockSize), 8)

	return w.bw.Err()
}

func (w *Writer) writeFooter() error {
	w.bw.Write(endStreamMagic, 48)
	w.bw.Write(uint64(w.checksum), 32)

	return w.bw.Err()
}
