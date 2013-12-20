package md

// xgo/md/definition.go

import (
	"fmt"
	u "unicode"
)

type Definition struct {
	uri   []rune
	title []rune
}

// Given a candidate ID in text, strip off leading and trailing spaces
// and then check that there are no spaces in the ID.  Return a valid
// ID in string form or an error.
func ValidID(text []rune) (validID string, err error) {
	id := make([]rune, len(text))
	copy(id, text)
	// get rid of any leading spaces
	for len(id) > 0 && u.IsSpace(id[0]) {
		id = id[1:]
	}
	if len(id) == 0 {
		err = NilID
	} else {
		// get rid of any trailing spaces
		for err == nil {
			if len(id) == 0 {
				err = EmptyID
			} else {
				ndxLast := len(id) - 1
				if u.IsSpace(id[ndxLast]) {
					id = id[:ndxLast]
				}
			}
		}
	}
	// this is a very loose definition of a valid ID!
	// XXX AND IT'S WRONG: SPACES WITHIN THE ID ARE OK
	if err == nil {
		for i := 0; i < len(id); i++ {
			if u.IsSpace(id[i]) {
				err = InvalidCharInID
			}
		}
	}
	if err == nil {
		validID = string(id)
	}
	return
}

// We are at the beginning of a line (possiblly with up to three leading
// spaces) ahd have seen a left square bracket.  If we find the rest of
//   [id]:\s+uri\s?("title")?(CR|LF)+
// where the uri may be delimited with angle brackets and the title
// may be delimited with DQUOTE or PAREN, then we absorb all of
// these, adding id => DEF to the dictionary for the document.  That
// is, a successful parse produces no output.
//
// If the parse fails, we push all characters scanned here back onto the input
// and returns collected == false.
//
func (p *Parser) parseDefinition() (collected bool, err error) {

	lx := p.lexer
	var (
		atEOF     bool
		cantParse bool
		id        string
	)
	// Enter having seen a left square bracket ('[') at the beginning
	// of a line, possibly preceded by up to three spaces.
	var runes []rune
	ch, err := lx.NextCh()
	for err == nil {
		runes = append(runes, ch)
		if ch == CR || ch == LF {
			cantParse = true
		} else if ch == ']' {
			id, err = ValidID(runes[:len(runes)-1])
			// DEBUG
			fmt.Printf("processIDRef: id is '%s'\n")
			// END
			break
		}
		ch, err = lx.NextCh()
	}
	// XXX STUB: skip spaces
	// XXX STUB: collect uri possibly enclosed in angle brackets
	// XXX STUB: skip spaces
	// XXX STUB: try to collect title delimited by DQUOTE or PAREN
	// XXX STUB: expect EOL, which will be absorbed if reference collected
	// XXX STUB: if no err, bulld Reference and add to dictionary
	// XXX STUB: if no err, collected = true

	_, _, _ = atEOF, cantParse, id

	// We push everything  back on the input and signal that the parse failed.
	// The caller knows that the current char is '['
	if !collected {
		lx.PushBackChars(runes)
	}
	return
}
