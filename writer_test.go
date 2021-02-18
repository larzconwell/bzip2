package bzip2

import (
	"testing"
)

// TODO: Add Writer tests

func TestNewWriter(t *testing.T) {
	t.Parallel()
}

func TestNewWriterLevel(t *testing.T) {
	t.Run("with DefaultCompression level", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with level lower than BestSpeed", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with BestSpeed level", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with level within BestSpeed and BestCompression", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with BestCompression level", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with level greater than BestCompression", func(t *testing.T) {
		t.Parallel()
	})
}

func TestWriterErr(t *testing.T) {
	t.Run("no error has occurred", func(t *testing.T) {
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
}

func TestWriterWrite(t *testing.T) {
	t.Run("Writer has errored previously", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("Writer has been closed", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with an empty Writer less than a block of data is written", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with some data in the Writer enough data is written to fill a block of data", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with some data in the Writer enough data is written to fill a block of data with left over", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with an empty Writer enough data is written to fill a block of data", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with an empty Writer enough data is written to fill a block of data with left over", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with an empty Writer enough data is written to fill multiple blocks of data", func(t *testing.T) {
		t.Parallel()
	})
}
