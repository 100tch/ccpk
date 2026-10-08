package luagen

import "strings"

type writer struct {
	text  strings.Builder
	depth int
}

func (w *writer) String() string { return w.text.String() }

func (w *writer) line(s string) {
	w.text.WriteString(strings.Repeat("  ", w.depth))
	w.text.WriteString(s)
	w.text.WriteByte('\n')
}

func (w *writer) open(key string) {
	if key != "" {
		w.line(key + " =")
	}
	w.line("{")
	w.depth++
}

func (w *writer) close() {
	w.depth--
	if w.depth == 0 {
		w.line("}")
	} else {
		w.line("},")
	}
}

func (w *writer) sprite(key string, sprite Sprite) {
	w.open(key)
	for _, field := range sprite.fields() {
		w.line(field + ",")
	}
	w.close()
}
