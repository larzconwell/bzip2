package rle

import (
	"testing"
)

func TestRunEncodeShort(t *testing.T) {
	run := Run{
		Byte: 'b',
		Len:  3,
	}
	expected := []byte("bbb")
	actual := run.Encode()

	if len(actual) != len(expected) {
		t.Error("RLE length doesn't match expected length")
	}

	for i, b := range actual {
		if b != expected[i] {
			t.Error("Byte value", string(b), "isn't the expected value",
				string(expected[i]))
		}
	}
}

func TestRunEncodeLong(t *testing.T) {
	run := Run{
		Byte: 'b',
		Len:  14,
	}
	expected := []byte("bbbb\x0a")
	actual := run.Encode()

	if len(actual) != len(expected) {
		t.Error("RLE length doesn't match expected length")
	}

	for i, b := range actual {
		if b != expected[i] {
			t.Error("Byte value", string(b), "isn't the expected value",
				string(expected[i]))
		}
	}
}

func TestRunEncodeExtraLong(t *testing.T) {
	run := Run{
		Byte: 'b',
		Len:  259,
	}
	expected := []byte("bbbb\xff")
	actual := run.Encode()

	if len(actual) != len(expected) {
		t.Error("RLE length doesn't match expected length")
	}

	for i, b := range actual {
		if b != expected[i] {
			t.Error("Byte value", string(b), "isn't the expected value",
				string(expected[i]))
		}
	}
}

func TestRunEncodedLen(t *testing.T) {
	run := Run{
		Byte: 'b',
		Len:  3,
	}
	if run.EncodedLen() != 3 {
		t.Error("RLE length doesn't match expected length, 3")
	}

	run.Len = 14
	if run.EncodedLen() != 5 {
		t.Error("RLE length doesn't match expected length, 5")
	}

	run.Len = 259
	if run.EncodedLen() != 5 {
		t.Error("RLE length doesn't match expected length, 5")
	}
}
