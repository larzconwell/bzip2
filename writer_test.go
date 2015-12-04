package bzip2

import (
	"bytes"
	"compress/bzip2"
	"io"
	"io/ioutil"
	"testing"

	"github.com/larzconwell/bzip2/internal/testhelpers"
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
	expected := testhelpers.RandomRunData(baseBlockSize)

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
	expected := testhelpers.NoRunData(baseBlockSize)

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
	expected := testhelpers.RandomRunData(2 * baseBlockSize)

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
