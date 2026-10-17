package packer

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"ccpk/internal/fsutil"
)

type Config struct {
	Kind    Kind
	Input   string
	Pattern string

	// trimming
	Trim           bool
	AlphaThreshold int
	Padding        int

	// spritesheet
	UseSheet  bool
	Columns   int
	OutputDir string

	// lua
	Name           string
	ModPath        string
	LineLength     int
	AnimationSpeed float64

	Force   bool
	Verbose bool
}

func (c Config) Validate() error {
	switch {
	case c.Kind == "":
		return fmt.Errorf("-kind is required (one of %s)", KindNames())
	case !slices.Contains(allKinds, c.Kind):
		return fmt.Errorf("unknown -kind %q (one of %s)", c.Kind, KindNames())
	case c.Input == "":
		return errors.New("input path is required")
	case !fsutil.Exists(c.Input):
		return fmt.Errorf("input %q does not exist", c.Input)
	case c.Kind != kindSprite && !fsutil.IsDir(c.Input):
		return fmt.Errorf("-kind %s needs a directory of frames, got a file: %q", c.Kind, c.Input)
	case c.Name == "" || c.Name == "." || c.Name == ".." || strings.ContainsAny(c.Name, `/\`):
		return fmt.Errorf("-name %q must be a plain name without directories", c.Name)
	case c.AlphaThreshold < 0 || c.AlphaThreshold > 255:
		return fmt.Errorf("-alpha must be between 0 and 255, got %d", c.AlphaThreshold)
	case c.Kind == kindAnimation && !c.UseSheet:
		return errors.New("-kind animation requires -sheet (a single spritesheet file walked by frame_count)")
	case c.UseSheet && c.OutputDir == "":
		return errors.New("-out is required when -sheet is set")
	case c.Kind == kindAnimationVariations && c.OutputDir == "":
		return errors.New("-out is required: each animation variation still needs its own spritesheet PNG")
	case c.Trim && c.OutputDir == "":
		return errors.New("-out is required with -trim")
	case c.outputInsideInput():
		return fmt.Errorf("-out %q must not be where the source files", c.OutputDir)
	case c.Padding != 0 && !c.Trim:
		return errors.New("-padding only applies together with -trim")
	case fsutil.HasPNGExt(c.ModPath) && c.multiFile():
		return errors.New("-mod-path is a single file path but this mode references several files")
	}
	return nil
}

func (c Config) luaPath(diskPath string) string {
	switch {
	case c.ModPath == "":
		return filepath.ToSlash(diskPath)
	case fsutil.HasPNGExt(c.ModPath):
		return c.ModPath
	}
	return path.Join(c.ModPath, filepath.Base(diskPath))
}

func (c Config) lineLength(sheetColumns int) int {
	if c.LineLength > 0 {
		return c.LineLength
	}
	return sheetColumns
}

func (c Config) outputInsideInput() bool {
	if c.OutputDir == "" {
		return false
	}
	if fsutil.IsDir(c.Input) {
		return fsutil.Inside(c.Input, c.OutputDir)
	}
	return fsutil.SamePath(filepath.Dir(c.Input), c.OutputDir)
}

func (c Config) multiFile() bool {
	return c.Kind != kindSprite && c.Kind != kindAnimation && !c.UseSheet
}

func DefaultConfig() Config {
	return Config{
		Pattern:        "*.png",
		AlphaThreshold: 25,
		Name:           "sprite",
	}
}
