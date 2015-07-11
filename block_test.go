package bzip2

import (
	"strconv"
	"testing"
)

// noRunData produces data to write with no runs in it.
func noRunData(size int) []byte {
	data := make([]byte, size)
	b := 0

	for i := range data {
		data[i] = byte(b)

		b++
		if b > 255 {
			b = 0
		}
	}

	return data
}

func TestBlockFullWrite(t *testing.T) {
	block := newBlock(1000)

	_, err := block.Write(noRunData(block.size))
	if err == nil {
		t.Error("Block full size write should return size reached error")
	}
}

func TestBlockMultiWrite(t *testing.T) {
	block := newBlock(1000)

	n, err := block.Write(noRunData(block.size / 2))
	if err != nil {
		t.Fatal(err)
	}

	if n != block.size/2 {
		t.Error("Block half write wrote unexpected number of bytes. Got " + strconv.Itoa(n) + " wanted " + strconv.Itoa(block.size/2))
	}

	_, err = block.Write(noRunData(block.size / 2))
	if err == nil {
		t.Error("Block full size write should return size reached error")
	}
}

func TestBlockOverWrite(t *testing.T) {
	block := newBlock(1000)

	n, err := block.Write(noRunData(block.size + 500))
	if err == nil {
		t.Error("Block full size write should return size reached error")
	}

	if n != block.size {
		t.Error("Block write wrote unexpected number of bytes. Got " + strconv.Itoa(n) + " wanted " + strconv.Itoa(block.size))
	}
}
