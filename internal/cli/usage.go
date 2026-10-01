package cli

import (
	"flag"
	"fmt"
	"os"
)

const usageText = `cc - cc packer

Usage:
	ccpk -kind <kind> -in <dir> [options]

Kinds:
	sprite                 Sprite (exactly one source image)
	sprite-variations      SpriteVariations (sheet or array[Sprite])
	sprite-4way            Sprite4Way (frames split by -N/-E/-S/-W suffix)
	sprite-8way            Sprite8Way (frames split by -N/-NE/.../-NW suffix)
	animation              Animation (one "walked" spritesheet, requires -sheet)
	animation-variations   AnimationVariations (subdirs = variations)

Examples:
	TODO

Options:
`

func printUsage(fls *flag.FlagSet) {
	fmt.Fprint(os.Stderr, usageText)
	fls.SetOutput(os.Stderr)
	fls.PrintDefaults()
}
