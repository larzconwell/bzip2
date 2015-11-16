package huffman

// Node is a single node in the tree.
type Node struct {
	Left  *Node
	Right *Node

	Value     uint16
	Frequency int
}

// Leaf is used to check if the node has children.
func (n Node) Leaf() bool {
	return n.Left == nil && n.Right == nil
}

// NodeQueue is a priority queue that keeps track of nodes.
type NodeQueue []*Node

// Len gets the number of nodes in the queue.
func (nq NodeQueue) Len() int {
	return len(nq)
}

// Less checks if the priority of node i is less than the priority of node j.
func (nq NodeQueue) Less(i, j int) bool {
	return nq[i].Frequency < nq[j].Frequency
}

// Swap swaps the nodes i and j.
func (nq NodeQueue) Swap(i, j int) {
	nq[i], nq[j] = nq[j], nq[i]
}

// Push pushes a new node to the queue.
func (nq *NodeQueue) Push(x interface{}) {
	*nq = append(*nq, x.(*Node))
}

// Pop pops the lowest priority node from the queue.
func (nq *NodeQueue) Pop() interface{} {
	queue := *nq
	n := len(queue)
	node := queue[n-1]

	*nq = queue[:n-1]
	return node
}
