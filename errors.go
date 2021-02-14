package bzip2

import (
	"errors"
)

var (
	// ErrInvalidCompressionLevel is returned from NewWriterLevel when an invalid compression level is used.
	ErrInvalidCompressionLevel = errors.New("bzip2: invalid compression level")

	// ErrClosed is returned when the stream has been closed previously.
	ErrClosed = errors.New("bzip2: stream already closed")
)
