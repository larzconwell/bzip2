package rle

import (
	"testing"
)

func TestIndexOf(t *testing.T) {
	src := []byte("sk\x02\x02\x02\x02\x02\x02\x02")

	dst := Encode(src)
	idx := IndexOf(len(dst)-1, dst)
	if idx != len(src)-1 {
		t.Error("Index value is incorrect. Got", idx, "wanted", len(src)-1)
	}
}

func TestIndexOfComplex(t *testing.T) {
	src := []byte("sk\x02\x02\x02\x02\x02\x02vrrrrrkc")

	dst := Encode(src)
	idx := IndexOf(len(dst)-2, dst)
	if idx != len(src)-2 {
		t.Error("Index value is incorrect. Got", idx, "wanted", len(src)-2)
	}
}
