package bzip2

import (
	"flag"
)

var fuzzRetrySeed = flag.Int64("seed", 0, "seed to use running fuzz tests, use in conjunction with '-run TestFuzzWriter' or '-run TestFuzzReader'")
