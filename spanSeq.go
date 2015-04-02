package md

// xgo/md/spanSeq.go

import (
	"fmt"
	"strings"
	u "unicode"
)

var _ = fmt.Print

type SpanSeq struct {
	spans   []SpanI
	lineSep []rune
}

func (sq *SpanSeq) String() string {
	var ss []string
	for i := 0; i < len(sq.spans); i++ {
		ss = append(ss, sq.spans[i].String())
	}
	s := strings.Join(ss, " ")
	return s + string(sq.lineSep)
}

// Advance down the line.  If a special character is encountered,
// invoke the parser(s) associated with the character (leaving the
// offset pointing to the character.  If a parse fails, it returns
// a nil SpanI and leaves the offset unchanged.  The special
// character is added to curText if all parses fail.  If a parse
// succeeds it returns a non-nil SpanI.  In such a case any curText
// is converted to a Text object and appended to the spans output
// slice, followed by the SpanI.
//
func (q *Line) parseSpanSeq(opt *Options, doc *Document, from uint,
	leftTrim bool) (seq *SpanSeq, err error) {

	var (
		curText          []rune
		testing, verbose bool
	)
	eol := uint(len(q.runes))
	if opt == nil {
		err = NilOptions
	} else {
		testing = opt.Testing
		verbose = opt.Verbose
		_ = verbose // DEBUG
		q.offset = from
		seq = new(SpanSeq)
		firstSpan := true
		for q.offset < eol {
			var span SpanI
			ch := q.runes[q.offset]

			// run through all candidate parsers ------------------------
			if ch == '_' || ch == '*' {
				span, _ = q.parseEmphSpan()
			} else if ch == '`' {
				span, _ = q.parseCodeSpan()
			} else if ch == '[' {
				span, _ = q.parseLinkSpan(opt)
				if span == nil {
					span, _ = q.parseLinkRefSpan(doc)
				}
			} else if ch == '!' {
				span, _ = q.parseImageSpan(opt)
				if span == nil {
					span, _ = q.parseImageRefSpan(doc)
				}
			} else if ch == '&' {
				span, _ = q.parseEntitySpan()
			} else if ch == '<' {
				span, _ = q.parseAutomaticLink()
			}

			// handle any parse results ---------------------------------
			if span == nil {
				if ch == '\\' && q.offset < eol-1 &&
					escaped(q.runes[q.offset+1]) {

					q.offset++
					ch = q.runes[q.offset]
				}
				curText = append(curText, ch)
				q.offset++
			} else {
				if len(curText) > 0 {
					if firstSpan && leftTrim {
						if testing {
							fmt.Println("LEFT-TRIMMING")
						}
						// get rid of any leading spaces
						for len(curText) > 0 {
							if u.IsSpace(curText[0]) {
								curText = curText[1:]
							} else {
								break
							}
						}
					}
					if len(curText) > 0 {
						seq.spans = append(seq.spans, NewTextSpan(curText))
						curText = curText[:0]
					}
				}
				seq.spans = append(seq.spans, span)
				firstSpan = false
			}
		} // end for loop

		if len(curText) > 0 {
			if firstSpan && leftTrim {
				if testing {
					fmt.Println("LEFT-TRIMMING")
				}
				// get rid of any leading spaces
				for len(curText) > 0 {
					if u.IsSpace(curText[0]) {
						curText = curText[1:]
					} else {
						break
					}
				}
			}
			txtLen := len(curText)
			if txtLen == 1 {
				seq.spans = append(seq.spans, NewTextSpan(curText))
			} else if txtLen == 2 {
				if curText[txtLen-1] == '\t' {
					// XXX DOESN'T CATCH space-tab
					seq.spans = append(seq.spans, NewTextSpan(curText[:2]))

				} else {
					seq.spans = append(seq.spans, NewTextSpan(curText))
					seq.spans = append(seq.spans, NewBreakSpan())
				}
			} else if txtLen > 2 {
				// convert trailing 2 spaces or tab to <br />
				matched := false
				if curText[txtLen-2] == ' ' && curText[txtLen-1] == ' ' {
					curText = curText[:txtLen-2]
					matched = true
				} else if curText[txtLen-1] == '\t' {
					curText = curText[:txtLen-1]
					matched = true
				}
				if matched {
					// drop any other trailing spaces too
					for len(curText) > 0 {
						ndxLast := len(curText) - 1
						if u.IsSpace(curText[ndxLast]) {
							curText = curText[:ndxLast]
						} else {
							break
						}
					}
				}
				if len(curText) > 0 {
					seq.spans = append(seq.spans, NewTextSpan(curText))
				}
				if matched {
					seq.spans = append(seq.spans, NewBreakSpan())
				}
			}
		}
		ls := q.lineSep
		for i := 0; i < len(ls); i++ {
			ch := ls[i]
			if ch != rune(0) {
				seq.lineSep = append(seq.lineSep, ch)
			}
		} // FOO
	}
	return
}
