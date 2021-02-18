package crc32

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	raw string
	crc uint32
}

// These crcs were created using code from the bzip2 source.
var (
	tests = []test{
		{raw: "UNIX is basically a simple operating system, but you have to be a genius to understand the simplicity. - Dennis Ritchie", crc: 0x9ba345ab},
		{raw: "Code is like humor. When you have to explain it, it’s bad. - Cory House", crc: 0x88e217ff},
		{raw: "It works on my machine.", crc: 0x70c0b4c3},
		{raw: "Cosmic rays must have caused a bit flip.", crc: 0xd6180a5f},
		{raw: "Nano? Real programmers use Vim.", crc: 0x89a39dbf},
		{raw: "There’s no place like 127.0.0.1.", crc: 0x92966114},
		{raw: "There’s no place like ::1.", crc: 0xd3d80dea},
		{raw: "Don't be clever.", crc: 0x67b46868},
	}
	combinedCRC uint32 = 0x57f00daf
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	for _, test := range tests {
		var actual uint32
		data := []byte(test.raw)

		// crcs should result in the same value
		// regardless of the number of updates.
		for len(data) != 0 {
			count := rand.Intn(len(data)) + 1
			batch := data[0:count]
			data = data[count:]

			actual = Update(actual, batch)
		}

		assert.Equal(t, test.crc, actual)
	}
}

func TestCombine(t *testing.T) {
	t.Parallel()

	var actual uint32

	for _, test := range tests {
		actual = Combine(actual, test.crc)
	}

	assert.Equal(t, combinedCRC, actual)
}
