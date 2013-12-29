package md

// xgo/md/unordered.go

import (
	"fmt"
	u "unicode"
)

var _ = fmt.Print

// This must implement BlockI
type Unordered struct {
	runes []rune
}

func NewUnordered(body []rune) (h BlockI, err error) {
	if len(body) == 0 {
		err = EmptyUnorderedItem
	} else {
		runes := make([]rune, len(body))
		copy(runes, body)
		h = &Unordered{
			runes: runes,
		}
	}
	return
}

func (h *Unordered) Get() (r []rune) {
	r = append(r, LI_OPEN...)
	r = append(r, h.runes...)
	r = append(r, LI_CLOSE...)
	return
}

// Parse a line beginning with asterisk, plus, or minus (*, +, -) as an
// LI entity.
func (q *Line) parseUnordered(from int) (b BlockI, err error) {

	var (
		bodyStart int
		eol       int = len(q.runes)
		offset    int = from
	)

	// Enter with the offset set to the first character after
	//   \s*[\*\+\-]\s

	// skip leading spaces ------------------------------------------
	for offset = from; offset < eol; offset++ {
		ch := q.runes[offset]
		if !u.IsSpace(ch) {
			bodyStart = offset
			break
		}
	}

	// if we have a body -------------------------------------------
	if bodyStart > 0 {

		// drop any trailing spaces -----------------------
		for u.IsSpace(q.runes[eol-1]) {
			eol--
		}
		// create the object --------------------------
		b, err = NewUnordered(q.runes[bodyStart:eol])
	}
	return
}
