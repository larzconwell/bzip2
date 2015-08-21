package bzip2

import (
	"testing"
)

func TestHuffmanTreeCodeLength(t *testing.T) {
	_, reduced := symbolSet([]byte("banana"))
	data := []uint16{'\x03', '\x00', '\x03', '\x00', '\x01', '\x04'}
	freq := symbolFrequencies(reduced, data)

	lowestLen := 0
	huffmanTree := newHuffmanTree(freq)
	for i, code := range huffmanTree.Codes {
		if i == 0 || code.Len() < lowestLen {
			lowestLen = code.Len()
		}
	}

	if lowestLen != huffmanTree.Codes['\x03'].Len() {
		t.Error("The lowest code-length isn't the most used symbol")
	}
}

func TestGenerateHuffmanTreesLowest(t *testing.T) {
	_, reduced := symbolSet([]byte("banana"))
	data := []uint16{'\x03', '\x00', '\x03', '\x00', '\x01', '\x04'}
	freq := symbolFrequencies(reduced, data)

	huffmanTrees, treeIndexes := generateHuffmanTrees(freq, data)
	if len(huffmanTrees) < 2 {
		t.Error("Not enough huffman trees generated")
	}
	if len(huffmanTrees) > 6 {
		t.Error("Too many huffman trees generated")
	}

	if len(treeIndexes) != 1 {
		t.Error("The wrong number of huffman tree switches was returned")
	}
}

func TestGenerateHuffmanTreesMultipleSwitches(t *testing.T) {
	_, reduced := symbolSet([]byte("banana"))
	data := []uint16{'\x03', '\x00', '\x03', '\x00', '\x01', '\x04'}
	base := data
	for len(data) <= 200 {
		data = append(data, base...)
	}
	freq := symbolFrequencies(reduced, data)

	huffmanTrees, treeIndexes := generateHuffmanTrees(freq, data)
	if len(huffmanTrees) < 2 {
		t.Error("Not enough huffman trees generated")
	}
	if len(huffmanTrees) > 6 {
		t.Error("Too many huffman trees generated")
	}

	if len(treeIndexes) != 5 {
		t.Error("The wrong number of huffman tree switches was returned")
	}
}
