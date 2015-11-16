package huffman

import (
	"container/heap"

	"github.com/larzconwell/bzip2/internal/rle2"
)

// Tree is a binary tree that is navigated to produce bits for
// the frequencies of symbols.
type Tree struct {
	Codes []*Code
	root  *Node
}

// NewTree creates a huffman tree and gets the codes for the symbol
// frequencies given.
func NewTree(freqs rle2.Frequencies) *Tree {
	tree := &Tree{Codes: make([]*Code, len(freqs))}

	var queue NodeQueue
	for i, f := range freqs {
		queue = append(queue, &Node{Value: uint16(i), Frequency: f})
	}
	heap.Init(&queue)

	// As long as we have multiple nodes, remove the two lowest frequency
	// symbols and create a new node with them as children.
	for queue.Len() > 1 {
		left := heap.Pop(&queue).(*Node)
		right := heap.Pop(&queue).(*Node)

		heap.Push(&queue, &Node{
			Left:      left,
			Right:     right,
			Frequency: left.Frequency + right.Frequency,
		})
	}

	tree.root = heap.Pop(&queue).(*Node)
	tree.getCodes(tree.root, 0, 0)
	return tree
}

// getCodes finds the codes for the frequencies.
func (t Tree) getCodes(node *Node, n int, bits uint64) {
	if node.Leaf() {
		t.Codes[node.Value] = &Code{Len: n, Bits: bits}
		return
	}

	n++
	t.getCodes(node.Left, n, (bits<<1)|0)
	t.getCodes(node.Right, n, (bits<<1)|1)
}
