package ui

import "strings"

const BigFontHeight = 5

var bigFont = map[rune][]string{
	'0': {
		"██████",
		"██  ██",
		"██  ██",
		"██  ██",
		"██████",
	},
	'1': {
		"  ██  ",
		" ███  ",
		"  ██  ",
		"  ██  ",
		"██████",
	},
	'2': {
		"██████",
		"    ██",
		"██████",
		"██    ",
		"██████",
	},
	'3': {
		"██████",
		"    ██",
		" █████",
		"    ██",
		"██████",
	},
	'4': {
		"██  ██",
		"██  ██",
		"██████",
		"    ██",
		"    ██",
	},
	'5': {
		"██████",
		"██    ",
		"██████",
		"    ██",
		"██████",
	},
	'6': {
		"██████",
		"██    ",
		"██████",
		"██  ██",
		"██████",
	},
	'7': {
		"██████",
		"    ██",
		"   ██ ",
		"  ██  ",
		"  ██  ",
	},
	'8': {
		"██████",
		"██  ██",
		"██████",
		"██  ██",
		"██████",
	},
	'9': {
		"██████",
		"██  ██",
		"██████",
		"    ██",
		"██████",
	},
	':': {
		"  ",
		"██",
		"  ",
		"██",
		"  ",
	},
	' ': {
		"      ",
		"      ",
		"      ",
		"      ",
		"      ",
	},
}

// BigText renders s as multi-line block art
func BigText(s string) string {
	rows := make([]string, BigFontHeight)
	for i, r := range []rune(s) {
		glyph, ok := bigFont[r]
		if !ok {
			glyph = bigFont[' ']
		}
		for row := range BigFontHeight {
			if i > 0 {
				rows[row] += " "
			}
			rows[row] += glyph[row]
		}
	}
	return strings.Join(rows, "\n")
}
