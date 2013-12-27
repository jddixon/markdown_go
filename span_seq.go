package md

// xgo/md/span_seq.go

import ()

type SpanSeq struct {
	spans   []SpanI
	lineSep []rune
}

// Advance down the line.  If a special character is encountered,
// invoke the parser(s) associated with the character (leaving the
// offset pointing to the character.  If a parse fails, it returns
// a nil SpanI and leaves the offset unchanged.  The special
// character is added to curText if all parses fail.  If a parse
// succeeds it returns a non-nil SpanI.  In such a case any curText
// is converted to a Text object and appended to the spans output
// slice, followed by the SpanI.

func (q *Line) parseSpanSeq() (seq *SpanSeq, err error) {

	var (
		curText []rune
	)
	seq = new(SpanSeq)
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
				span, _ = q.parseLinkRefSpan(q.doc)
			}
		} else if ch == '!' {
			span, _ = q.parseImageSpan()
			if span == nil {
				span, _ = q.parseImageRefSpan(q.doc)
			}
		}

		// handle any parse results ---------------------------------
		if span == nil {
			curText = append(curText, ch)
			q.offset++
		} else {
			if len(curText) > 0 {
				seq.spans = append(seq.spans, NewText(curText))
				curText = curText[:0]
			}
			seq.spans = append(seq.spans, span)
		}
	}
	if len(curText) > 0 {
		seq.spans = append(seq.spans, NewText(curText))
	}
	ls := q.lineSep
	for i := 0; i < len(ls); i++ {
		ch := ls[i]
		if ch != rune(0) {
			seq.lineSep = append(seq.lineSep, ch)
		}
	}
	return
}
