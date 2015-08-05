package bzip2

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestBWTransform(t *testing.T) {
	src := []byte("banana")
	dst := make([]byte, len(src))

	ptr := bwTransform(dst, src)
	if ptr != 3 {
		t.Error("Value ptr is incorrect. Got " + strconv.Itoa(ptr) + " wanted 3")
	}

	if string(dst) != "nnbaaa" {
		t.Error("Output is incorrect. Got " + string(dst) + " wanted nnbaaa")
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
