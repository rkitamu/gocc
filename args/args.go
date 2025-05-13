package args

import (
	"flag"
	"fmt"
)

type Args struct {
	Input  string
	Output string
	Debug  bool
}

// ParseArgs parses command line arguments and returns an Args struct.
// Returns an error if required arguments are missing.
func ParseArgs() (*Args, error) {
	input := flag.String("i", "", "Input file name")
	output := flag.String("o", "tmp.s", "Output file name")
	debug := flag.Bool("d", false, "Enable debug mode")

	flag.Parse()

	if *input == "" {
		return nil, fmt.Errorf("input file name is required")
	}

	args := &Args{
		Input: *input,
		Output: *output,
		Debug: *debug,
	}

	return args, nil;
}
