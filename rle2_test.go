package bzip2

import (
	"math/rand"
	"testing"
	"time"
)

func TestRL2Encode(t *testing.T) {
	src := []byte("\x02\x00\x02\x02\x00\x00\x00\x00\x00")
	expected := []byte("\x03\x00\x03\x03\x00\x01\x04")

	_, reduced := symbolSet([]byte("banana"))
	dst := rl2Encode(reduced, src)
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

	_, reduced := symbolSet(symbolSetData)
	dst := rl2Encode(reduced, src)
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

	_, reduced := symbolSet([]byte("banana"))
	dst := rl2Encode(reduced, src)
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
	_, reduced := symbolSet(src)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl2Encode(reduced, src)
	}
}

func TestSymbolFrequencies(t *testing.T) {
	data := []uint16{'\x03', '\x00', '\x03', '\x03', '\x00', '\x01', '\x04'}
	expectedFreq := []int{'\x00': 2, '\x01': 1, '\x03': 3, '\x04': 1}

	_, reduced := symbolSet([]byte("banana"))
	freq := symbolFrequencies(reduced, data)
	for i, f := range freq {
		if f != expectedFreq[i] {
			t.Error("Frequency", i, "isn't the expected value", f, "should be", expectedFreq[i])
		}
	}
}
