package bzip2

import (
	"testing"
)

func TestEncodeRLE(t *testing.T) {
	src := []byte("aeecccbhzzzzkkkkkkkkvvvvvrvv")
	expected := []byte("aeecccbhzzzz\x00kkkk\x04vvvv\x01rvv")

	dst := encodeRLE(src)
	if len(dst) != len(expected) {
		t.Error("RLE length doesn't match expected length")
	}
	for i, d := range dst {
		if d != expected[i] {
			t.Error("Byte value " + string(d) + " isn't the expected value " + string(expected[i]))
		}
	}
}

func TestEncodeRLELong(t *testing.T) {
	expected := []byte("bbbb\xFBbbbb\x01")
	src := make([]byte, 260)
	for i := range src {
		src[i] = 'b'
	}

	dst := encodeRLE(src)
	if len(dst) != len(expected) {
		t.Error("RLE length doesn't match expected length")
	}
	for i, d := range dst {
		if d != expected[i] {
			t.Error("Byte value " + string(d) + " isn't the expected value " + string(expected[i]))
		}
	}
}

func BenchmarkEncodeRLE(b *testing.B) {
	src := make([]byte, 260)
	for i := range src {
		src[i] = 'b'
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encodeRLE(src)
	}
}
