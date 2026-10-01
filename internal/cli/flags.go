package cli

import (
	"flag"
	"io"

	"ccpk/internal/packer"
)

type options struct {
	packer.Config
	luaPath string
}

func parseArgs(args []string) (options, error) {
	opts := options{Config: packer.DefaultConfig()}

	fls := flag.NewFlagSet("ccpk", flag.ContinueOnError)
	fls.SetOutput(io.Discard)
	fls.Usage = func() { printUsage(fls) }
	defineFlags(fls, &opts)

	if err := fls.Parse(args); err != nil {
		return opts, err
	}
	return opts, opts.Validate()
}

func defineFlags(fls *flag.FlagSet, opts *options) {
}
