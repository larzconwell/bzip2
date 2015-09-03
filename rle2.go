package bzip2

// rl2Encode encodes src using the second run-length encoding form after
// the MTF transformation. The return type uses uint16 because it's possible
// to write values greater than 255 if the full byte range is used.
func rl2Encode(symbols []byte, src []byte) []uint16 {
	var repeats uint
	dst := make([]uint16, 0, len(src))

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

		repeats = 0
	}

	for _, b := range src {
		if b == '\x00' {
			repeats++
			continue
		}

		finishRun()

		v := uint16(b) + 1
		dst = append(dst, v)
	}

	finishRun()
	return append(dst, uint16(len(symbols)+1))
}

// symbolFrequencies gets the frequencies for a slice of symbols after it has
// been encoded in the second run-length encoding form.
func symbolFrequencies(symbols []byte, data []uint16) []int {
	freq := make([]int, len(symbols)+2)

	for _, b := range data {
		freq[b]++
	}

	return freq
}
