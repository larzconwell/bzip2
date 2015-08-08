package bzip2

// rlEncode encodes src using the RLE format.
func rlEncode(src []byte) []byte {
	var lastb byte
	repeats := 0
	dst := make([]byte, 0, len(src))

	// finishRun writes the repeats for the last byte.
	var finishRun func()
	finishRun = func() {
		if repeats < 4 {
			for i := 0; i < repeats; i++ {
				dst = append(dst, lastb)
			}
		} else if repeats <= 259 {
			list := []byte{lastb, lastb, lastb, lastb, byte(repeats - 4)}
			dst = append(dst, list...)
		} else {
			list := []byte{lastb, lastb, lastb, lastb, byte(255)}
			dst = append(dst, list...)

			repeats -= 259
			finishRun()
		}
	}

	// Gather the repeats for the bytes in src.
	for i, b := range src {
		if i == 0 || b != lastb {
			if i > 0 {
				finishRun()
			}

			lastb = b
			repeats = 1
			continue
		}

		repeats++
	}
	finishRun()

	return dst
}

// rlIndexOf gets the index of the decoded form of data at index n.
func rlIndexOf(data []byte, n int) int {
	var lastb byte
	repeats := 0
	idx := -1

	for i := 0; i <= n; i++ {
		b := data[i]
		if i == 0 || repeats == 0 {
			lastb = b
			repeats = 1
			idx++
			continue
		}

		// Check repeats as a special case rather than a simple byte check
		// because lastb could be the same as the number of bytes in
		// the encoded length.
		if repeats == 4 {
			repeats = 0
			idx += int(b)
			continue
		}

		if b != lastb {
			lastb = b
			repeats = 1
			idx++
			continue
		}

		repeats++
		idx++
	}

	return idx
}
