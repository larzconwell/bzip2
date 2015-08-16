package bzip2

// rl2Encode encodes src using the RLE2 format after the MTF transformation.
// The return type uses uint16 because it's possible to write values greater
// than 255 if the full byte range is used. The frequency of values is also
// returned.
func rl2Encode(symbols [256]int, src []byte) ([]int, []uint16) {
	var repeats uint
	dst := make([]uint16, 0, len(src))

	numSymbols := uint16(0)
	for _, present := range symbols {
		if present > 0 {
			numSymbols++
		}
	}
	freq := make([]int, numSymbols+2)

	finishRun := func() {
		for repeats > 0 {
			if repeats&1 > 0 {
				freq['\x00']++
				dst = append(dst, '\x00')
				repeats--
			} else {
				freq['\x01']++
				dst = append(dst, '\x01')
				repeats -= 2
			}

			repeats >>= 1
		}
	}

	for _, b := range src {
		if b == '\x00' {
			repeats++
			continue
		}

		finishRun()
		repeats = 0

		v := uint16(b) + 1
		freq[v]++
		dst = append(dst, v)
	}

	finishRun()
	freq[numSymbols+1]++
	return freq, append(dst, numSymbols+1)
}
