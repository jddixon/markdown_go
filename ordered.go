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

func (h *Ordered) String() string {
	return "1. " + string(h.runes) + "\n"
}

func (h *Ordered) GetHtml() (r []rune) {
	r = append(r, LI_OPEN...)
	r = append(r, h.runes...)
	r = append(r, LI_CLOSE...)
	return
}

// Parse a line beginning with a digit followed by a dot followed by a space.
func (q *Line) parseOrdered(from uint) (b BlockI, err error) {

	var (
		bodyStart uint
		eol       uint = uint(len(q.runes))
		offset    uint
	)

	// Enter with the offset set to the first character after
	//   \s*\d+\.\s

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
		b, err = NewOrdered(q.runes[bodyStart:eol])
	}
	return
}
