package mtf

import (
	"github.com/larzconwell/bzip2/internal/symbols"
)

// Transform performs the move-to-front transform on the src slice and
// writes the results to dst. Dst and src may point to the same memory.
func Transform(syms symbols.ReducedSet, dst, src []byte) {
	symbols := make(symbols.ReducedSet, len(syms))
	copy(symbols, syms)

	for i, b := range src {
		// Get the index where the byte b exists.
		symidx := 0
		for i, s := range symbols {
			if s == b {
				symidx = i
				break
			}
		}

		// Move the byte b to the front of the set.
		copy(symbols[1:], symbols[:symidx])
		symbols[0] = b

		dst[i] = byte(symidx)
	}
}
