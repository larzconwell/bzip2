package bzip2

// mtfTransform performs the move-to-front transform on the src slice and
// writes the results to dst. Dst and src may point to the same memory.
func mtfTransform(symbols []byte, dst, src []byte) {
	reduced := make([]byte, len(symbols))
	copy(reduced, symbols)

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
