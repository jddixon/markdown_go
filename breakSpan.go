package md

// xgo/md/breakSpan.go

var BREAK_SPAN = []rune("<br />")

// A run of simple breakSpan.
type BreakSpan struct {
	runes []rune
}

func NewBreakSpan() (t *BreakSpan) {
	out := make([]rune, len(BREAK_SPAN))
	copy(out, BREAK_SPAN)
	return &BreakSpan{out}
}

// XXX This doesn't make much sense
func (br *BreakSpan) String() string {
	return string(br.runes)
}

func (br *BreakSpan) GetHtml() (out []rune) {
	return br.runes
}
