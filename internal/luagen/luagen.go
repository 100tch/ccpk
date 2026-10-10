package luagen

import (
	"strconv"
)

type Sprite struct {
	Filename       string
	Width, Height  int
	FrameCount     int
	LineLength     int
	AnimationSpeed float64
	VariationCount int
}

type Field struct {
	Key    string
	Sprite Sprite
}

// name = { filename = ..., width = ..., height = ... }
func Single(name string, sprite Sprite) string {
	var w writer
	w.sprite(name, sprite)
	return w.String()
}

// name = { sheet = { filename = ..., ... } }
func Sheet(name string, sprite Sprite) string {
	var w writer
	w.open(name)
	w.sprite("sheet", sprite)
	w.close()
	return w.String()
}

// name = { { filename = ... }, { filename = ... } }
func List(name string, sprites []Sprite) string {
	var w writer
	w.open(name)
	for _, sprite := range sprites {
		w.sprite("", sprite)
	}
	w.close()
	return w.String()
}

// name = { north = { filename = ... }, east = { filename = ... } }
func Keyed(name string, fields []Field) string {
	var w writer
	w.open(name)
	for _, field := range fields {
		w.sprite(field.Key, field.Sprite)
	}
	w.close()
	return w.String()
}

func (s Sprite) fields() []string {
	fields := []string{
		"filename = " + strconv.Quote(s.Filename),
		"width = " + strconv.Itoa(s.Width),
		"height = " + strconv.Itoa(s.Height),
	}
	if s.FrameCount > 1 {
		fields = append(fields, "frame_count = "+strconv.Itoa(s.FrameCount))
	}
	if s.LineLength > 1 {
		fields = append(fields, "line_length = "+strconv.Itoa(s.LineLength))
	}
	if s.AnimationSpeed > 0 && s.AnimationSpeed != 1 {
		fields = append(fields, "animation_speed = "+strconv.FormatFloat(s.AnimationSpeed, 'f', -1, 64))
	}
	if s.VariationCount > 1 {
		fields = append(fields, "variation_count = "+strconv.Itoa(s.VariationCount))
	}
	return fields
}
