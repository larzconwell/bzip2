package bzip2

// rl2Encode encodes src using the RLE2 format after the MTF transformation.
// The return type uses uint16 because it's possible to write values greater
// than 255 if the full byte range is used.
func rl2Encode(symbols [256]int, src []byte) []uint16 {
	var repeats uint
	dst := make([]uint16, 0, len(src))

	numSymbols := uint16(0)
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
		dst = append(dst, uint16(b)+1)
		repeats = 0
	}

	finishRun()
	return append(dst, numSymbols+1)
}
