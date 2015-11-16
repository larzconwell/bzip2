package bits

import (
	"bytes"
	"testing"
)

func TestWriteBits(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	w.WriteBits(4, 11)
	w.WriteBits(4, 13)
	if buf.Len() != 1 {
		t.Error("First byte should have been written but didn't")
	}

	w.WriteBits(5, 22)
	w.WriteBits(7, 93)
	if buf.Len() != 2 {
		t.Error("Second byte should have been written but didn't")
	}

	w.WriteBits(4, 2)
	w.WriteBits(11, 1458)
	if buf.Len() != 4 {
		t.Error("Bytes should have been written but didn't")
	}

	w.WriteBits(5, 16)
	if buf.Len() != 5 {
		t.Error("Last byte should have been written but didn't")
	}

	expected := []byte{'\xbd', '\xb5', '\xd2', '\xb6', '\x50'}
	for i, actual := range buf.Bytes() {
		if actual != expected[i] {
			t.Error("Byte doesn't match expected value")
		}
	}
}

func TestWriteBytes(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	w.WriteBytes([]byte{'\x01'})
	if buf.Len() != 1 {
		t.Error("Wrong number of bytes written")
	}
}

func TestWriteBytesUnfinished(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	w.WriteBits(4, 8)
	w.WriteBytes([]byte{'\x01'})
	err := w.Err()
	if err == nil {
		t.Error("No error returned for unfinished bit write")
	}
	if err != ErrWriteBufferedBits {
		t.Fatal(err)
	}
}

func TestMixWriteWriteBits(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	w.WriteBits(4, 8)
	w.WriteBits(4, 7)
	w.WriteBytes([]byte{'\x55', '\x64'})
	w.WriteBits(8, 189)
	if buf.Len() != 4 {
		t.Error("Wrong number of bits and bytes written")
	}

	expected := []byte{'\x87', '\x55', '\x64', '\xbd'}
	for i, actual := range buf.Bytes() {
		if actual != expected[i] {
			t.Error("Byte doesn't match expected value")
		}
	}
}
