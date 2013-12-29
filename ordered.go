package md

// xgo/md/ordered.go

import (
	"fmt"
	u "unicode"
)

var _ = fmt.Print

// This must implement BlockI
type Ordered struct {
	runes []rune
}

func NewOrdered(body []rune) (h BlockI, err error) {
	if len(body) == 0 {
		err = EmptyOrderedItem
	} else {
		runes := make([]rune, len(body))
		copy(runes, body)
		h = &Ordered{
			runes: runes,
		}
	}
	return
}

func (h *Ordered) Get() (r []rune) {
	r = append(r, LI_OPEN...)
	r = append(r, h.runes...)
	r = append(r, LI_CLOSE...)
	return
}

// XXX JUST COPIED FROM UNORDERED: TRUST ME NOT!

// Parse a line beginning with a digit followed by a dot followed by a space.
func (q *Line) parseOrdered() (b BlockI, err error) {

	var (
		bodyStart int
		eol       int = len(q.runes)
		offset    int
	)

	// skip leading spaces ------------------------------------------
	for offset = 1; offset < eol; offset++ {
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
		b, err = NewOrdered(q.runes[bodyStart:eol])
	}
	return
}
