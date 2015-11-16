package bwt

import (
	"math/rand"
	"testing"
	"time"
)

func TestTransformEven(t *testing.T) {
	src := []byte("banana")
	dst := make([]byte, len(src))

	idx := Transform(dst, src)
	if idx != 3 {
		t.Error("Value idx is incorrect. Got", idx, "wanted 3")
	}

	if string(dst) != "nnbaaa" {
		t.Error("Output is incorrect. Got", string(dst), "wanted nnbaaa")
	}
}

func TestTransformOdd(t *testing.T) {
	src := []byte("baanana")
	dst := make([]byte, len(src))

	idx := Transform(dst, src)
	if idx != 4 {
		t.Error("Value idx is incorrect. Got", idx, "wanted 4")
	}

	if string(dst) != "bnnaaaa" {
		t.Error("Output is incorrect. Got", string(dst), "wanted bnnaaaa")
	}
}

func TestTransformAfterRLE(t *testing.T) {
	src := []byte("baaaa\x00nana")
	dst := make([]byte, len(src))

	idx := Transform(dst, src)
	if idx != 7 {
		t.Error("Value idx is incorrect. Got", idx, "wanted 7")
	}

	if string(dst) != "aaaabnnaa\x00" {
		t.Error("Output is incorrect. Got", string(dst), "wanted aaaabnnaa\\0")
	}
}

func BenchmarkTransform(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	src := make([]byte, 100000)
	dst := make([]byte, len(src))
	for i := range src {
		src[i] = byte(rand.Intn(256))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Transform(dst, src)
	}
}

func BenchmarkTransformLarge(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	src := make([]byte, 100000*6)
	dst := make([]byte, len(src))
	for i := range src {
		src[i] = byte(rand.Intn(256))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Transform(dst, src)
	}
}
