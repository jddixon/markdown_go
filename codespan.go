package md

// xgo/md/codespan.go

// In Markdown input text a code span begins with one or two backticks
// and ends with the same number of backticks.  The delimiting backticks
// may themselves be separated by spaces.  If there is more than one
// opening backtick, a single backtick within the code segment is
// interpreted as a literal backtick.  The opening delimiting backtick(s)
// may be followed by a space; the closing delimiting backtick(s) may be
// preceded by a space.  These spaces will be dropped by the parser.

type CodeSpan struct {
	runes []rune
}

var (
	AMP_ENTITY = []rune("&amp;")
	LT_ENTITY  = []rune("&lt;")
	GT_ENTITY  = []rune("&gt;")
)

func NewCodeSpan(runes []rune) (t *CodeSpan) {
	txt := make([]rune, len(runes))
	copy(txt, runes)
	return &CodeSpan{runes: txt}
}

// Any backticks in CodeSpan.Runes are literal backticks.
// In the current implementation, < and > are automatically 'escaped',
// in the sense that they are converted to character entities here.
func (p *CodeSpan) Get() (out []rune) {

	out = append(out, []rune("<code>")...)
	for i := 0; i < len(p.runes); i++ {
		r := p.runes[i]
		if r == '&' {
			out = append(out, AMP_ENTITY...)
		} else if r == '<' {
			out = append(out, LT_ENTITY...)
		} else if r == '>' {
			out = append(out, GT_ENTITY...)
		} else {
			out = append(out, r)
		}
	}
	out = append(out, []rune("</code>")...)
	return
}

// Attempt to parse out a CodeSpan, returning a SpanI reference
// to it on success and nil and possibly an error on failure.  If the parse
// fails but there is no input error, leave the line offset unchanged and
// return a nil SpanI.  If the parse succeeds, return a SpanI and advance
// the offset accordingly.
//
func (q *Line) parseCodeSpan() (span SpanI, err error) {

	codeChar := q.runes[q.offset]
	offset := q.offset
	found := false

	// look for the end of the span
	for offset++; offset < len(q.runes); offset++ {
		ch := q.runes[offset]
		if ch == codeChar {
			found = true
			break
		}
	}
	if found {
		var start, end int
		start = q.offset + 1
		end = offset
		q.offset = offset + 1
		span = &CodeSpan{
			runes: q.runes[start:end],
		}
	}
	return
}
