package bzip2

import (
	"bytes"
	"compress/bzip2"
	"io"
	"io/ioutil"
	"testing"

	"github.com/larzconwell/bzip2/internal/block"
	"github.com/larzconwell/bzip2/internal/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestNewWriter(t *testing.T) {
	t.Parallel()

	writer := NewWriter(ioutil.Discard)
	defer writer.Close()

	assert.NotNil(t, writer.bw)
	assert.NotNil(t, writer.block)
	assert.Equal(t, defaultLevel, writer.level)
	assert.Equal(t, uint32(0), writer.checksum)
	assert.Equal(t, false, writer.wroteHeader)
	assert.NoError(t, writer.err)
}

func TestNewWriterLevel(t *testing.T) {
	t.Run("sets the struct fields correctly", func(t *testing.T) {
		t.Parallel()

		writer, err := NewWriterLevel(ioutil.Discard, 1)
		assert.NoError(t, err)
		defer writer.Close()

		assert.NotNil(t, writer.bw)
		assert.NotNil(t, writer.block)
		assert.Equal(t, 1, writer.level)
		assert.Equal(t, uint32(0), writer.checksum)
		assert.Equal(t, false, writer.wroteHeader)
	})

	t.Run("with DefaultCompression level", func(t *testing.T) {
		t.Parallel()

		writer, err := NewWriterLevel(ioutil.Discard, DefaultCompression)
		assert.NoError(t, err)
		defer writer.Close()

		assert.Equal(t, defaultLevel, writer.level)
	})

	t.Run("with level lower than BestSpeed", func(t *testing.T) {
		t.Parallel()

		writer, err := NewWriterLevel(ioutil.Discard, BestSpeed-1)
		assert.Nil(t, writer)
		assert.ErrorIs(t, err, ErrInvalidCompressionLevel)
	})

	t.Run("with BestSpeed level", func(t *testing.T) {
		t.Parallel()

		writer, err := NewWriterLevel(ioutil.Discard, BestSpeed)
		assert.NoError(t, err)
		defer writer.Close()

		assert.Equal(t, BestSpeed, writer.level)
	})

	t.Run("with level within BestSpeed and BestCompression", func(t *testing.T) {
		t.Parallel()

		writer, err := NewWriterLevel(ioutil.Discard, 5)
		assert.NoError(t, err)
		defer writer.Close()

		assert.Equal(t, 5, writer.level)
	})

	t.Run("with BestCompression level", func(t *testing.T) {
		t.Parallel()

		writer, err := NewWriterLevel(ioutil.Discard, BestCompression)
		assert.NoError(t, err)
		defer writer.Close()

		assert.Equal(t, BestCompression, writer.level)
	})

	t.Run("with level greater than BestCompression", func(t *testing.T) {
		t.Parallel()

		writer, err := NewWriterLevel(ioutil.Discard, BestCompression+1)
		assert.Nil(t, writer)
		assert.ErrorIs(t, err, ErrInvalidCompressionLevel)
	})
}

func TestWriterErr(t *testing.T) {
	t.Run("no error has occurred", func(t *testing.T) {
		t.Parallel()

		writer := NewWriter(ioutil.Discard)
		writer.Write(make([]byte, 1))
		writer.Close()

		assert.NoError(t, writer.Err())
	})

	t.Run("Writer has been closed", func(t *testing.T) {
		t.Parallel()

		writer := NewWriter(ioutil.Discard)
		writer.Close()

		assert.NoError(t, writer.Err())
	})

	t.Run("error occurred during Write", func(t *testing.T) {
		t.Parallel()

		writer, err := NewWriterLevel(testhelpers.ErrReadWriter(io.ErrNoProgress), BestSpeed)
		assert.NoError(t, err)

		writer.Write(testhelpers.NoRunData(BestSpeed * block.BaseSize))
		assert.ErrorIs(t, writer.Err(), io.ErrNoProgress)
	})

	t.Run("error occurred during Close", func(t *testing.T) {
		t.Parallel()

		writer := NewWriter(testhelpers.ErrReadWriter(io.ErrNoProgress))

		writer.Write(make([]byte, 1))
		assert.NoError(t, writer.Err())

		writer.Close()
		assert.ErrorIs(t, writer.Err(), io.ErrNoProgress)
	})
}

func TestWriterReset(t *testing.T) {
	t.Parallel()

	writer := NewWriter(ioutil.Discard)
	_, err := writer.Write(make([]byte, 1))
	assert.NoError(t, err)
	writer.err = io.ErrNoProgress

	var out bytes.Buffer
	writer.Reset(&out)

	expected := testhelpers.NoRunData(10)

	_, err = writer.Write(expected)
	assert.NoError(t, err)
	err = writer.Close()
	assert.NoError(t, err)

	var actual bytes.Buffer
	reader := bzip2.NewReader(&out)
	_, err = io.Copy(&actual, reader)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual.Bytes())
}

func TestWriterClose(t *testing.T) {
	t.Run("Writer has errored previously", func(t *testing.T) {
		t.Parallel()

		writer := NewWriter(ioutil.Discard)
		writer.err = io.ErrNoProgress

		err := writer.Close()
		assert.ErrorIs(t, err, io.ErrNoProgress)
	})

	t.Run("Writer has been closed", func(t *testing.T) {
		t.Parallel()

		writer := NewWriter(ioutil.Discard)
		writer.Close()

		err := writer.Close()
		assert.ErrorIs(t, err, ErrClosed)
	})

	t.Run("with no data written", func(t *testing.T) {
		t.Parallel()

		var out bytes.Buffer
		writer := NewWriter(&out)
		err := writer.Close()
		assert.NoError(t, err)

		var actual bytes.Buffer
		reader := bzip2.NewReader(&out)
		_, err = io.Copy(&actual, reader)
		assert.NoError(t, err)

		assert.Equal(t, 0, actual.Len())
	})

	t.Run("with an unfilled block", func(t *testing.T) {
		t.Parallel()

		var out bytes.Buffer
		writer := NewWriter(&out)

		expected := testhelpers.NoRunData(10)
		_, err := writer.Write(expected)
		assert.NoError(t, err)

		err = writer.Close()
		assert.NoError(t, err)

		var actual bytes.Buffer
		reader := bzip2.NewReader(&out)
		_, err = io.Copy(&actual, reader)
		assert.NoError(t, err)

		assert.Equal(t, expected, actual.Bytes())
	})
}

func TestWriterWrite(t *testing.T) {
	t.Run("Writer has errored previously", func(t *testing.T) {
		t.Parallel()

		writer := NewWriter(ioutil.Discard)
		writer.err = io.ErrNoProgress

		_, err := writer.Write(make([]byte, 1))
		assert.ErrorIs(t, err, io.ErrNoProgress)
	})

	t.Run("Writer has been closed", func(t *testing.T) {
		t.Parallel()

		writer := NewWriter(ioutil.Discard)
		writer.Close()

		_, err := writer.Write(make([]byte, 1))
		assert.ErrorIs(t, err, ErrClosed)
	})

	t.Run("with an empty Writer less than a block of data is written", func(t *testing.T) {
		t.Parallel()

		var out bytes.Buffer
		writer, err := NewWriterLevel(&out, BestSpeed)
		assert.NoError(t, err)

		expected := testhelpers.NoRunData(10)
		_, err = writer.Write(expected)
		assert.NoError(t, err)

		err = writer.Close()
		assert.NoError(t, err)

		var actual bytes.Buffer
		reader := bzip2.NewReader(&out)
		_, err = io.Copy(&actual, reader)
		assert.NoError(t, err)

		assert.Equal(t, expected, actual.Bytes())
	})

	t.Run("with some data in the Writer enough data is written to fill a block of data", func(t *testing.T) {
		t.Parallel()

		var out bytes.Buffer
		writer, err := NewWriterLevel(&out, BestSpeed)
		assert.NoError(t, err)

		expected := testhelpers.NoRunData(BestSpeed * block.BaseSize)
		_, err = writer.Write(expected[:10])
		assert.NoError(t, err)
		_, err = writer.Write(expected[10:])
		assert.NoError(t, err)

		err = writer.Close()
		assert.NoError(t, err)

		var actual bytes.Buffer
		reader := bzip2.NewReader(&out)
		_, err = io.Copy(&actual, reader)
		assert.NoError(t, err)

		assert.Equal(t, expected, actual.Bytes())
	})

	t.Run("with some data in the Writer enough data is written to fill a block of data with left over", func(t *testing.T) {
		t.Parallel()

		var out bytes.Buffer
		writer, err := NewWriterLevel(&out, BestSpeed)
		assert.NoError(t, err)

		expected := testhelpers.NoRunData((BestSpeed * block.BaseSize) + 10)
		_, err = writer.Write(expected[:10])
		assert.NoError(t, err)
		_, err = writer.Write(expected[:10])
		assert.NoError(t, err)

		err = writer.Close()
		assert.NoError(t, err)

		var actual bytes.Buffer
		reader := bzip2.NewReader(&out)
		_, err = io.Copy(&actual, reader)
		assert.NoError(t, err)

		assert.Equal(t, expected, actual.Bytes())
	})

	t.Run("with an empty Writer enough data is written to fill a block of data", func(t *testing.T) {
		t.Parallel()

		var out bytes.Buffer
		writer, err := NewWriterLevel(&out, BestSpeed)
		assert.NoError(t, err)

		expected := testhelpers.NoRunData(BestSpeed * block.BaseSize)
		_, err = writer.Write(expected)
		assert.NoError(t, err)

		err = writer.Close()
		assert.NoError(t, err)

		var actual bytes.Buffer
		reader := bzip2.NewReader(&out)
		_, err = io.Copy(&actual, reader)
		assert.NoError(t, err)

		assert.Equal(t, expected, actual.Bytes())
	})

	t.Run("with an empty Writer enough data is written to fill a block of data with left over", func(t *testing.T) {
		t.Parallel()

		var out bytes.Buffer
		writer, err := NewWriterLevel(&out, BestSpeed)
		assert.NoError(t, err)

		expected := testhelpers.NoRunData((BestSpeed * block.BaseSize) + 10)
		_, err = writer.Write(expected)
		assert.NoError(t, err)

		err = writer.Close()
		assert.NoError(t, err)

		var actual bytes.Buffer
		reader := bzip2.NewReader(&out)
		_, err = io.Copy(&actual, reader)
		assert.NoError(t, err)

		assert.Equal(t, expected, actual.Bytes())
	})

	t.Run("with an empty Writer enough data is written to fill multiple blocks of data", func(t *testing.T) {
		t.Parallel()

		var out bytes.Buffer
		writer, err := NewWriterLevel(&out, BestSpeed)
		assert.NoError(t, err)

		expected := testhelpers.NoRunData(2 * block.BaseSize)
		_, err = writer.Write(expected)
		assert.NoError(t, err)

		err = writer.Close()
		assert.NoError(t, err)

		var actual bytes.Buffer
		reader := bzip2.NewReader(&out)
		_, err = io.Copy(&actual, reader)
		assert.NoError(t, err)

		assert.Equal(t, expected, actual.Bytes())
	})
}
