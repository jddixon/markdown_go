package md

// xgo/md/textSpan.go

// A run of simple text.
type TextSpan struct {
	runes []rune
}

func NewTextSpan(runes []rune) (t *TextSpan) {
	txt := make([]rune, len(runes))
	// copy(txt, runes)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == rune(0x00a0) {
			r = rune(0x20)
		}
		txt[i] = r
	}
	return &TextSpan{runes: txt}
}

func (t *TextSpan) Get() []rune {
	return t.runes
}
