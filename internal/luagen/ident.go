package luagen

import "strings"

func identifier(name string) string {
	var b strings.Builder
	for i := 0; i < len(name); i++ {
		c := name[i]

		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c == '_':
			b.WriteByte(c)
		case c >= '0' && c <= '9':
			if b.Len() == 0 {
				b.WriteByte('_')
			}
			b.WriteByte(c)
		default:
			b.WriteByte('_')
		}
	}

	if b.Len() == 0 {
		return "_"
	}
	return b.String()
}
