package bzip2

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestRL2Encode(t *testing.T) {
	src := []byte("\x02\x00\x02\x02\x00\x00\x00\x00\x00")
	expected := []byte("\x03\x00\x03\x03\x00\x01\x04")

	dst := rl2Encode(symbolSet([]byte("banana")), src)
	if len(dst) != len(expected) {
		t.Error("RLE2 length doesn't match expected length")
	}
	for i, d := range dst {
		if d != expected[i] {
			t.Error("Byte value " + strconv.Itoa(int(d)) + " isn't the expected value " + strconv.Itoa(int(expected[i])))
		}
	}
}

func TestRL2EncodeShortRun(t *testing.T) {
	src := []byte("\x02\x00\x02\x02\x00\x00")
	expected := []byte("\x03\x00\x03\x03\x01\x04")

	dst := rl2Encode(symbolSet([]byte("banana")), src)
	if len(dst) != len(expected) {
		t.Error("RLE2 length doesn't match expected length")
	}
	for i, d := range dst {
		if d != expected[i] {
			t.Error("Byte value " + strconv.Itoa(int(d)) + " isn't the expected value " + strconv.Itoa(int(expected[i])))
		}
	}
}

func BenchmarkRL2Encode(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	src := make([]byte, 1000000)
	for i := range src {
		src[i] = byte(rand.Intn(256))
	}
	symbols := symbolSet(src)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl2Encode(symbols, src)
	}
}
