package md

// xgo/md/line.go

import ()

type Line struct {
	runes   []rune
	offset  uint   // offset of current rune within this line
	lineSep []rune // CR, LF, or both; or 0
	Err     error
}

// Used to pass a line for further processing.  The raw slice is not copied
// here, and so should be copied by using code.  lineSep _is_ copied here.
func NewLine(raw []rune, lineSep []rune) (q *Line) {
	var err error
	for i := 0; i < len(lineSep); i++ {
		sep := lineSep[i]
		if sep != rune(0) && sep != CR && sep != LF {
			err = InvalidLineSeparator
			break
		}
	}
	if err == nil {
		ls := make([]rune, len(lineSep))
		copy(ls, lineSep)
		q = &Line{
			runes:   raw,
			lineSep: ls,
		}
	} else {
		q = &Line{
			Err: err,
		}
	}
	return
}
