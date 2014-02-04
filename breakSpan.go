package md

// xgo/md/breakSpan.go

var BREAK_SPAN = []rune("<br />")

// A run of simple breakSpan.
type BreakSpan struct {
	runes *[]rune
}

func NewBreakSpan() (t *BreakSpan) {
	out := make([]rune, len(BREAK_SPAN))
	copy(out, BREAK_SPAN)
	return &BreakSpan{&out}
}

func (br *BreakSpan) GetHtml() (out []rune) {
	return *br.runes
}
