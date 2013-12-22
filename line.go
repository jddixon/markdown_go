package md

import ()

type Line struct {
	runes   []rune
	offset  int
	lineSep rune // CR, LF, or 0
	p	*Parser
}

func NewLine(p *Parser, raw []rune, lineSep rune) (q *Line, err error) {
	if p == nil {
		err = NilParser
	} else if lineSep != rune(0) && lineSep != CR && lineSep != LF {
		err = InvalidLineSeparator
	} else {
		q = &Line{
			runes:   raw,
			lineSep: lineSep,
		}
	}
	return
}

// Advance down the line.  If a special character is encountered,
// invoke the parser(s) associated with the character (leaving the
// offset pointing to the character.  If a parse fails, it returns
// a nil SpanI and leaves the offset unchanged.  The special
// character is added to curText if all parses fail.  If a parse
// succeeds it returns a non-nil SpanI.  In such a case any curText
// is converted to a Text object and appended to the spans output
// slice, followed by the SpanI.

func (q *Line) parseToSpans() (spans []SpanI, err error) {

	var (
		curText []rune
	)

	for q.offset < len(q.runes) {
		var span SpanI
		ch := q.runes[q.offset]

		// run through all candidate parsers ------------------------
		if ch == '_' || ch == '*' {
			span, _ = q.parseEmphSpan()
		} else if ch == '`' {
			span, _ = q.parseCodeSpan()
		} else if ch == '[' {
			span, _ = q.parseLinkSpan()
			if span == nil {
				span, _ = q.parseLinkRefSpan(q.p)
			}
		}

		// handle any parse results ---------------------------------
		if span == nil {
			curText = append(curText, ch)
			q.offset++
		} else {
			if len(curText) > 0 {
				spans = append(spans, NewText(curText))
				curText = curText[:0]
			}
			spans = append(spans, span)
		}
	}
	if len(curText) > 0 {
		spans = append(spans, NewText(curText))
		// curText = curText[:0]
	}
	return
}
