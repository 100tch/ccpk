package cli

import (
	"flag"
	"fmt"
	"os"
)

const usageText = `ccpk - cc packer tool

Usage:
	ccpk <kind> [options] <input>

Kinds:
	sprite                 Sprite (exactly one source image)
	sprite-variations      SpriteVariations (sheet or array[Sprite])
	sprite-4way            Sprite4Way (frames split by -N/-E/-S/-W suffix)
	sprite-8way            Sprite8Way (frames split by -N/-NE/.../-NW suffix)
	animation              Animation (one "walked" spritesheet, requires -sheet)
	animation-variations   AnimationVariations (subdirs = variations)

Examples:
	ccpk sprite-variations -sheet -out ./output/ -name coal \
		-mod-path "__base__/graphics/icons/coal.png" \
		-lua ./output/coal.lua ./frames/


	ccpk sprite-4way -name pipe -trim -out ./output/ ./frames/


	ccpk animation -sheet -out ./output/ -name gman-walk \
		-anim-speed 0.5 ./frames/

Options:
`

func printUsage(fls *flag.FlagSet) {
	fmt.Fprint(os.Stderr, usageText)
	fls.SetOutput(os.Stderr)
	fls.PrintDefaults()
}
