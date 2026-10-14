package packer

import "strings"

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

func KindNames() string {
	names := make([]string, len(allKinds))
	for i, kind := range allKinds {
		names[i] = string(kind)
	}
	return strings.Join(names, ", ")
}
