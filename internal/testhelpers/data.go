package testhelpers

// NoRunData produces data to write with no runs in it.
func NoRunData(size int) []byte {
	var b byte
	data := make([]byte, size)

	for i := range data {
		data[i] = b
		b++
	}

	return data
}
