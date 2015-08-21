package bzip2

var (
	poly     uint32 = 0x04c11db7
	crcTable [256]uint32
)

func init() {
	// Build the byte table from the polynomial.
	for i := range crcTable {
		c := uint32(i) << 24

		for j := 0; j < 8; j++ {
			if c&0x80000000 == 0 {
				c <<= 1
			} else {
				c = (c << 1) ^ poly
			}
		}

		crcTable[i] = c
	}
}

// crcUpdate updates the given crc by adding the bytes in p to it.
func crcUpdate(crc uint32, p []byte) uint32 {
	crc = ^crc

	for _, b := range p {
		crc = (crc << 8) ^ crcTable[byte(crc>>24)^b]
	}

	return ^crc
}
