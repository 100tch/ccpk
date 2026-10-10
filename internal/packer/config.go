package packer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Kind string

var allKinds []Kind

func regKind(s string) Kind {
	k := Kind(s)
	allKinds = append(allKinds, k)
	return k
}

var (
	kindSprite              Kind = regKind("sprite")
	kindSpriteVariations    Kind = regKind("sprite-variations")
	kindSprite4Way          Kind = regKind("sprite-4way")
	kindSprite8Way          Kind = regKind("sprite-8way")
	kindAnimation           Kind = regKind("animation")
	kindAnimationVariations Kind = regKind("animation-variations")
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
	case !isDir(c.InputDir):
		return fmt.Errorf("-in %q is not a directory", c.InputDir)
	case c.AlphaThreshold < 0 || c.AlphaThreshold > 255:
		return fmt.Errorf("-alpha must be between 0 and 255, got %d", c.AlphaThreshold)
	case c.Kind == kindAnimation && !c.UseSheet:
		return errors.New("-kind animation requires -sheet (a single spritesheet file walked by frame_count)")
	case c.UseSheet && c.OutputDir == "":
		return errors.New("-out is required when -sheet is set")
	case c.Kind == kindAnimationVariations && c.OutputDir == "":
		return errors.New("-out is required: each animation variation still needs its own spritesheet PNG")
	}
	return nil
}

func (c Config) luaPath(diskPath string) string {
	if c.ModPath != "" {
		return c.ModPath
	}
	return filepath.ToSlash(diskPath)
}

func (c Config) lineLength(sheetColumns int) int {
	if c.LineLength > 0 {
		return c.LineLength
	}
	return sheetColumns
}

func KindNames() string {
	names := make([]string, len(allKinds))
	for i, kind := range allKinds {
		names[i] = string(kind)
	}
	return strings.Join(names, ",")
}

func DefaultConfig() Config {
	return Config{
		Pattern:        "*.png",
		AlphaThreshold: 25,
		Name:           "sprite",
	}
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
