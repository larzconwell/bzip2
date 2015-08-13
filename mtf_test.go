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

func TestMTFTransformEven(t *testing.T) {
	data := []byte("banana")
	mtfTransform(symbolSet(data), data, data)

	if string(data) != "\x01\x01\x02\x01\x01\x01" {
		t.Error("Output is incorrect")
	}
}

func TestMTFTransformOdd(t *testing.T) {
	data := []byte("baanana")
	mtfTransform(symbolSet(data), data, data)

	if string(data) != "\x01\x01\x00\x02\x01\x01\x01" {
		t.Error("Output is incorrect")
	}
}

func TestMTFTransformAfterBWT(t *testing.T) {
	data := []byte("nnbaaa")
	mtfTransform(symbolSet(data), data, data)

	if string(data) != "\x02\x00\x02\x02\x00\x00" {
		t.Error("Output is incorrect")
	}
}

func BenchmarkMTFTransform(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	src := make([]byte, 1000000)
	dst := make([]byte, len(src))
	for i := range src {
		src[i] = byte(rand.Intn(256))
	}
	symbols := symbolSet(src)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mtfTransform(symbols, dst, src)
	}
}

func BenchmarkMTFTransformLarge(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	src := make([]byte, 1000000*6)
	dst := make([]byte, len(src))
	for i := range src {
		src[i] = byte(rand.Intn(256))
	}
	symbols := symbolSet(src)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mtfTransform(symbols, dst, src)
	}
}
