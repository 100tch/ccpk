package packer

import (
	"fmt"

	"ccpk/internal/frames"
	"ccpk/internal/luagen"
)

func packSprite(cfg Config) (Result, error) {
	paths, err := frames.ListFiles(cfg.InputDir, cfg.Pattern)
	if err != nil {
		return Result{}, err
	}
	if len(paths) != 1 {
		return Result{}, fmt.Errorf("-kind sprite expects exactly one matching image, found %d (use sprite-variations for multiple)", len(paths))
	}

	set, err := loadFrames(cfg, paths)
	if err != nil {
		return Result{}, err
	}
	sprite := set.sprite(cfg.luaPath(paths[0]))
	return set.result(luagen.Single(cfg.Name, sprite)), nil
}

func packSpriteVariations(cfg Config) (Result, error) {
	set, err := loadDir(cfg, cfg.InputDir)
	if err != nil {
		return Result{}, err
	}
	if cfg.UseSheet {
		return packAsVariationSheet(cfg, set)
	}

	sprites := make([]luagen.Sprite, len(set.frames))
	for i, frame := range set.frames {
		sprites[i] = set.sprite(cfg.luaPath(frame.Path))
	}
	return set.result(luagen.List(cfg.Name, sprites)), nil
}

func packDirectional(cfg Config, directions []frames.Direction) (Result, error) {
	all, err := frames.ListFiles(cfg.InputDir, cfg.Pattern)
	if err != nil {
		return Result{}, err
	}

	paths, err := frames.OrderByDirection(all, directions)
	if err != nil {
		return Result{}, err
	}

	set, err := loadFrames(cfg, paths)
	if err != nil {
		return Result{}, err
	}
	if cfg.UseSheet {
		return packAsVariationSheet(cfg, set)
	}

	fields := make([]luagen.Field, len(directions))
	for i, direction := range directions {
		fields[i] = luagen.Field{
			Key:    string(direction),
			Sprite: set.sprite(cfg.luaPath(paths[i])),
		}
	}
	return set.result(luagen.Keyed(cfg.Name, fields)), nil
}

func packAsVariationSheet(cfg Config, set frameSet) (Result, error) {
	sheet, err := writeSheet(cfg, cfg.Name, set.cells())
	if err != nil {
		return Result{}, err
	}

	sprite := set.sprite(cfg.luaPath(sheet.Path))
	sprite.LineLength = cfg.lineLength(sheet.Columns)
	sprite.VariationCount = len(set.frames)

	result := set.result(luagen.Sheet(cfg.Name, sprite))
	result.Sheet = &sheet
	return result, nil
}
