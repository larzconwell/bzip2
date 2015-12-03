package bzip2

import (
	"bytes"
	"compress/bzip2"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

func TestWriteAfterClose(t *testing.T) {
	var buf bytes.Buffer
	writer := NewWriter(&buf)
	err := writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	_, err = writer.Write([]byte{})
	if err != nil && err != ErrWriteAfterClose {
		t.Fatal(err)
	}

	if err == nil {
		t.Error("Write after closing should error but didn't")
	}

	// Minimally test reset.
	writer.Reset(&buf)
	_, err = writer.Write([]byte{})
	if err != nil && err != ErrWriteAfterClose {
		t.Fatal(err)
	}

	if err != nil {
		t.Error("Write after reset shouldn't return ErrWriterAfterClose but did")
	}
}

func TestEmptyValid(t *testing.T) {
	var buf bytes.Buffer
	writer := NewWriter(&buf)
	_, err := writer.Write([]byte{})
	if err == nil {
		err = writer.Close()
	}
	if err != nil {
		t.Fatal(err)
	}

	reader := bzip2.NewReader(&buf)
	_, err = io.Copy(ioutil.Discard, reader)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIncompleteBlock(t *testing.T) {
	var buf bytes.Buffer
	var out bytes.Buffer
	expected := []byte("banana")

	writer := NewWriter(&buf)
	_, err := writer.Write(expected)
	if err == nil {
		err = writer.Close()
	}
	if err != nil {
		t.Fatal(err)
	}

	reader := bzip2.NewReader(&buf)
	_, err = io.Copy(&out, reader)
	if err != nil {
		t.Fatal(err)
	}

	if out.String() != string(expected) {
		t.Error("Output is incorrect. Got", out.String(), "wanted",
			string(expected))
	}
}

func TestFilledBlock(t *testing.T) {
	var buf bytes.Buffer
	var out bytes.Buffer
	expected := randomRunData(baseBlockSize)

	writer := NewWriterLevel(&buf, 1)
	_, err := writer.Write(expected)
	if err == nil {
		err = writer.Close()
	}
	if err != nil {
		t.Fatal(err)
	}

	reader := bzip2.NewReader(&buf)
	_, err = io.Copy(&out, reader)
	if err != nil {
		t.Fatal(err)
	}

	if out.String() != string(expected) {
		t.Error("Output is incorrect.")
	}
}

func TestMultiBlock(t *testing.T) {
	var buf bytes.Buffer
	var out bytes.Buffer
	expected := randomRunData(2 * baseBlockSize)

	writer := NewWriterLevel(&buf, 1)
	_, err := writer.Write(expected)
	if err == nil {
		err = writer.Close()
	}
	if err != nil {
		t.Fatal(err)
	}

	reader := bzip2.NewReader(&buf)
	_, err = io.Copy(&out, reader)
	if err != nil {
		t.Fatal(err)
	}

	if out.String() != string(expected) {
		t.Error("Output is incorrect.")
	}
}

func TestFilledNoRunsBlock(t *testing.T) {
	var buf bytes.Buffer
	var out bytes.Buffer
	expected := noRunData(baseBlockSize)

	writer := NewWriterLevel(&buf, 1)
	_, err := writer.Write(expected)
	if err == nil {
		err = writer.Close()
	}
	if err != nil {
		t.Fatal(err)
	}

	reader := bzip2.NewReader(&buf)
	_, err = io.Copy(&out, reader)
	if err != nil {
		t.Fatal(err)
	}

	if out.String() != string(expected) {
		t.Error("Output is incorrect.")
	}
}

// noRunData produces data to write with no runs in it.
func noRunData(size int) []byte {
	data := make([]byte, size)
	b := byte('\x00')

	for i := range data {
		data[i] = b

		b++
	}

	return data
}

// randomRunData produces data to write with
// random bytes possibly with long runs.
func randomRunData(size int) []byte {
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, size)

	for i := range data {
		b := byte(rand.Intn(256))

		// Get the last bytes run length.
		runb := byte('\x00')
		runlen := 0
		if i > 0 {
			for j := i - 1; j >= 0; j-- {
				if j == i-1 {
					runb = data[j]
					runlen = 1
					continue
				}

				if data[j] != runb {
					break
				}

				runlen++
			}
		}

		// Detect if we shoud make the run longer.
		if shouldIncRun(runlen) {
			b = runb
		}

		data[i] = b
	}

	return data
}

// shouldIncRun decides if the run should be made longer.
func shouldIncRun(runlen int) bool {
	if runlen == 0 {
		return false
	}
	trues := 0
	falses := 0

	// Count runlen number of random trues/falses.
	for i := 0; i < runlen; i++ {
		if rand.Intn(2) == 1 {
			trues++
		} else {
			falses++
		}
	}

	if trues > falses {
		return true
	} else if trues < falses {
		return false
	}

	return rand.Intn(2) == 1
}
