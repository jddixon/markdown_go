package md

// xgo/md/text.go

// A run of simple text.
type Text struct {
	runes []rune
}

func NewText(runes []rune) (t *Text) {
	txt := make([]rune, len(runes))
	// copy(txt, runes)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == rune(0x00a0) {
			r = rune(0x20)
		}
		txt[i] = r
	}
	return &Text{runes: txt}
}

func (t *Text) Get() []rune {
	return t.runes
}
