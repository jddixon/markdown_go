package md

// xgo/md/definition.go

import (
	"fmt"
	u "unicode"
)

var _ = fmt.Print

// We use the same data structure for both link and image defs.
type Definition struct {
	uri   []rune
	title []rune
}

func (def *Definition) GetURI() string {
	return string(def.uri)
}

func (def *Definition) GetTitle() string {
	return string(def.title)
}

// XXX THIS IS CURRENTLY NOT USED

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
// spaces) and have seen a left square bracket.  If we find the rest of
//   [id]:\s+uri\s?("title")
// where the uri may be delimited with angle brackets and the title
// may be delimited with DQUOTE or PAREN, then we absorb all of
// these, adding id => DEF to the dictionary for the document.  That
// is, a successful parse produces no output.
//
// If there is any deviation from the spec, we leave the offset where it
// is and return a nil definition.  If the parse succeeds, we add the
// definition to the parser's dictionary, set the offset, and return a
// non-nil definition.
//
func (line *Line) parseLinkDefinition(p *Parser) (def *Definition, err error) {

	var (
		ch                   rune
		idStart, idEnd       int
		offset               int
		uriStart, uriEnd     int
		titleStart, titleEnd int
	)
	// Enter having seen a left square bracket ('[') at the beginning
	// of a line, possibly preceded by up to three spaces.  The offset
	// is on the bracket.

	eol := len(line.runes)
	offset = line.offset + 1 // just beyond the bracket

	// collect the id -----------------------------------------------
	for idStart = offset; offset < eol; offset++ {
		ch = line.runes[offset]
		if ch == ']' {
			idEnd = offset // exclusive end
			offset++       // position beyond right bracket
			break
		}
	}
	// expect a colon and one or more spaces ------------------------
	if idEnd > 0 && offset < eol-3 {
		if line.runes[offset] == ':' {
			offset++
			// skip any spaces
			for ch = line.runes[offset]; offset < eol && u.IsSpace(ch); ch = line.runes[offset] {

				offset++
			}
			if offset < eol {
				uriStart = offset
			}
		}
	}
	// collect the uri ----------------------------------------------
	if uriStart > 0 {
		// assume that a uri contains no spaces
		for offset < eol && !u.IsSpace(line.runes[offset]) {
			offset++
		}
		uriEnd = offset
	}
	// collect any title
	if uriEnd > 0 && offset < eol {
		// skip any spaces
		for ch = line.runes[offset]; offset < eol && u.IsSpace(ch); ch = line.runes[offset] {

			offset++
		}
		if offset < eol {
			if ch == '\'' || ch == '"' {
				quote := ch
				offset++
				if offset < eol {
					titleStart = offset
					for ch = line.runes[offset]; offset < eol && ch != quote; ch = line.runes[offset] {

						offset++
					}
				}
				if ch == quote {
					titleEnd = offset
				}
			}
		}
	}
	// XXX IF titleStart > 0 but titleEnd == 0, abort parse

	// XXX FOR STRICTNESS require offset = eol - 1
	if uriEnd > 0 {
		id := string(line.runes[idStart:idEnd])
		uri := line.runes[uriStart:uriEnd]
		var title []rune
		if titleEnd > 0 {
			title = line.runes[titleStart:titleEnd]
		}
		def, err = p.doc.addDefinition(id, uri, title)
	}
	return
}

// We are at the beginning of a line (possiblly with up to three leading
// spaces) and have seen an exclamation point followed by aleft square bracket.
// If we find the rest of
//   [id]:\s+(uri\s+("title"))
// where the optional title may be delimited with DQUOTE or SQUOTE, then we
// absorb all of these, adding id => DEF to the dictionary for the image.
// That is, a successful parse produces no output; it just affects the
// document dictionary.
//
// If there is any deviation from the spec, we leave the offset where it
// is and return a nil definition.  If the parse succeeds, we add the
// definition to the parser's dictionary, set the offset, and return a
// non-nil definition.
//
func (line *Line) parseImageDefinition(p *Parser) (def *Definition, err error) {

	var (
		ch                   rune
		idStart, idEnd       int
		offset               int
		uriStart, uriEnd     int
		titleStart, titleEnd int
	)
	// Enter having seen an exclamation point followed by a left square
	// bracket ('![') at the beginning of a line, possibly preceded by up
	// to three spaces.  The offset is on the exclamation point.

	eol := len(line.runes)
	offset = line.offset + 2 // just beyond the bracket

	// collect the id -----------------------------------------------
	for idStart = offset; offset < eol; offset++ {
		ch = line.runes[offset]
		if ch == ']' {
			idEnd = offset // exclusive end
			// DEBUG
			fmt.Printf("imageDef: idStart = %d, idEnd = %d\n",
				idStart, idEnd)
			// END
			offset++ // position beyond right bracket
			break
		}
	}
	// expect a colon and zero or more spaces -----------------------
	if idEnd > 0 && offset < eol-3 {
		if line.runes[offset] == ':' {
			offset++
			// skip any spaces
			for ch = line.runes[offset]; offset < eol && u.IsSpace(ch); ch = line.runes[offset] {

				offset++
			}
			if offset < eol-1 && ch == '(' {
				offset++
				uriStart = offset
			}
		}
	}
	// collect the uri ----------------------------------------------
	if uriStart > 0 {
		// assume that a uri contains no spaces
		for ; offset < eol; offset++ {
			ch = line.runes[offset]
			if u.IsSpace(ch) || ch == ')' {
				break
			}
		}
		if offset < eol {
			uriEnd = offset
			// DEBUG
			fmt.Printf("imageDef: uriStart = %d, uriEnd = %d\n",
				uriStart, uriEnd)
			// END
		}
	}
	// collect any title
	if uriEnd > 0 && offset < eol {
		// skip any spaces
		for ; offset < eol; offset++ {
			ch = line.runes[offset]
			if !u.IsSpace(ch) {
				break
			}
		}
		if offset < eol {
			if ch == '\'' || ch == '"' {
				quote := ch
				offset++
				if offset < eol {
					titleStart = offset
					for ch = line.runes[offset]; offset < eol && ch != quote; ch = line.runes[offset] {

						offset++
					}
				}
				if ch == quote {
					titleEnd = offset
					// DEBUG
					fmt.Printf("imageDef: uriStart = %d, uriEnd = %d\n",
						uriStart, uriEnd)
					// END
					offset++
				}
			}
		}
	}
	if uriEnd > 0 && offset < eol {
		if line.runes[offset] != ')' {
			// expect a closing RPAREN
			// DEBUG
			fmt.Printf("expected closing paren but found '%c'\n",
				line.runes[offset])
			// END
			uriEnd = 0
		} else if titleStart > 0 && titleEnd == 0 {
			// abort parse
			fmt.Printf("problem with title\n") // DEBUG
			uriEnd = 0
		} else if offset != eol-1 {
			// DEBUG
			fmt.Printf("offset %d but eol is %d\n",
				offset, eol)
			// END
			uriEnd = 0
		}
	}
	if uriEnd > 0 {
		id := string(line.runes[idStart:idEnd])
		uri := line.runes[uriStart:uriEnd]
		var title []rune
		if titleEnd > 0 {
			title = line.runes[titleStart:titleEnd]
		}
		def, err = p.doc.addDefinition(id, uri, title)
	}
	return
}
