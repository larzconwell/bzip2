package rle2

import (
	"github.com/larzconwell/bzip2/internal/symbols"
)

// Frequencies contains the number of times each symbol is used,
// the index being the symbol.
type Frequencies []int

// GetFrequencies gets the frequencies for a slice of symbols
func GetFrequencies(syms symbols.ReducedSet, data []uint16) Frequencies {
	freqs := make(Frequencies, len(syms)+2)

	for _, b := range data {
		freqs[b]++
	}

	return freqs
}
