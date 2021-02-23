package testhelpers

import (
	"io"
)

// ErrReadWriter is an io.ReadWriter that returns err for all Read and Write calls.
func ErrReadWriter(err error) io.ReadWriter {
	return errReadWriter{err: err}
}

type errReadWriter struct {
	err error
}

func (eio errReadWriter) Read([]byte) (int, error) {
	return 0, eio.err
}

func (eio errReadWriter) Write([]byte) (int, error) {
	return 0, eio.err
}
