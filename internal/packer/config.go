package packer

import (
	"ccpk/internal/fsutil"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"slices"
)

type Config struct {
	Kind     Kind
	InputDir string
	Pattern  string

	// trimming
	NoTrim         bool
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

	Verbose bool
}

func (c Config) Validate() error {
	switch {
	case c.Kind == "":
		return fmt.Errorf("-kind is required (one of %s)", KindNames())
	case !slices.Contains(allKinds, c.Kind):
		return fmt.Errorf("unknown -kind %q (one of %s)", c.Kind, KindNames())
	case c.InputDir == "":
		return errors.New("-in is required")
	case !fsutil.IsDir(c.InputDir):
		return fmt.Errorf("-in %q is not a directory", c.InputDir)
	case c.AlphaThreshold < 0 || c.AlphaThreshold > 255:
		return fmt.Errorf("-alpha must be between 0 and 255, got %d", c.AlphaThreshold)
	case c.Kind == kindAnimation && !c.UseSheet:
		return errors.New("-kind animation requires -sheet (a single spritesheet file walked by frame_count)")
	case c.UseSheet && c.OutputDir == "":
		return errors.New("-out is required when -sheet is set")
	case c.Kind == kindAnimationVariations && c.OutputDir == "":
		return errors.New("-out is required: each animation variation still needs its own spritesheet PNG")
	case !c.NoTrim && !c.writesPNG():
		return errors.New("trimming needs rewritten PNGs: use -sheet, or pass with -no-trim")
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

func (c Config) writesPNG() bool {
	switch c.Kind {
	case kindAnimation, kindAnimationVariations:
		return true
	case kindSprite:
		return false
	}
	return c.UseSheet
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
