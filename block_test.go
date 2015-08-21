package bzip2

import (
	"testing"
)

// noRunData produces data to write with no runs in it.
func noRunData(size int) []byte {
	data := make([]byte, size)
	b := byte('\x00')

	for i := range data {
		data[i] = b

		b++
	}

	return data
}

func TestBlockFullWrite(t *testing.T) {
	block := newBlock(1000)

	_, err := block.Write(noRunData(block.size))
	if err == nil {
		t.Error("Block full size write should return size reached error")
	}
	if err != errBlockSizeReached {
		t.Fatal(err)
	}
}

func TestBlockMultiWrite(t *testing.T) {
	block := newBlock(1000)

	n, err := block.Write(noRunData(block.size / 2))
	if err != nil {
		t.Fatal(err)
	}

	if n != block.size/2 {
		t.Error("Block half write wrote unexpected number of bytes. Got", n,
			"wanted", block.size/2)
	}

	_, err = block.Write(noRunData(block.size / 2))
	if err == nil {
		t.Error("Block full size write should return size reached error")
	}
	if err != errBlockSizeReached {
		t.Fatal(err)
	}
}

func TestBlockOverWrite(t *testing.T) {
	block := newBlock(1000)

	n, err := block.Write(noRunData(block.size + 500))
	if err == nil {
		t.Error("Block full size write should return size reached error")
	}
	if err != errBlockSizeReached {
		t.Fatal(err)
	}

	if n != block.size {
		t.Error("Block write wrote unexpected number of bytes. Got", n,
			"wanted", block.size)
	}
}

func TestSymbolSet(t *testing.T) {
	symbols, reducedSymbols := symbolSet([]byte("banana"))
	if string(reducedSymbols) != "abn" {
		t.Error("The reduced symbol set doesn't include the correct bytes")
	}

	for i, present := range symbols {
		if present > 1 && (i != 'a' || i != 'b' || i != 'n') {
			t.Error("Symbol set includes a byte that should be set")
		}
	}
}
