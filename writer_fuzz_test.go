package bzip2

import (
	"bytes"
	"compress/bzip2"
	"context"
	"errors"
	"io"
	"math/rand"
	"os/exec"
	"testing"
	"time"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
)

func TestFuzzWriter(t *testing.T) {
	maxSize := 1_000_000
	iters := 1
	seed := *fuzzRetrySeed

	if seed == 0 {
		iters = 50
		seed = time.Now().UnixNano()
	}

	t.Run("reading with the bzip2 binary", func(t *testing.T) {
		t.Parallel()
		bin := "bzip2"

		_, err := exec.LookPath(bin)
		if errors.Is(err, exec.ErrNotFound) {
			t.Skip("Skipping bzip2 binary fuzzer test since the bzip2 binary cannot be found")
		}

		source := rand.NewSource(seed)
		rng := rand.New(source)
		fuzzer := fuzz.New().RandSource(source).NilChance(0).NumElements(1, maxSize)

		for i := 0; i < iters; i++ {
			raw, compressed, err := fuzzWriter(rng, fuzzer)
			assert.NoError(t, err)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			var actual bytes.Buffer
			cmd := exec.CommandContext(ctx, bin, "-d", "-c", "-q")
			cmd.Stdin = bytes.NewReader(compressed)
			cmd.Stdout = &actual

			err = cmd.Run()
			assert.NoError(t, err)

			if !assert.Equal(t, raw, actual.Bytes()) {
				t.Logf("retry with '-seed %#v -run TestFuzzWriter'", seed)
			}
		}
	})

	t.Run("reading with the Go stdlib bzip2 package", func(t *testing.T) {
		t.Parallel()

		source := rand.NewSource(seed)
		rng := rand.New(source)
		fuzzer := fuzz.New().RandSource(source).NilChance(0).NumElements(1, maxSize)

		for i := 0; i < iters; i++ {
			raw, compressed, err := fuzzWriter(rng, fuzzer)
			assert.NoError(t, err)

			var actual bytes.Buffer
			reader := bzip2.NewReader(bytes.NewReader(compressed))
			_, err = io.Copy(&actual, reader)
			assert.NoError(t, err)

			if !assert.Equal(t, raw, actual.Bytes()) {
				t.Logf("retry with '-seed %#v -run TestFuzzWriter'", seed)
			}
		}
	})

	t.Run("reading with Reader", func(t *testing.T) {
		t.Parallel()
	})
}

func fuzzWriter(rng *rand.Rand, fuzzer *fuzz.Fuzzer) ([]byte, []byte, error) {
	var compressed bytes.Buffer
	writer := NewWriter(&compressed)

	var raw []byte
	fuzzer.Fuzz(&raw)

	data := raw
	for len(data) != 0 {
		count := rng.Intn(len(data)) + 1
		batch := data[0:count]
		data = data[count:]

		_, err := writer.Write(batch)
		if err != nil {
			return nil, nil, err
		}
	}

	return raw, compressed.Bytes(), nil
}
