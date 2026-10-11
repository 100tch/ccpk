package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"ccpk/internal/packer"
)

func Run(args []string) int {
	if len(args) == 0 {
		return fail(2, errors.New("no arguments, use -help"))
	}

	opts, err := parseArgs(args)
	if errors.Is(err, flag.ErrHelp) {
		return 0
	}
	if err != nil {
		return fail(2, err)
	}

	// packing
	result, err := packer.Pack(opts.Config)
	if err != nil {
		return fail(1, err)
	}
	// lua
	if err := writeLua(opts.LuaPath, result.LuaSource); err != nil {
		return fail(1, err)
	}

	// summary
	printSummary(opts, result)
	return 0
}

func fail(exitCode int, err error) int {
	fmt.Fprintln(os.Stderr, "[error]", err)
	return exitCode
}

func writeLua(path, lua string) error {
	if path == "" {
		_, err := fmt.Print(lua)
		return err
	}
	return os.WriteFile(path, []byte(lua), 0o644)
}

func printSummary(opts options, result packer.Result) {
	if !opts.Verbose {
		return
	}

	lua := "stdout"
	if opts.LuaPath != "" {
		lua = opts.LuaPath
	}

	fmt.Fprintln(os.Stderr, "result:")
	fmt.Fprintf(os.Stderr, " frames %d\n", result.FrameCount)
	if result.Width > 0 {
		fmt.Fprintf(os.Stderr, " size %dx%d px\n", result.Width, result.Height)
	}
	if result.Sheet != nil {
		fmt.Fprintf(os.Stderr, " sheet %s (%d columns x %d rows)\n", result.Sheet.Path, result.Sheet.Columns, result.Sheet.Rows)
	}
	fmt.Fprintf(os.Stderr, " lua %s\n", lua)
}
