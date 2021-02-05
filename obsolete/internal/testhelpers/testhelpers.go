package testhelpers

import (
	"math/rand"
	"time"
)

// NoRunData produces data to write with no runs in it.
func NoRunData(size int) []byte {
	data := make([]byte, size)
	b := byte('\x00')

	for i := range data {
		data[i] = b

		b++
	}

	return data
}

// RandomRunData produces data to write with random bytes
// including runs.
func RandomRunData(size int) []byte {
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, size)

	for i := range data {
		b := byte(rand.Intn(256))

		// Get the last bytes run length.
		runb := byte('\x00')
		runlen := 0
		if i > 0 {
			for j := i - 1; j >= 0; j-- {
				if j == i-1 {
					runb = data[j]
					runlen = 1
					continue
				}

				if data[j] != runb {
					break
				}

				runlen++
			}
		}

		// Detect if we shoud make the run longer.
		if shouldIncRun(runlen) {
			b = runb
		}

		data[i] = b
	}

	return data
}

// shouldIncRun decides if the run should be made longer.
func shouldIncRun(runlen int) bool {
	if runlen == 0 {
		return false
	}
	trues := 0
	falses := 0

	// Count runlen number of random trues/falses.
	for i := 0; i < runlen; i++ {
		if rand.Intn(2) == 1 {
			trues++
		} else {
			falses++
		}
	}

	if trues > falses {
		return true
	} else if trues < falses {
		return false
	}

	return rand.Intn(2) == 1
}
