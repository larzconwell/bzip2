package crc32

// Update returns the result of adding the bytes in data to crc.
func Update(crc uint32, data []byte) uint32 {
	if table == nil {
		makeTable()
	}

	crc = ^crc

	for _, b := range data {
		crc = (crc << 8) ^ table[byte(crc>>24)^b]
	}

	return ^crc
}

// Combine returns the resulting crcs combined into one.
func Combine(crc, other uint32) uint32 {
	return (crc<<1 | crc>>31) ^ other
}
