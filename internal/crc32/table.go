package crc32

import (
	"hash/crc32"
)

const (
	poly uint32 = 0x04c11db7
)

var (
	table *crc32.Table
)

func makeTable() {
	table = new(crc32.Table)

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
