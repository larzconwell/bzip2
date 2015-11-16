package rle2

import (
	"math/rand"
	"testing"
	"time"

	"github.com/larzconwell/bzip2/internal/symbols"
)

func TestRL2Encode(t *testing.T) {
	src := []byte("\x02\x00\x02\x02\x00\x00\x00\x00\x00")
	expected := []byte("\x03\x00\x03\x03\x00\x01\x04")

	_, reduced := symbols.Get([]byte("banana"))
	dst := Encode(reduced, src)
	if len(dst) != len(expected) {
		t.Error("RLE2 length doesn't match expected length")
	}

	for i, b := range dst {
		if b != uint16(expected[i]) {
			t.Error("Value", int(b), "isn't the expected value", int(expected[i]))
		}
	}
}

func TestRL2EncodeFullRange(t *testing.T) {
	symbolSetData := make([]byte, 256)
	for i := range symbolSetData {
		symbolSetData[i] = byte(i)
	}

	// src is the results of MTF where each byte is used in reverse order.
	src := make([]byte, 256)
	expected := make([]uint16, 257)
	for i := range src {
		src[i] = '\xff'
		expected[i] = '\u0100'
	}
	expected[len(expected)-1] = uint16(len(expected))

	_, reduced := symbols.Get(symbolSetData)
	dst := Encode(reduced, src)
	if len(dst) != len(expected) {
		t.Error("RLE2 length doesn't match expected length")
	}

	for i, b := range dst {
		if b != expected[i] {
			t.Error("Value", int(b), "isn't the expected value", int(expected[i]))
		}
	}
}

func TestRL2EncodeShortRun(t *testing.T) {
	src := []byte("\x02\x00\x02\x02\x00\x00")
	expected := []byte("\x03\x00\x03\x03\x01\x04")

	_, reduced := symbols.Get([]byte("banana"))
	dst := Encode(reduced, src)
	if len(dst) != len(expected) {
		t.Error("RLE2 length doesn't match expected length")
	}

	for i, b := range dst {
		if b != uint16(expected[i]) {
			t.Error("Value", int(b), "isn't the expected value", int(expected[i]))
		}
	}
}

func BenchmarkRL2Encode(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	src := make([]byte, 1000000)
	for i := range src {
		src[i] = byte(rand.Intn(256))
	}
	_, reduced := symbols.Get(src)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(reduced, src)
	}
}
