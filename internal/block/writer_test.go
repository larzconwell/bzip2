package block

import (
	"testing"
)

// TODO: Add Writer tests

func TestNewWriter(t *testing.T) {
	t.Parallel()
}

func TestWriterErr(t *testing.T) {
	t.Run("no error has occured", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("Writer has been closed", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("error occurred during Write", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("error occurred during Close", func(t *testing.T) {
		t.Parallel()
	})
}

func TestWriterReset(t *testing.T) {
	t.Parallel()
}

func TestWriterLen(t *testing.T) {
	t.Parallel()
}

func TestWriterWrite(t *testing.T) {
	t.Run("Writer has errored previously", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("Writer has been closed", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with less than the block size limit of data", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with exactly the block size limit of data", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with more than the block size limit of data", func(t *testing.T) {
		t.Parallel()
	})
}

func TestWriterClose(t *testing.T) {
	t.Run("Writer has errored previously", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("Writer has been closed", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with no data written", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with an unfilled block", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with a filled block", func(t *testing.T) {
		t.Parallel()
	})
}
