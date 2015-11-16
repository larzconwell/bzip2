package crc32

var (
	poly  uint32 = 0x04c11db7
	table [256]uint32
)

// Build the byte table from the polynomial.
func init() {
	for i := range table {
		c := uint32(i) << 24

		for j := 0; j < 8; j++ {
			if c&0x80000000 == 0 {
				c <<= 1
			} else {
				c = (c << 1) ^ poly
			}
		}

		table[i] = c
	}
}

// Update returns the result of adding the bytes in p to the crc.
func Update(crc uint32, p []byte) uint32 {
	crc = ^crc

	for _, b := range p {
		crc = (crc << 8) ^ table[byte(crc>>24)^b]
	}

	return ^crc
}
