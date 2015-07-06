package bzip2

import (
	"bytes"
	"testing"
)

func TestWriteBits(t *testing.T) {
	var buf bytes.Buffer
	bw := newBitWriter(&buf)
	bw.WriteBits(4, 11) // 1011
	bw.WriteBits(4, 13) // 1101
	if buf.Len() != 1 {
		t.Error("First byte should have been written but didn't")
	}

	bw.WriteBits(5, 22) // 10110
	bw.WriteBits(7, 93) // 1011101
	if buf.Len() != 2 {
		t.Error("Second byte should have been written but didn't")
	}

	bw.WriteBits(4, 2)     // 0010
	bw.WriteBits(11, 1458) // 10110110010
	if buf.Len() != 4 {
		t.Error("Bytes should have been written but didn't")
	}

	bw.WriteBits(5, 16) // 10000
	if buf.Len() != 5 {
		t.Error("Last byte should have been written but didn't")
	}

	expected := []byte{189, 181, 210, 182, 80}
	for i, got := range buf.Bytes() {
		if got != expected[i] {
			t.Error("Byte doesn't match expected value")
		}
	}
}

func TestWriteBytes(t *testing.T) {
	var buf bytes.Buffer
	bw := newBitWriter(&buf)
	bw.WriteBytes([]byte{1})
	if buf.Len() != 1 {
		t.Error("Wrong number of bytes written")
	}
}

func TestWriteBytesUnfinished(t *testing.T) {
	var buf bytes.Buffer
	bw := newBitWriter(&buf)
	bw.WriteBits(4, 8)
	bw.WriteBytes([]byte{1})
	err := bw.Err()
	if err == nil {
		t.Error("No error returned for unfinished bit write")
	}
	if err != errWriteUnfinishedBits {
		t.Fatal(err)
	}
}

func TestMixWriteWriteBits(t *testing.T) {
	var buf bytes.Buffer
	bw := newBitWriter(&buf)
	bw.WriteBits(4, 8) // 1000
	bw.WriteBits(4, 7) // 0111
	bw.WriteBytes([]byte{85, 100})
	bw.WriteBits(8, 189)
	if buf.Len() != 4 {
		t.Error("Wrong number of bits and bytes written")
	}

	expected := []byte{135, 85, 100, 189}
	for i, got := range buf.Bytes() {
		if got != expected[i] {
			t.Error("Byte doesn't match expected value")
		}
	}
}
