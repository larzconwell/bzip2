package symbols

import (
	"testing"
)

func TestSymbolSet(t *testing.T) {
	symbols, reduced := Get([]byte("banana"))
	if string(reduced) != "abn" {
		t.Error("The reduced symbol set doesn't include the correct bytes")
	}

	for i, present := range symbols {
		if present == 0 {
			continue
		}

		switch i {
		case 'a', 'b', 'n':
		default:
			t.Error("Symbol set includes a byte that should be set")
		}
	}
}
