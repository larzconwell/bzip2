package rle

import (
	"testing"

	"github.com/larzconwell/bzip2/internal/testhelpers"
)

func TestRunListUpdate(t *testing.T) {
	runlist := NewRunList()

	expected := 6
	actual := runlist.Update([]byte("banana"))
	if actual != expected {
		t.Error("Update encoded length is not the expected length")
	}

	if runlist.EncodedLen() != expected {
		t.Error("EncodedLen returned an unexpected length")
	}

}

func TestRunListRunAcrossUpdate(t *testing.T) {
	runlist := NewRunList()
	expected := 20

	runlist.Update([]byte("banana"))
	actual := runlist.Update([]byte("aaabbbbbbanana"))
	if actual != expected {
		t.Error("Update encoded length is not the expected length")
	}

	if runlist.EncodedLen() != expected {
		t.Error("EncodedLen returned an unexpected length")
	}
}

func TestRunListUpdateRunLimit(t *testing.T) {
	runlist := NewRunList()
	data := make([]byte, 520)
	for i := range data {
		data[i] = 'b'
	}

	expected := 3
	runlist.Update(data)
	if len(runlist.runs) != expected {
		t.Error("Update should split runs larger than the max into multiple runs but isn't")
	}
}

func TestRunListTrim(t *testing.T) {
	runlist := NewRunList()
	runlist.Update([]byte("bananaaa"))

	expectedTrimmed := 2
	actualTrimmed := runlist.Trim(2)
	if actualTrimmed != expectedTrimmed {
		t.Error("Trim trimmed an unexpected number of actual bytes.")
	}

	expected := []byte("banana")
	actual := runlist.Encode()
	if len(actual) != len(expected) {
		t.Error("Trimmed encode data length doesn't match expected length")
	}

	for i, b := range actual {
		if b != expected[i] {
			t.Error("Byte value", string(b), "isn't the expected value",
				string(expected[i]))
		}
	}
}

func TestRunListTrimLong(t *testing.T) {
	runlist := NewRunList()
	runlist.Update([]byte("bananaaaaaaaaa"))

	expectedTrimmed := 6
	actualTrimmed := runlist.Trim(2)
	if actualTrimmed != expectedTrimmed {
		t.Error("Trim trimmed an unexpected number of actual bytes.")
	}

	expected := []byte("bananaaa")
	actual := runlist.Encode()
	if len(actual) != len(expected) {
		t.Error("Trimmed encode data length doesn't match expected length")
	}

	for i, b := range actual {
		if b != expected[i] {
			t.Error("Byte value", string(b), "isn't the expected value",
				string(expected[i]))
		}
	}
}

func TestRunListTrimAcrossRuns(t *testing.T) {
	runlist := NewRunList()
	runlist.Update([]byte("banannnaaaaa"))

	expectedTrimmed := 7
	actualTrimmed := runlist.Trim(7)
	if actualTrimmed != expectedTrimmed {
		t.Error("Trim trimmed an unexpected number of actual bytes.")
	}

	expected := []byte("banan")
	actual := runlist.Encode()
	if len(actual) != len(expected) {
		t.Error("Trimmed encode data length doesn't match expected length")
	}

	for i, b := range actual {
		if b != expected[i] {
			t.Error("Byte value", string(b), "isn't the expected value",
				string(expected[i]))
		}
	}
}

func TestRunListTrimLongLimit(t *testing.T) {
	runlist := NewRunList()
	runlist.Update([]byte("bananaaaa"))

	// This trim is a special case where it ends with an empty length
	// byte, we can't just trim the byte because that makes it invalid
	// so we remove the last actual byte too so it's a short run without
	// the length byte.
	expectedTrimmed := 1
	actualTrimmed := runlist.Trim(1)
	if actualTrimmed != expectedTrimmed {
		t.Error("Trim trimmed an unexpected number of actual bytes.")
	}

	expected := []byte("bananaaa")
	actual := runlist.Encode()
	if len(actual) != len(expected) {
		t.Error("Trimmed encode data length doesn't match expected length")
	}

	for i, b := range actual {
		if b != expected[i] {
			t.Error("Byte value", string(b), "isn't the expected value",
				string(expected[i]))
		}
	}
}

func TestRunListEncode(t *testing.T) {
	runlist := NewRunList()
	expected := []byte("bananaaaa\x00bbbb\x02anana")

	runlist.Update([]byte("banana"))
	runlist.Update([]byte("aaabbbbbbanana"))
	actual := runlist.Encode()
	if len(actual) != len(expected) {
		t.Error("Encode data length doesn't match expected length")
	}

	for i, b := range actual {
		if b != expected[i] {
			t.Error("Byte value", string(b), "isn't the expected value",
				string(expected[i]))
		}
	}
}

func BenchmarkRunListUpdate(b *testing.B) {
	data := testhelpers.RandomRunData(100000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runlist := NewRunList()
		runlist.Update(data)
	}
}

func BenchmarkRunListEncode(b *testing.B) {
	runlist := NewRunList()
	runlist.Update(testhelpers.RandomRunData(100000))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runlist.Encode()
	}
}
