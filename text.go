package md

// xgo/md/text.go

// A run of simple text.
type Text struct {
	runes []rune
}

func NewText(runes []rune) (t *Text) {
	txt := make([]rune, len(runes))
	copy(txt, runes)
	return &Text{runes: txt}
}

func (t *Text) Get() []rune {
	return t.runes
}
