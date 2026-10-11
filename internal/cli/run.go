package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"ccpk/internal/packer"
)

// return error code
func Run(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "[error] ", errors.New("no arguments, use -help"))
		return 2
	}

	opts, err := parseArgs(args)
	if errors.Is(err, flag.ErrHelp) {
		return 0
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "[error] ", err)
		return 2
	}

	// packing
	result, err := packer.Pack(opts.Config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[error] ", err)
		return 1
	}
	// lua
	if err := writeLua(opts.LuaPath, result.LuaSource); err != nil {
		fmt.Fprintln(os.Stderr, "[error] ", err)
		return 1
	}

	// summary
	printSummary(opts, result)
	return 0
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
