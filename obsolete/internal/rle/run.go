package rle

// maxRunLen is the max length for a single run.
// If a run is longer it's just spread out into
// multiple runs to simplify things.
const maxRunLen = 259

// Run contains a run byte along with its length.
type Run struct {
	Byte byte
	Len  int
}

// Encode encodes a single run.
func (r Run) Encode() []byte {
	// Less than the long encode form.
	if r.Len < 4 {
		data := make([]byte, r.Len)
		for i := range data {
			data[i] = r.Byte
		}

		return data
	}

	// Long encode form up to the max byte.
	return []byte{r.Byte, r.Byte, r.Byte, r.Byte, byte(r.Len - 4)}
}

// EncodedLen gets the length of the encoded run.
func (r Run) EncodedLen() int {
	// Less than the long encode form
	if r.Len < 4 {
		return r.Len
	}

	// Long encode form up to the max byte.
	return 5 // 4 repeats, along with a length byte.
}
