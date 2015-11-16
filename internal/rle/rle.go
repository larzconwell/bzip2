package rle

// IndexOf gets the index equivalent to n for the
// decoded form of data.
func IndexOf(n int, data []byte) int {
	var lastb byte
	var repeat int
	idx := -1

	for i := 0; i <= n; i++ {
		b := data[i]
		if i == 0 || repeat == 0 {
			lastb = b
			repeat = 1
			idx++
			continue
		}

		// Check repeats as a special case rather than a
		// simple byte check because lastb could be the
		// same as the number of bytes in the encoded length.
		if repeat == 4 {
			repeat = 0
			idx += int(b)
			continue
		}

		if b != lastb {
			lastb = b
			repeat = 1
			idx++
			continue
		}

		repeat++
		idx++
	}

	return idx
}
