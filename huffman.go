package bzip2

import (
	"container/heap"
	"math"
)

// huffmanCode contains the bit code for a value.
type huffmanCode struct {
	n    int
	Bits uint64
}

func (hc huffmanCode) Len() int {
	return hc.n
}

// huffmanNode is a single node in the tree.
type huffmanNode struct {
	Left  *huffmanNode
	Right *huffmanNode

	Value     uint16
	Frequency int
}

// Leaf is used to check if the node has children.
func (hn huffmanNode) Leaf() bool {
	return hn.Left == nil && hn.Right == nil
}

// huffmanNodeQueue is a priority queue that keeps track of nodes.
type huffmanNodeQueue []*huffmanNode

// Let gets the number of items in the queue.
func (hnq huffmanNodeQueue) Len() int {
	return len(hnq)
}

// Less checks if the priority of node i is less than the priority of node j.
func (hnq huffmanNodeQueue) Less(i, j int) bool {
	return hnq[i].Frequency < hnq[j].Frequency
}

// Swap swaps the nodes i and j.
func (hnq huffmanNodeQueue) Swap(i, j int) {
	hnq[i], hnq[j] = hnq[j], hnq[i]
}

// Push pushes a new node to the queue.
func (hnq *huffmanNodeQueue) Push(x interface{}) {
	*hnq = append(*hnq, x.(*huffmanNode))
}

// Pop pops the lowest priority node from the queue.
func (hnq *huffmanNodeQueue) Pop() interface{} {
	queue := *hnq
	n := len(queue)
	node := queue[n-1]

	*hnq = queue[:n-1]
	return node
}

// huffmanTree is a binary tree that is navigated to produce bits for
// the frequencies of symbols.
type huffmanTree struct {
	Codes []*huffmanCode
	root  *huffmanNode
}

// newHuffmanTree creates the huffman tree and gets the codes for the symbol
// frequencies given.
func newHuffmanTree(freq []int) *huffmanTree {
	tree := &huffmanTree{Codes: make([]*huffmanCode, len(freq))}

	var queue huffmanNodeQueue
	for i, f := range freq {
		queue = append(queue, &huffmanNode{Value: uint16(i), Frequency: f})
	}
	heap.Init(&queue)

	// As long as we have multiple nodes, remove the two lowest frequency
	// symbols and create a new node with them as children.
	for queue.Len() > 1 {
		left := heap.Pop(&queue).(*huffmanNode)
		right := heap.Pop(&queue).(*huffmanNode)

		heap.Push(&queue, &huffmanNode{
			Left:      left,
			Right:     right,
			Frequency: left.Frequency + right.Frequency,
		})
	}

	tree.root = heap.Pop(&queue).(*huffmanNode)
	tree.getCodes(tree.root, 0, 0)
	return tree
}

// getCodes finds the codes for the frequencies.
func (ht huffmanTree) getCodes(node *huffmanNode, n int, bits uint64) {
	if node.Leaf() {
		ht.Codes[node.Value] = &huffmanCode{n: n, Bits: bits}
		return
	}

	n++
	ht.getCodes(node.Left, n, (bits<<1)|0)
	ht.getCodes(node.Right, n, (bits<<1)|1)
}

// generateHuffmanTrees generates the trees required to encode the data, and
// which tree to use for each 50 symbol block of data in src.
func generateHuffmanTrees(freq []int, src []uint16) ([]*huffmanTree, []int) {
	// Get the number of huffman tree changes.
	numSelectors := 1
	if len(src) > 50 {
		numSelectors = int(math.Ceil(float64(len(src)) / 50))
	}

	// Get the number of huffman trees to use.
	numHuffmanTrees := 2
	if numSelectors > 6 {
		numHuffmanTrees = 6
	} else if numSelectors > 2 {
		numHuffmanTrees = numSelectors
	}

	// Create the huffman trees generating the codes for the frequencies.
	trees := make([]*huffmanTree, numHuffmanTrees)
	for i := range trees {
		trees[i] = newHuffmanTree(freq)
	}

	// Get the tree index to use for each 50 symbol block.
	idx := 0
	treeIndexes := make([]int, numSelectors)
	for i := range treeIndexes {
		treeIndexes[i] = idx

		idx++
		if idx == numHuffmanTrees {
			idx = 0
		}
	}

	return trees, treeIndexes
}
