package md

// xgo/md/codeSpan.go

import (
	//"fmt"
)

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

func NewCodeSpan(runes []rune) (t *CodeSpan) {
	txt := make([]rune, len(runes))
	copy(txt, runes)
	return &CodeSpan{runes: txt}
}

func (p *CodeSpan) String() string {
	return "`" + string(p.runes) + "`"
}

// Any backticks in CodeSpan.Runes are literal backticks.
// In the current implementation, < and > are automatically 'escaped',
// in the sense that they are converted to character entities here.
func (p *CodeSpan) GetHtml() (out []rune) {

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

	const BACKTICK = '`'

	offset := q.offset
	EOL := uint(len(q.runes))
	found := false
	doubled := false

	// we require that the cursor is on the first BACKTICK
	offset++
	if offset < EOL && q.runes[offset] == BACKTICK {
		doubled = true
		offset++
	}

	// look for the end of the span
	for offset < EOL {
		ch := q.runes[offset]
		if ch == BACKTICK {
			if doubled {
				if offset < EOL-1 && q.runes[offset+1] == BACKTICK {
					offset++
					found = true
					break
				}
			} else {
				found = true
				break
			}
		} // FOO
		offset++
	}
	if found {
		var start, end uint
		start = q.offset + 1
		end = offset
		if doubled {
			start++
			end--
		}
		q.offset = offset + 1
		span = &CodeSpan{
			runes: q.runes[start:end],
		}
	}
	return
}
