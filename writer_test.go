package bzip2

import (
	"bytes"
	"compress/bzip2"
	"io"
	"io/ioutil"
	"testing"
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
