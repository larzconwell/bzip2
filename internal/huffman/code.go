package huffman

// Code contains the bit code for a symbol.
type Code struct {
	n    int
	Bits uint64
}

// Len gets the length of the bit code.
func (c Code) Len() int {
	return c.n
}
