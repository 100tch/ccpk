package packer

import (
	"fmt"
	"image"

	"ccpk/internal/frames"
	"ccpk/internal/luagen"
)

func packAnimation(cfg Config) (Result, error) {
	set, err := loadDir(cfg, cfg.InputDir)
	if err != nil {
		return Result{}, err
	}

	sheet, err := writeSheet(cfg, cfg.Name, set.cells())
	if err != nil {
		return Result{}, err
	}

	result := set.result(luagen.Single(cfg.Name, set.animation(cfg, sheet)))
	result.Sheet = &sheet
	return result, nil
}

func packAnimationVariations(cfg Config) (Result, error) {
	dirs, err := frames.ListSubdirs(cfg.InputDir)
	if err != nil {
		return Result{}, err
	}
	if len(dirs) == 0 {
		return Result{}, fmt.Errorf("no variation subdirectories found in %s", cfg.InputDir)
	}

	variations := make([]frameSet, len(dirs))
	for i, dir := range dirs {
		variations[i], err = loadDir(cfg, dir)
		if err != nil {
			return Result{}, fmt.Errorf("variation dir %s: %w", dir, err)
		}
	}

	if cfg.UseSheet {
		return packVariationsInOneSheet(cfg, variations)
	}
	return packVariationsInSeparateSheets(cfg, variations)
}

func packVariationsInOneSheet(cfg Config, variations []frameSet) (Result, error) {
	first := variations[0]

	var cells []*image.RGBA
	for _, variation := range variations {
		if variation.crop.Size() != first.crop.Size() {
			return Result{}, fmt.Errorf("all variations must share the same frame size for -sheet (got %dx%d and %dx%d)",
				first.width(), first.height(), variation.width(), variation.height())
		}
		cells = append(cells, variation.cells()...)
	}

	sheet, err := writeSheet(cfg, cfg.Name, cells)
	if err != nil {
		return Result{}, err
	}

	animation := first.animation(cfg, sheet)
	animation.VariationCount = len(variations)

	return Result{
		LuaSource:  luagen.Sheet(cfg.Name, animation),
		FrameCount: len(cells),
		Width:      first.width(),
		Height:     first.height(),
		Sheet:      &sheet,
	}, nil
}

func packVariationsInSeparateSheets(cfg Config, variations []frameSet) (Result, error) {
	animations := make([]luagen.Sprite, len(variations))
	frameCount := 0

	for i, variation := range variations {
		sheet, err := writeSheet(cfg, fmt.Sprintf("%s-%d", cfg.Name, i), variation.cells())
		if err != nil {
			return Result{}, err
		}
		animations[i] = variation.animation(cfg, sheet)
		frameCount += len(variation.frames)
	}

	return Result{
		LuaSource:  luagen.List(cfg.Name, animations),
		FrameCount: frameCount,
	}, nil
}

func (s frameSet) animation(cfg Config, sheet SheetFile) luagen.Sprite {
	animation := s.sprite(cfg.luaPath(sheet.Path))
	animation.FrameCount = len(s.frames)
	animation.LineLength = cfg.lineLength(sheet.Columns)
	animation.AnimationSpeed = cfg.AnimationSpeed
	return animation
}
