package bzip2

// rl2Encode encodes src using the RLE2 format
// after the MTF transformation.
func rl2Encode(symbols []int, src []byte) []byte {
	var repeats uint
	dst := make([]byte, 0, len(src))

	numSymbols := 0
	for _, present := range symbols {
		if present > 0 {
			numSymbols++
		}
	}

	finishRun := func() {
		for repeats > 0 {
			if repeats&1 > 0 {
				dst = append(dst, '\x00')
				repeats--
			} else {
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
		dst = append(dst, b+1)
		repeats = 0
	}

	finishRun()
	return append(dst, byte(numSymbols+1))
}
