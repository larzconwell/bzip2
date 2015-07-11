package bzip2

import (
	"strconv"
	"testing"
)

func TestRLEncode(t *testing.T) {
	data := []byte("aeecccbhzzzzkkkkkkkkvvvvvrvv")
	expected := []byte("aeecccbhzzzz\x00kkkk\x04vvvv\x01rvv")

	out := rlEncode(data)
	if len(out) != len(expected) {
		t.Error("RLE length doesn't match expected length")
	}
	for i, d := range out {
		if d != expected[i] {
			t.Error("Byte value " + string(d) + " isn't the expected value " + string(expected[i]))
		}
	}
}

func TestRLEncodeLong(t *testing.T) {
	expected := []byte("bbbb\xFFb")
	data := make([]byte, 260)
	for i := range data {
		data[i] = 'b'
	}

	out := rlEncode(data)
	if len(out) != len(expected) {
		t.Error("RLE length doesn't match expected length")
	}
	for i, d := range out {
		if d != expected[i] {
			t.Error("Byte value " + string(d) + " isn't the expected value " + string(expected[i]))
		}
	}
}

func TestRLIndexOf(t *testing.T) {
	data := []byte("sk\x02\x02\x02\x02\x02\x02\x02")

	out := rlEncode(data)
	idx := rlIndexOf(out, len(out)-1)
	if idx != len(data)-1 {
		t.Error("Index value is incorrect. Got " + strconv.Itoa(idx) + " wanted " + strconv.Itoa(len(data)-1))
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
