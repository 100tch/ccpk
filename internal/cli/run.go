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
		fmt.Fprintln(os.Stderr, "[error] us -help")
		return 1
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
}
