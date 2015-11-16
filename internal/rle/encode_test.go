package rle

import (
	"math/rand"
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	src := []byte("aeecccbhzzzzkkkkkkkkvvvvvrvv")
	expected := []byte("aeecccbhzzzz\x00kkkk\x04vvvv\x01rvv")

	dst := Encode(src)
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

func TestEncodeLong(t *testing.T) {
	expected := []byte("bbbb\xffb")
	src := make([]byte, 260)
	for i := range src {
		src[i] = 'b'
	}

	dst := Encode(src)
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

func BenchmarkEncode(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	data := make([]byte, 100000)
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(data)
	}
}
