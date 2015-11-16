package huffman

import (
	"math"

	"github.com/larzconwell/bzip2/internal/rle2"
)

// TreeSelectionLimit is the symbol limit for each tree selection.
const TreeSelectionLimit = 50

// GenerateTrees generates the trees required to encode the data, and
// which tree to use for each 50 symbol block of data in src.
func GenerateTrees(freqs rle2.Frequencies, src []uint16) ([]*Tree, []int) {
	// Get the number of huffman tree selections.
	numSelections := 1
	if len(src) > TreeSelectionLimit {
		numSelections = int(math.Ceil(float64(len(src)) / TreeSelectionLimit))
	}

	// Get the number of trees to use.
	numTrees := 2
	if numSelections > 6 {
		numTrees = 6
	} else if numSelections > 2 {
		numTrees = numSelections
	}

	// Create the huffman trees generating the codes for the frequencies.
	trees := make([]*Tree, numTrees)
	for i := range trees {
		trees[i] = NewTree(freqs)
	}

	// Get the tree selection to use for each 50 symbol block.
	idx := 0
	selections := make([]int, numSelections)
	for i := range selections {
		selections[i] = idx

		idx++
		if idx == numTrees {
			idx = 0
		}
	}

	return trees, selections
}
