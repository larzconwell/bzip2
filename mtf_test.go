package bzip2

import (
	"math/rand"
	"testing"
	"time"
)

// symbolSet gets the symbol set for a slice of bytes.
func symbolSet(data []byte) []int {
	symbols := make([]int, 256)
	for _, b := range data {
		symbols[int(b)] = 1
	}

	return symbols
}

func TestMTFTransform(t *testing.T) {
	data := []byte("nnbaaa")
	mtfTransform(symbolSet(data), data, data)

	if string(data) != "\x02\x00\x02\x02\x00\x00" {
		t.Error("Output is incorrect")
	}
}

func BenchmarkMTFTransform(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	src := make([]byte, 1000000)
	for i := range src {
		src[i] = byte(rand.Intn(256))
	}
	symbols := symbolSet(src)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst := make([]byte, len(src))
		mtfTransform(symbols, dst, src)
	}
}

func BenchmarkMTFTransformLarge(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	src := make([]byte, 1000000*6)
	for i := range src {
		src[i] = byte(rand.Intn(256))
	}
	symbols := symbolSet(src)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst := make([]byte, len(src))
		mtfTransform(symbols, dst, src)
	}
}
