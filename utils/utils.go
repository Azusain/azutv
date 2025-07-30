package utils

import (
	"strings"

	"github.com/rivo/uniseg"
)

func WrapEveryNRunes(s string, n int) string {
	runes := []rune(s)
	var b strings.Builder
	for i := 0; i < len(runes); i += n {
		end := i + n
		if end > len(runes) {
			end = len(runes)
		}
		b.WriteString(string(runes[i:end]))
		if end < len(runes) {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// wrap lines without breaking words.
func WrapAtWidth(s string, maxWidth int) string {
	state := -1
	remaining := s
	var buf, line strings.Builder
	width := 0

	for len(remaining) > 0 {
		seg, rest, mustBreak, newState := uniseg.FirstLineSegmentInString(remaining, state)
		remaining = rest
		state = newState

		segWidth := uniseg.StringWidth(seg)
		if width+segWidth > maxWidth {
			buf.WriteString(line.String())
			buf.WriteRune('\n')
			line.Reset()
			width = 0
		}

		line.WriteString(seg)
		width += segWidth

		if mustBreak {
			buf.WriteString(line.String())
			buf.WriteRune('\n')
			line.Reset()
			width = 0
		}
	}
	buf.WriteString(line.String())
	return buf.String()
}
