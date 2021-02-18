package block

import (
	"errors"
)

var (
	// ErrLimitReached is returned once the blocks limit has been reached.
	ErrLimitReached = errors.New("bzip2: block size limit reached")
)
