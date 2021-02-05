package rle

// RunList contains a list of runs from data as they occurred.
type RunList struct {
	runs       []*Run
	encodedlen int
}

// NewRunList creates a RunList ready to read runs.
func NewRunList() *RunList {
	return &RunList{runs: make([]*Run, 0, 100)}
}

// Update updates the runs with the data given. Returning
// the encoded length of the runs afterwards.
func (rl *RunList) Update(data []byte) int {
	// The last run in the list is used if one exists
	// so that runs across updates are correct.
	var lastRun *Run
	if len(rl.runs) > 0 {
		lastRun = rl.runs[len(rl.runs)-1]

		// Remove the last runs encoded length. It'll be added later.
		rl.encodedlen -= lastRun.EncodedLen()
	}

	for _, b := range data {
		if lastRun == nil || b != lastRun.Byte || lastRun.Len == maxRunLen {
			if lastRun != nil {
				rl.encodedlen += lastRun.EncodedLen()
			}

			lastRun = &Run{
				Byte: b,
				Len:  1,
			}

			rl.runs = append(rl.runs, lastRun)
			continue
		}

		lastRun.Len++
	}
	if lastRun != nil {
		rl.encodedlen += lastRun.EncodedLen()
	}

	return rl.encodedlen
}

// Trim trims a number of bytes from the encoded form of the run list.
// The number of bytes trimmed from the original input is returned.
func (rl *RunList) Trim(n int) int {
	trimmed := 0

	for i := len(rl.runs) - 1; i >= 0; i-- {
		run := rl.runs[i]
		encodedlen := run.EncodedLen()

		// If n is bigger than the runs encoded length
		// just remove the run entirely.
		if n >= encodedlen {
			rl.runs = rl.runs[:i]
			trimmed += run.Len

			n -= encodedlen
			if n == 0 {
				break
			}

			continue
		}

		// Less than the long encode form, just remove
		// them from the length.
		if encodedlen < 4 {
			trimmed += n
			run.Len -= n
			break
		}

		// Long encode form, subtracting the encoding length
		// with the number of bytes to remove will yield
		// the new length before encoding, since the only
		// valid results are 1-3, if the result is 4 that
		// isn't a valid encoded length so we have to trim
		// an extra byte.
		newencodedlen := encodedlen - n
		if newencodedlen == 4 {
			newencodedlen--
		}

		trimmed += run.Len - newencodedlen
		run.Len = newencodedlen
		break
	}

	return trimmed
}

// Encode encodes the runs reconstructing the original data in
// the encoded form.
func (rl RunList) Encode() []byte {
	var data []byte

	for _, run := range rl.runs {
		data = append(data, run.Encode()...)
	}

	return data
}

// EncodedLen gets the encoded length of the runs.
func (rl RunList) EncodedLen() int {
	return rl.encodedlen
}
