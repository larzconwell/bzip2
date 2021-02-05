package huffman

import (
	"testing"

	"github.com/larzconwell/bzip2/internal/rle2"
	"github.com/larzconwell/bzip2/internal/symbols"
)

func TestGenerateTreesLowest(t *testing.T) {
	_, reduced := symbols.Get([]byte("banana"))
	data := []uint16{'\x03', '\x00', '\x03', '\x00', '\x01', '\x04'}
	freqs := rle2.GetFrequencies(reduced, data)

	trees, selections := GenerateTrees(freqs, data)
	if len(trees) < 2 {
		t.Error("Not enough huffman trees generated")
	}
	if len(trees) > 6 {
		t.Error("Too many huffman trees generated")
	}

	if len(selections) != 1 {
		t.Error("The wrong number of huffman tree selections was returned")
	}
}

func TestGenerateTreesMultipleSelections(t *testing.T) {
	_, reduced := symbols.Get([]byte("banana"))
	data := []uint16{'\x03', '\x00', '\x03', '\x00', '\x01', '\x04'}
	base := data
	for len(data) <= 200 {
		data = append(data, base...)
	}
	freqs := rle2.GetFrequencies(reduced, data)

	trees, selections := GenerateTrees(freqs, data)
	if len(trees) < 2 {
		t.Error("Not enough huffman trees generated")
	}
	if len(trees) > 6 {
		t.Error("Too many huffman trees generated")
	}

	if len(selections) != 5 {
		t.Error("The wrong number of huffman tree selections was returned")
	}
}
