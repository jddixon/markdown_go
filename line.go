package md

// xgo/md/line.go

import ()

type Line struct {
	runes   []rune
	offset  int  // offset of current rune within this line
	lineSep rune // CR, LF, or 0
	doc     *Document
}

func NewLine(doc *Document, raw []rune, lineSep rune) (q *Line, err error) {
	if doc == nil {
		err = NilDocument
	} else if lineSep != rune(0) && lineSep != CR && lineSep != LF {
		err = InvalidLineSeparator
	} else {
		q = &Line{
			runes:   raw,
			lineSep: lineSep,
			doc:     doc,
		}
	}
	return
}
