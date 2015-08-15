package bzip2

// mtfTransform performs the move-to-front transform on the src slice and
// writes the results to dst. The symbol set given is reduced to the bytes
// set to 1. Dst and src may point to the same memory.
func mtfTransform(symbols [256]int, dst, src []byte) {
	var reduced []byte
	for i, present := range symbols {
		if present > 0 {
			reduced = append(reduced, byte(i))
		}
	}

	for i, b := range src {
		symidx := 0
		for i, s := range reduced {
			if s == b {
				symidx = i
				break
			}
		}

		copy(reduced[1:], reduced[:symidx])
		reduced[0] = b
		dst[i] = byte(symidx)
	}
}
