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

func TestEmptyValid(t *testing.T) {
	var buf bytes.Buffer
	writer := NewWriter(&buf)
	_, err := writer.Write(make([]byte, 0))
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
	rand.Seed(time.Now().UnixNano())

	var buf bytes.Buffer
	var out bytes.Buffer
	expected := make([]byte, baseBlockSize)
	for i := range expected {
		expected[i] = byte(rand.Intn(256))
	}

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
	rand.Seed(time.Now().UnixNano())

	var buf bytes.Buffer
	var out bytes.Buffer
	expected := make([]byte, 2*baseBlockSize)
	for i := range expected {
		expected[i] = byte(rand.Intn(256))
	}

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
