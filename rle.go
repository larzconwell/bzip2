package bzip2

// encodeRLE encodes data using the RLE format.
func encodeRLE(src []byte) []byte {
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
		} else if repeats <= 255 {
			list := []byte{lastb, lastb, lastb, lastb, byte(repeats - 4)}
			dst = append(dst, list...)
		} else {
			list := []byte{lastb, lastb, lastb, lastb, byte(251)}
			dst = append(dst, list...)

			repeats -= 255
			finishRun()
		}
	}

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
