package bzip2

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestBWTransform(t *testing.T) {
	data := []byte("banana")
	ptr := bwTransform(data, data)
	if ptr != 3 {
		t.Error("Value ptr is incorrect. Got " + strconv.Itoa(ptr) + " wanted 3")
	}

	if string(data) != "nnbaaa" {
		t.Error("Output is incorrect. Got " + string(data) + " wanted nnbaaa")
	}
}

func BenchmarkBWTransform(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	src := make([]byte, 100000)
	for i := range src {
		src[i] = byte(rand.Intn(256))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst := make([]byte, len(src))
		bwTransform(dst, src)
	}
}

func BenchmarkBWTransformLarge(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	src := make([]byte, 100000*6)
	for i := range src {
		src[i] = byte(rand.Intn(256))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst := make([]byte, len(src))
		bwTransform(dst, src)
	}
}
