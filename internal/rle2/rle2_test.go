package rle2

import (
	"testing"

	"github.com/larzconwell/bzip2/internal/symbols"
)

func TestFrequencies(t *testing.T) {
	data := []uint16{'\x03', '\x00', '\x03', '\x03', '\x00', '\x01', '\x04'}
	expectedFreqs := Frequencies{'\x00': 2, '\x01': 1, '\x03': 3, '\x04': 1}

	_, reduced := symbols.Get([]byte("banana"))
	freqs := GetFrequencies(reduced, data)
	for i, f := range freqs {
		if f != expectedFreqs[i] {
			t.Error("Frequency", i, "isn't the expected value", f, "should be", expectedFreqs[i])
		}
	}
}
