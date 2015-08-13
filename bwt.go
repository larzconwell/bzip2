package bzip2

import (
	"bytes"
	"sort"
)

// rotateSort is a sort.Interface that sorts the rotations of the given data
// lexicographically.
type rotateSort struct {
	data    []byte
	rotates []int
	scrap   []byte
}

// newRotateSort creates a rotateSort generating the rotations.
func newRotateSort(data []byte) *rotateSort {
	rs := &rotateSort{
		data:    data,
		rotates: make([]int, len(data)),
		scrap:   make([]byte, len(data)),
	}
	for i := range rs.rotates {
		rs.rotates[i] = i
	}

	return rs
}

func (rs *rotateSort) Len() int {
	return len(rs.rotates)
}

func (rs *rotateSort) Less(i, j int) bool {
	var (
		a []byte
		b []byte
	)
	irotate := rs.rotates[i]
	jrotate := rs.rotates[j]

	if irotate > jrotate { // Data at i is short than j.
		b = rs.data[jrotate:]

		apref := rs.data[irotate:]
		copy(rs.scrap, apref)
		copy(rs.scrap[len(apref):], rs.data[:len(b)-len(apref)])
		a = rs.scrap[:len(b)]
	} else if irotate < jrotate { // Data at i is longer than j.
		a = rs.data[irotate:]

		bpref := rs.data[jrotate:]
		copy(rs.scrap, bpref)
		copy(rs.scrap[len(bpref):], rs.data[:len(a)-len(bpref)])
		b = rs.scrap[:len(a)]
	} else {
		a = rs.data[rs.rotates[i]:]
		b = rs.data[rs.rotates[j]:]
	}

	return bytes.Compare(a, b) == -1
}

func (rs *rotateSort) Swap(i, j int) {
	rs.rotates[i], rs.rotates[j] = rs.rotates[j], rs.rotates[i]
}

// bwTransform performs the Burrows-Wheeler Transform on the src slice and
// writes the results to dst, the index to the original src after sorting
// is returned.
func bwTransform(dst, src []byte) int {
	srclen := len(src)
	rs := newRotateSort(src)
	sort.Sort(rs)
	idx := -1

	for i, r := range rs.rotates {
		data := src[r:]
		datalen := len(data)

		// If it's the original input, set the index and the last character.
		if datalen == srclen {
			idx = i
			dst[i] = data[srclen-1]

			continue
		}

		// Get the last character in the suffix of the rotation.
		suffix := src[:srclen-datalen]
		dst[i] = suffix[srclen-datalen-1]
	}

	return idx
}
