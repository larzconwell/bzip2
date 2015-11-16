package rle

// Encode encodes src using run-length encoding.
func Encode(src []byte) []byte {
	var lastb byte
	var repeat int
	dst := make([]byte, 0, len(src))

	// Gather the repeats for the bytes in src.
	for i, b := range src {
		if i == 0 || b != lastb {
			if i > 0 {
				dst = finishRun(lastb, repeat, dst)
			}

			lastb = b
			repeat = 1
			continue
		}

		repeat++
	}

	return finishRun(lastb, repeat, dst)
}

// finishRun repeats a run repeat number of times writing to dst.
func finishRun(run byte, repeat int, dst []byte) []byte {
	if repeat < 4 {
		for i := 0; i < repeat; i++ {
			dst = append(dst, run)
		}

		return dst
	} else if repeat <= 259 {
		list := []byte{run, run, run, run, byte(repeat - 4)}
		return append(dst, list...)
	}

	list := []byte{run, run, run, run, '\xff'}
	return finishRun(run, repeat-259, append(dst, list...))
}
