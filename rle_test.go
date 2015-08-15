package bzip2

import (
	"testing"
)

func TestRLEncode(t *testing.T) {
	src := []byte("aeecccbhzzzzkkkkkkkkvvvvvrvv")
	expected := []byte("aeecccbhzzzz\x00kkkk\x04vvvv\x01rvv")

	dst := rlEncode(src)
	if len(dst) != len(expected) {
		t.Error("RLE length doesn't match expected length")
	}
	for i, d := range dst {
		if d != expected[i] {
			t.Error("Byte value", string(d), "isn't the expected value", string(expected[i]))
		}
	}
}

func TestRLEncodeLong(t *testing.T) {
	expected := []byte("bbbb\xFFb")
	src := make([]byte, 260)
	for i := range src {
		src[i] = 'b'
	}

	dst := rlEncode(src)
	if len(dst) != len(expected) {
		t.Error("RLE length doesn't match expected length")
	}
	for i, actual := range dst {
		if actual != expected[i] {
			t.Error("Byte value", string(actual), "isn't the expected value", string(expected[i]))
		}
	}
}

func TestRLIndexOf(t *testing.T) {
	src := []byte("sk\x02\x02\x02\x02\x02\x02\x02")

	dst := rlEncode(src)
	idx := rlIndexOf(len(dst)-1, dst)
	if idx != len(src)-1 {
		t.Error("Index value is incorrect. Got", idx, "wanted", len(src)-1)
	}
}

func TestRLIndexOfComplex(t *testing.T) {
	src := []byte("sk\x02\x02\x02\x02\x02\x02vrrrrrkc")

	dst := rlEncode(src)
	idx := rlIndexOf(len(dst)-2, dst)
	if idx != len(src)-2 {
		t.Error("Index value is incorrect. Got", idx, "wanted", len(src)-2)
	}
}

func BenchmarkRLEncode(b *testing.B) {
	data := make([]byte, 260)
	for i := range data {
		data[i] = 'b'
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rlEncode(data)
	}
}
