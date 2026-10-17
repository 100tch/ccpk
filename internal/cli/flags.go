package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
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

	if err := fls.Parse(args); err != nil {
		return opts, err
	}

	// verbose config
	if opts.Verbose {
		explicit := make(map[string]bool)
		fls.Visit(func(f *flag.Flag) { explicit[f.Name] = true })

		tw := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "config:")
		fls.VisitAll(func(f *flag.Flag) {
			source := "default"
			if explicit[f.Name] {
				source = "set"
			}
			fmt.Fprintf(tw, " -%s\t%s\t(%s)\n", f.Name, displayValue(f.Value.String()), source)
		})
		tw.Flush()
	}

	switch fls.NArg() {
	case 0:
	case 1:
		opts.Input = fls.Arg(0)
	default:
		return opts, fmt.Errorf("expected one input path, got %d: %q (all options must start with \"-\" and come before path)", fls.NArg(), fls.Args())
	}

	return opts, opts.Validate()
}

func displayValue(v string) string {
	if v == "" {
		return `""`
	}
	return v
}

func defineFlags(fls *flag.FlagSet, opts *options) {
	fls.StringVar((*string)(&opts.Kind), "kind", "", "output shape: one of"+packer.KindNames())
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
	fls.Float64Var(&opts.AnimationSpeed, "anim-speed", opts.AnimationSpeed, "animation_speed for -kind animation/animation-variations (0 = omit, default 1)")
	fls.StringVar(&opts.LuaPath, "lua", "", "write the generated Lua to this file instead of stdout")
	fls.BoolVar(&opts.Verbose, "v", false, "print debug info to stderr")
}
