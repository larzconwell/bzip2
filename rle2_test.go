package bzip2

import (
	"math/rand"
	"testing"
	"time"
)

func TestRL2Encode(t *testing.T) {
	src := []byte("\x02\x00\x02\x02\x00\x00\x00\x00\x00")
	expected := []byte("\x03\x00\x03\x03\x00\x01\x04")
	expectedFreq := []int{'\x00': 2, '\x01': 1, '\x03': 3, '\x04': 1}

	freq, dst := rl2Encode(symbolSet([]byte("banana")), src)
	if len(dst) != len(expected) {
		t.Error("RLE2 length doesn't match expected length")
	}

	for i, d := range dst {
		if d != uint16(expected[i]) {
			t.Error("Value", int(d), "isn't the expected value", int(expected[i]))
		}
	}

	for i, f := range freq {
		if f != expectedFreq[i] {
			t.Error("Frequency", i, "isn't the expected value", f, "should be", expectedFreq[i])
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
	expectedFreq := []int{'\u0100': 256, '\u0101': 1}

	freq, dst := rl2Encode(symbolSet(symbolSetData), src)
	if len(dst) != len(expected) {
		t.Error("RLE2 length doesn't match expected length")
	}

	for i, d := range dst {
		if d != expected[i] {
			t.Error("Value", int(d), "isn't the expected value", int(expected[i]))
		}
	}

	for i, f := range freq {
		if f != expectedFreq[i] {
			t.Error("Frequency", i, "isn't the expected value", f, "should be", expectedFreq[i])
		}
	}
}

func TestRL2EncodeShortRun(t *testing.T) {
	src := []byte("\x02\x00\x02\x02\x00\x00")
	expected := []byte("\x03\x00\x03\x03\x01\x04")
	expectedFreq := []int{'\x00': 1, '\x01': 1, '\x03': 3, '\x04': 1}

	freq, dst := rl2Encode(symbolSet([]byte("banana")), src)
	if len(dst) != len(expected) {
		t.Error("RLE2 length doesn't match expected length")
	}

	for i, d := range dst {
		if d != uint16(expected[i]) {
			t.Error("Value", int(d), "isn't the expected value", int(expected[i]))
		}
	}

	for i, f := range freq {
		if f != expectedFreq[i] {
			t.Error("Frequency", i, "isn't the expected value", f, "should be", expectedFreq[i])
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
