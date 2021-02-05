package mtf

import (
	"math/rand"
	"testing"
	"time"

	"github.com/larzconwell/bzip2/internal/symbols"
)

func TestMTFTransformEven(t *testing.T) {
	data := []byte("banana")
	_, reduced := symbols.Get(data)
	Transform(reduced, data, data)

	if string(data) != "\x01\x01\x02\x01\x01\x01" {
		t.Error("Output is incorrect")
	}
}

func TestMTFTransformOdd(t *testing.T) {
	data := []byte("baanana")
	_, reduced := symbols.Get(data)
	Transform(reduced, data, data)

	if string(data) != "\x01\x01\x00\x02\x01\x01\x01" {
		t.Error("Output is incorrect")
	}
}

func TestMTFTransformFullRange(t *testing.T) {
	data := make([]byte, 256)
	for i := range data {
		data[i] = byte(255 - i)
	}

	_, reduced := symbols.Get(data)
	Transform(reduced, data, data)
	if data[0] != '\xff' {
		t.Error("Output is incorrect")
	}
}

func TestMTFTransformAfterBWT(t *testing.T) {
	data := []byte("nnbaaa")
	_, reduced := symbols.Get(data)
	Transform(reduced, data, data)

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
	_, reduced := symbols.Get(src)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Transform(reduced, dst, src)
	}
}

func BenchmarkMTFTransformLarge(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	src := make([]byte, 1000000*6)
	dst := make([]byte, len(src))
	for i := range src {
		src[i] = byte(rand.Intn(256))
	}
	_, reduced := symbols.Get(src)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Transform(reduced, dst, src)
	}
}
