package huffman

import (
	"testing"

	"github.com/larzconwell/bzip2/internal/rle2"
	"github.com/larzconwell/bzip2/internal/symbols"
)

func TestTreeCodeLength(t *testing.T) {
	_, reduced := symbols.Get([]byte("banana"))
	data := []uint16{'\x03', '\x00', '\x03', '\x00', '\x01', '\x04'}
	freqs := rle2.GetFrequencies(reduced, data)

	lowestlen := 0
	tree := NewTree(freqs)
	for i, code := range tree.Codes {
		if i == 0 || code.Len < lowestlen {
			lowestlen = code.Len
		}
	}

	if lowestlen != tree.Codes['\x03'].Len {
		t.Error("The lowest code-length isn't the most used symbol")
	}
}
