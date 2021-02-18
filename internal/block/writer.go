package block

import (
	"github.com/larzconwell/buffbits"
)

const (
	baseSize = 100_000
)

// Writer is an io.WriteCloser that handles writing compressed block data to an underlying
// bit writer up to a set size. If an error occurs while writing to a Writer, no more data
// will be written and all subsequent calls with return an error. The caller should call
// Close when done to complete the block.
type Writer struct {
	Checksum uint32
	bw       *buffbits.Writer
	size     int
	err      error
}

// NewWriter creates a Writer that compresses and writes to bw up to level * 100k bytes.
func NewWriter(bw *buffbits.Writer, level int) *Writer {
	return &Writer{bw: bw, size: level * baseSize}
}

// Err returns the first error that was encountered by the Writer.
func (w *Writer) Err() error {
	return w.err
}

// Reset discards any state and switches writing to the provided bit writer.
func (w *Writer) Reset(bw *buffbits.Writer) {
	w.Checksum = 0
	w.bw = bw
	w.err = nil
}

// Len returns the number of bytes written to the block.
func (w *Writer) Len() int {
	return 0
}

// Write writes the compressed form of data to the Writer, which is buffered
// until the caller closes the Writer. Once the Writer has reached it's limit
// ErrLimitReached is returned along with the number of bytes from data that
// were written.
func (w *Writer) Write(data []byte) (int, error) {
	return len(data), ErrLimitReached
}

// Close closes the Writer by compressing the buffered data and writing
// to the bit writer.
func (w *Writer) Close() error {
	return nil
}
