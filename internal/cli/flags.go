package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"ccpk/internal/packer"
)

type options struct {
	packer.Config
	LuaPath string
}

func parseArgs(args []string) (options, error) {
	opts := options{Config: packer.DefaultConfig()}

	fls := flag.NewFlagSet("ccpk", flag.ContinueOnError)
	fls.SetOutput(io.Discard)
	fls.Usage = func() { printUsage(fls) }
	defineFlags(fls, &opts)

	rest := args
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		if args[0] == "help" {
			fls.Usage()
			return opts, flag.ErrHelp
		}
		opts.Kind = packer.Kind(args[0])
		rest = args[1:]
	}

	positional, err := parseInterspersed(fls, rest)
	if err != nil {
		return opts, err
	}
	if opts.Kind == "" && len(positional) > 0 {
		return opts, fmt.Errorf("kind must be the frist argument")
	}

	// verbose config
	if opts.Verbose {
		explicit := make(map[string]bool)
		fls.Visit(func(f *flag.Flag) { explicit[f.Name] = true })

		tw := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "config:")
		fmt.Fprintf(tw, " kind\t%s\n", displayValue(string(opts.Kind)))
		fls.VisitAll(func(f *flag.Flag) {
			source := "default"
			if explicit[f.Name] {
				source = "set"
			}
			fmt.Fprintf(tw, " -%s\t%s\t(%s)\n", f.Name, displayValue(f.Value.String()), source)
		})
		tw.Flush()
	}

	switch len(positional) {
	case 0:
	case 1:
		opts.Input = positional[0]
	default:
		return opts, fmt.Errorf("expected one input path, got %d: %q", len(positional), positional)
	}

	return opts, opts.Validate()
}

func parseInterspersed(fls *flag.FlagSet, args []string) ([]string, error) {
	var positional []string
	for {
		if err := fls.Parse(args); err != nil {
			return nil, err
		}
		consumed := len(args) - fls.NArg()
		if consumed > 0 && args[consumed-1] == "--" {
			return append(positional, fls.Args()...), nil
		}

		args = fls.Args()
		if len(args) == 0 {
			return positional, nil
		}
		positional = append(positional, args[0])
		args = args[1:]
	}
}

func displayValue(v string) string {
	if v == "" {
		return `""`
	}
	return v
}

func defineFlags(fls *flag.FlagSet, opts *options) {
	fls.StringVar(&opts.Pattern, "pattern", opts.Pattern, "glob pattern for source files")
	fls.BoolVar(&opts.Trim, "trim", false, "trim frames to their alpha bounding box")
	fls.IntVar(&opts.AlphaThreshold, "alpha", opts.AlphaThreshold, "alpha value (0-255) at or above which a pixel counts as opaque for trimming")
	fls.IntVar(&opts.Padding, "padding", opts.Padding, "transparent padding (px) to keep around the trimmed bbox")
	fls.BoolVar(&opts.UseSheet, "sheet", false, "pack frames into one spritesheet PNG instead of an array of separate Sprite/Animation entries")
	fls.IntVar(&opts.Columns, "columns", opts.Columns, "columns in the packed sheet grid (0 = auto, ceil(sqrt(n)))")
	fls.StringVar(&opts.OutputDir, "out", "", "directory to output files")
	fls.StringVar(&opts.Name, "name", opts.Name, "base name: used for the output PNG file and as Lua variable name")
	fls.StringVar(&opts.ModPath, "mod-path", "", "file path to write into the Lua instead of the real one, e.g. __base__/graphics/icons/coal.png")
	fls.IntVar(&opts.LineLength, "line-length", opts.LineLength, "line_length written into the Lua (0 = auto)")
	fls.Float64Var(&opts.AnimationSpeed, "anim-speed", opts.AnimationSpeed, "animation_speed for animation/animation-variations (0 = omit, default 1)")
	fls.StringVar(&opts.LuaPath, "lua", "", "write the generated Lua to this file instead of stdout")
	fls.BoolVar(&opts.Force, "force", false, "overwrite existing output files")
	fls.BoolVar(&opts.Verbose, "v", false, "print debug info to stderr")
}
