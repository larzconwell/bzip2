package bzip2

import (
	"errors"
	"io"
)

var (
	errWriteUnfinishedBits = errors.New("bzip2: WriteBytes with unfinished bits")
)

// bitWriter wraps an io.Writer and provides the ability to write values
// bit-by-bit to it. It's Write* methods don't return the usual error because
// error handling is verbose. Instead, any error is kept and can be checked
// afterwards.
type bitWriter struct {
	w    io.Writer
	bits uint64
	n    uint
	err  error
}

// newbitWriter creates a bit writer writing to w.
func newBitWriter(w io.Writer) *bitWriter {
	return &bitWriter{w: w}
}

// WriteBytes writes bytes to the writer.
//
// If there are any pending bits in the buffer it is an error
// to write bytes.
func (bw *bitWriter) WriteBytes(b []byte) {
	if bw.err != nil {
		return
	}

	if bw.n != 0 {
		bw.err = errWriteUnfinishedBits
		return
	}

	_, bw.err = bw.w.Write(b)
	return
}

// WriteBits writes n bits to the writer, if n is less than a byte buffering
// may occur until enough bits are given to write.
func (bw *bitWriter) WriteBits(n uint, bits uint64) {
	if bw.err != nil {
		return
	}
	total := bw.n + n

	// Not enough bits to write, store until later.
	if total < 8 {
		bw.n = total
		bw.bits = (bw.bits << n) | bits
		return
	}

	// Bytes exists, but some bits may be left.
	bits = (bw.bits << n) | bits
	bw.n = total % 8
	bw.bits = bits & (1<<bw.n - 1)
	bits = (bits ^ uint64(bw.n)) >> bw.n

	list := byteList(total-bw.n, bits)
	_, bw.err = bw.w.Write(list)
}

// Buffered gets the number of buffered bits.
func (bw bitWriter) Buffered() uint {
	return bw.n
}

// Err gets the error for the bit writer.
func (bw bitWriter) Err() error {
	return bw.err
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
