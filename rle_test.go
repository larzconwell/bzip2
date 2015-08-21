package bzip2

import (
	"math/rand"
	"testing"
	"time"
)

func TestRLEncode(t *testing.T) {
	src := []byte("aeecccbhzzzzkkkkkkkkvvvvvrvv")
	expected := []byte("aeecccbhzzzz\x00kkkk\x04vvvv\x01rvv")

	dst := rlEncode(src)
	if len(dst) != len(expected) {
		t.Error("RLE length doesn't match expected length")
	}
	for i, b := range dst {
		if b != expected[i] {
			t.Error("Byte value", string(b), "isn't the expected value",
				string(expected[i]))
		}
	}
}

func TestRLEncodeLong(t *testing.T) {
	expected := []byte("bbbb\xffb")
	src := make([]byte, 260)
	for i := range src {
		src[i] = 'b'
	}

	dst := rlEncode(src)
	if len(dst) != len(expected) {
		t.Error("RLE length doesn't match expected length")
	}
	for i, b := range dst {
		if b != expected[i] {
			t.Error("Byte value", string(b), "isn't the expected value",
				string(expected[i]))
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
	rand.Seed(time.Now().UnixNano())

	data := make([]byte, 100000)
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rlEncode(data)
	}
}
