package md

// xgo/md/textSpan.go

// A run of simple text.
type TextSpan struct {
	runes []rune
}

func NewTextSpan(runes []rune) (t *TextSpan) {
	var txt []rune

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '&' {
			txt = append(txt, AMP_ENTITY...)
		} else if r == HARD_SPACE {
			txt = append(txt, rune(SPACE))
		} else {
			txt = append(txt, r)
		}
	}
	return &TextSpan{runes: txt}
}

func (t *TextSpan) String() string {
	return string(t.runes)
}
func (t *TextSpan) GetHtml() []rune {
	return t.runes
}
