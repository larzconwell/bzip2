package bzip2

import (
	"math/rand"
	"testing"
	"time"
)

func TestBWTransformEven(t *testing.T) {
	src := []byte("banana")
	dst := make([]byte, len(src))

	idx := bwTransform(dst, src)
	if idx != 3 {
		t.Error("Value idx is incorrect. Got", idx, "wanted 3")
	}

	if string(dst) != "nnbaaa" {
		t.Error("Output is incorrect. Got", string(dst), "wanted nnbaaa")
	}
}

func TestBWTransformOdd(t *testing.T) {
	src := []byte("baanana")
	dst := make([]byte, len(src))

	idx := bwTransform(dst, src)
	if idx != 4 {
		t.Error("Value idx is incorrect. Got", idx, "wanted 4")
	}

	if string(dst) != "bnnaaaa" {
		t.Error("Output is incorrect. Got", string(dst), "wanted bnnaaaa")
	}
}

func TestBWTransformAfterRLE(t *testing.T) {
	src := []byte("baaaa\x00nana")
	dst := make([]byte, len(src))

	idx := bwTransform(dst, src)
	if idx != 7 {
		t.Error("Value idx is incorrect. Got", idx, "wanted 7")
	}

	if string(dst) != "aaaabnnaa\x00" {
		t.Error("Output is incorrect. Got", string(dst), "wanted aaaabnnaa\\0")
	}
}

func BenchmarkBWTransform(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	src := make([]byte, 100000)
	dst := make([]byte, len(src))
	for i := range src {
		src[i] = byte(rand.Intn(256))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bwTransform(dst, src)
	}
}

func BenchmarkBWTransformLarge(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	src := make([]byte, 100000*6)
	dst := make([]byte, len(src))
	for i := range src {
		src[i] = byte(rand.Intn(256))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bwTransform(dst, src)
	}
}
