package bzip2

import (
	"testing"
)

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

func TestBlockRLEOverWrite(t *testing.T) {
	block := newBlock(1000)

	// Create data ending with a run that over writes
	// the block size with a repeat byte. This
	// is an issue because it wouldn't be a valid RLE
	// run because the repeat byte would be stripped off.
	// If this occurs we just strip off an extra byte
	// leaving a 3 byte run which is under the limit.
	data := noRunData(996)
	b := data[994]
	data = append(data, []byte{b, b, b, b}...)

	n, err := block.Write(data)
	if err == nil {
		t.Error("Block full size write should return size reached error")
	}
	if err != errBlockSizeReached {
		t.Fatal(err)
	}

	if n != block.size-1 {
		t.Error("Block write wrote unexpected number of bytes. Got", n,
			"wanted", block.size-1)
	}

	data = data[n:]
	if len(data) != 1 {
		t.Error("Wrong number of bytes left after overwrite. Got", len(data),
			"wanted", 1)
	}
}
