package symbols

// Set contains all possible bytes and if they've been set.
type Set [256]int

// ReducedSet contains all bytes that are set.
type ReducedSet []byte

// Get gets the symbol set for a slice of bytes.
func Get(data []byte) (Set, ReducedSet) {
	var symbols Set
	reduced := make(ReducedSet, 0, 256)

	for _, b := range data {
		symbols[b] = 1
	}

	for i, present := range symbols {
		if present > 0 {
			reduced = append(reduced, byte(i))
		}
	}

	return symbols, reduced
}
