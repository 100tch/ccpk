package packer

import (
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"

	"ccpk/internal/frames"
	"ccpk/internal/imgproc"
	"ccpk/internal/luagen"
)

type Result struct {
	LuaSource     string
	FrameCount    int
	Width, Height int
	Sheet         *SheetFile
}

type SheetFile struct {
	Path    string
	Columns int
	Rows    int
}

func Pack(cfg Config) (Result, error) {
	switch cfg.Kind {
	case kindSprite:
		return packSprite(cfg)
	case kindSpriteVariations:
		return packSpriteVariations(cfg)
	case kindSprite4Way:
		return packDirectional(cfg, frames.Dir4)
	case kindSprite8Way:
		return packDirectional(cfg, frames.Dir8)
	case kindAnimation:
		return packAnimation(cfg)
	case kindAnimationVariations:
		return packAnimationVariations(cfg)
	}
	return Result{}, fmt.Errorf("unknown kind %q", cfg.Kind)
}

type frameSet struct {
	frames []imgproc.Frame
	crop   image.Rectangle
}

func loadFrames(cfg Config, paths []string) (frameSet, error) {
	loaded, err := imgproc.Load(paths)
	if err != nil {
		return frameSet{}, err
	}

	var crop image.Rectangle
	if cfg.NoTrim {
		crop, err = imgproc.FullRect(loaded)
	} else {
		crop, err = imgproc.TrimRect(loaded, cfg.AlphaThreshold, cfg.Padding, debugOutput(cfg))
	}
	if err != nil {
		return frameSet{}, err
	}
	return frameSet{frames: loaded, crop: crop}, nil
}

func loadDir(cfg Config, dir string) (frameSet, error) {
	paths, err := frames.ListFiles(dir, cfg.Pattern)
	if err != nil {
		return frameSet{}, err
	}
	return loadFrames(cfg, paths)
}

func debugOutput(cfg Config) io.Writer {
	if cfg.Verbose {
		return os.Stderr
	}
	return nil
}

func (s frameSet) width() int  { return s.crop.Dx() }
func (s frameSet) height() int { return s.crop.Dy() }
func (s frameSet) cells() []*image.RGBA {
	return imgproc.Crop(s.frames, s.crop)
}

func (s frameSet) sprite(filename string) luagen.Sprite {
	return luagen.Sprite{Filename: filename, Width: s.width(), Height: s.height()}
}

func (s frameSet) result(lua string) Result {
	return Result{
		LuaSource:  lua,
		FrameCount: len(s.frames),
		Width:      s.width(),
		Height:     s.height(),
	}
}

func writeSheet(cfg Config, name string, cells []*image.RGBA) (SheetFile, error) {
	sheet := imgproc.BuildSheet(cells, cfg.Columns)
	path := filepath.Join(cfg.OutputDir, name+".png")
	if err := imgproc.SavePNG(path, sheet.Image); err != nil {
		return SheetFile{}, fmt.Errorf("write sheet png: %w", err)
	}
	return SheetFile{Path: path, Columns: sheet.Columns, Rows: sheet.Rows}, nil
}
