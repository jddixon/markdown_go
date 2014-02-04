package md

// xgo/md/header.go

import (
	"fmt"
	u "unicode"
)

// This must implement BlockI
type Header struct {
	n     int
	runes []rune
}

// XXX MODIFY TO TRIM LEADING AND TRAILING SPACES
//
func NewHeader(n int, title []rune) (h BlockI, err error) {
	if n < 1 || 6 < n {
		err = HeaderNOutOfRange
	} else if len(title) == 0 {
		err = EmptyHeaderTitle
	} else {
		runes := make([]rune, len(title))
		copy(runes, title)
		h = &Header{
			n:     n,
			runes: runes,
		}
	}
	return
}

func (h *Header) GetHtml() []rune {
	text := fmt.Sprintf("<h%d>%s</h%d>\n", h.n, string(h.runes), h.n)
	return []rune(text)
}

// Collect atx-style headers preceded by 1-6 hash signs ('#') and optionally
// terminated by some number of hash signes.  If the parse fails, but there
// is no other error, return nil for both Header and error.  If the parse
// succeeds, return a non-nil Header.
//
func (q *Line) parseHeader(from uint) (b BlockI, forceNL bool, err error) {

	var (
		eol                  uint = uint(len(q.runes))
		hashCount            int
		offset               uint
		titleStart, titleEnd uint
	)

	// count leading hashes -----------------------------------------
	hashCount = 1 // we enter having seen one '#'

	// enter with the offset set to the first hash sign on the line
	for offset = from; offset < eol; offset++ {
		ch := q.runes[offset]
		if ch != '#' {
			break
		}
		hashCount++
	}
	// skip leading spaces ------------------------------------------
	for ; offset < eol; offset++ {
		ch := q.runes[offset]
		if !u.IsSpace(ch) {
			break
		}
	}

	// collect the title --------------------------------------------
	for titleStart = offset; offset < eol; offset++ {
		ch := q.runes[offset]
		if ch == '#' {
			titleEnd = offset
			break
		}
	}
	if titleEnd == 0 {
		titleEnd = offset
	}
	// if there is anything other than # at the end of the line, we
	// have a parse error
	for ; offset < eol; offset++ {
		ch := q.runes[offset]
		if ch != '#' {
			titleStart = 0
			break
		}
	}

	// if we have a title -------------------------------------------
	spaceCount := 0
	if titleStart > 0 && titleEnd > 0 {

		// drop any trailing spaces -----------------------
		for q.runes[titleEnd-1] == ' ' {
			titleEnd--
			spaceCount++ // count of spaces seen at end of line
		}
		if titleEnd > titleStart {
			if spaceCount >= 2 {
				forceNL = true
			}
			// create the object --------------------------
			b, _ = NewHeader(hashCount, q.runes[titleStart:titleEnd])
		}
	}
	return
}
