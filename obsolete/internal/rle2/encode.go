package rle2

import (
	"github.com/larzconwell/bzip2/internal/symbols"
)

// Encode encodes src using the second run-length encoding
// form after the MTF transformation. The return type uses
// uint16 because it's possible to write values greater than
// 255 if the full byte range is used.
func Encode(syms symbols.ReducedSet, src []byte) []uint16 {
	var repeat int
	dst := make([]uint16, 0, len(src))

	for _, b := range src {
		if b == '\x00' {
			repeat++
			continue
		}

		dst = finishRun(repeat, dst)
		repeat = 0

		v := uint16(b) + 1
		dst = append(dst, v)
	}

	dst = finishRun(repeat, dst)
	return append(dst, uint16(len(syms)+1))
}

// finishRun performs the run-length encoding for a repeat
// of repeat times.
func finishRun(repeat int, dst []uint16) []uint16 {
	for repeat > 0 {
		if repeat&1 > 0 {
			dst = append(dst, '\x00')
			repeat--
		} else {
			dst = append(dst, '\x01')
			repeat -= 2
		}

		repeat >>= 1
	}

	return dst
}
