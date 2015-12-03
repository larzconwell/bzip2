package bits

import (
	"io"
)

// Writer wraps an io.Writer and provides the ability to write
// values bit-by-bit to it. It's Write* methods don't return the
// usual error because error handling is verbose. Instead, any
// error is kept and can be checked afterwards.
type Writer struct {
	w    io.Writer
	bits uint64
	n    uint
	err  error
}

// NewWriter creates a bit writer writing to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

// WriteBits writes n bits to the writer, if n is less than a byte
// buffering may occur until enough bits are given to write.
func (w *Writer) WriteBits(n uint, bits uint64) {
	if w.err != nil {
		return
	}
	total := w.n + n

	// Not enough bits to write, store until later.
	if total < 8 {
		w.n = total
		w.bits = (w.bits << n) | bits
		return
	}

	// Bytes exists, but some bits may be left.
	bits = (w.bits << n) | bits
	w.n = total % 8
	w.bits = bits & (1<<w.n - 1)
	bits = (bits ^ uint64(w.n)) >> w.n

	b := byteList(total-w.n, bits)
	_, w.err = w.w.Write(b)
}

// Buffered gets the number of buffered bits.
func (w Writer) Buffered() uint {
	return w.n
}

// Err gets the error for the bit writer.
func (w Writer) Err() error {
	return w.err
}

// byteList converts n bits to a byte slice.
func byteList(n uint, bits uint64) []byte {
	total := n / 8
	list := make([]byte, 0, total)

	for i := uint(1); i <= total; i++ {
		// Shift the bits over to LSB and then conversion to byte
		// masks anything left above.
		bits := byte(bits >> uint64((total-i)*8))

		list = append(list, bits)
	}

	return list
}
