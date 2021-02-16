package bzip2

import (
	"testing"
)

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

	t.Run("writer has been closed", func(t *testing.T) {
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
	t.Run("writer has errored previously", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("writer has been closed", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with no data written", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with an incomplete block", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with a completed block", func(t *testing.T) {
		t.Parallel()
	})
}

func TestWriterWrite(t *testing.T) {
	t.Run("writer has errored previously", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("writer has been closed", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with an empty writer less than a block of data is written", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with some data in the writer enough data is written to fill a block of data", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with some data in the writer enough data is written to fill a block of data with left over", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with an empty writer enough data is written to fill a block of data", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with an empty writer enough data is written to fill a block of data with left over", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("with an empty writer enough data is written to fill multiple blocks of data", func(t *testing.T) {
		t.Parallel()
	})
}
