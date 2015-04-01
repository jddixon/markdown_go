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
	isImg bool
}

func (def *Definition) GetURI() string {
	return string(def.uri)
}

func (def *Definition) GetTitle() string {
	return string(def.title)
}

// We are at the beginning of a line (possiblly with up to three leading
// spaces) and have seen a left square bracket.  If we find the rest of
//   [id]:\s+uri\s?("title")?
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
func (line *Line) parseLinkDefinition(opt *Options, doc *Document) (
	def *Definition, err error) {

	var (
		ch                   rune
		EOL                  uint
		idStart, idEnd       uint
		offset               uint
		uriStart, uriEnd     uint
		titleStart, titleEnd uint
		verbose, testing     bool
	)
	// Enter having seen a left square bracket ('[') at the beginning
	// of a line, possibly preceded by up to three spaces.  The offset
	// is on the bracket.
	if opt == nil {
		err = NilOptions
	} else {
		verbose = opt.Verbose
		testing = opt.Testing
		_ = verbose
		EOL = uint(len(line.runes))
		offset = line.offset + 1 // just beyond the bracket

		// collect the id -----------------------------------------------
		for idStart = offset; offset < EOL; offset++ {
			ch = line.runes[offset]
			if ch == ']' {
				idEnd = offset // exclusive end
				offset++       // position beyond right bracket
				break
			}
		}
		// expect a colon and one or more spaces ------------------------
		if (idEnd > 0) && (offset+3 < EOL) {
			if line.runes[offset] == ':' {		// XXX
				offset++
				// skip any spaces
				for ch = line.runes[offset]; offset < EOL && u.IsSpace(ch); ch = line.runes[offset] {

					offset++
				}
				if offset < EOL {
					uriStart = offset
				}
			}
		}
		// collect the uri ----------------------------------------------
		if uriStart > 0 {
			// assume that a uri contains no spaces
			for offset < EOL && !u.IsSpace(line.runes[offset]) {
				offset++
			}
			uriEnd = offset
			if line.runes[uriStart] == '<' {
				uriStart++
				if line.runes[uriEnd-1] == '>' { // end is exclusive
					uriEnd--
				} else {
					// an error, so force the parse to fail
					uriEnd = 0
				}
			}
		}
		// collect any title
		if uriEnd > 0 && offset < EOL {
			// skip any spaces
			for ch = line.runes[offset]; offset < EOL && u.IsSpace(ch); ch = line.runes[offset] {

				offset++
			}
			if offset < EOL {
				if ch == '\'' || ch == '"' || ch == '(' {
					openQuote := ch
					var closeQuote rune
					if openQuote == '(' {
						closeQuote = ')'
					} else {
						closeQuote = openQuote
					}
					offset++
					if offset < EOL {
						titleStart = offset
						for ch = line.runes[offset]; 
							(offset < EOL - 1) && (ch != closeQuote); 
							ch = line.runes[offset] {

							offset++
						}
					}
					if ch == closeQuote { // GEEP
						titleEnd = offset
					}
				}
			}
		}
		// XXX IF titleStart > 0 but titleEnd == 0, abort parse

		// XXX FOR STRICTNESS require offset = EOL - 1
		if uriEnd > 0 {
			id := string(line.runes[idStart:idEnd])
			uri := line.runes[uriStart:uriEnd]
			var title []rune
			if titleEnd > 0 {
				title = line.runes[titleStart:titleEnd]
			}
			def, err = doc.AddDefinition(id, uri, title, false) // isImg=false
			// DEBUG
			if testing {
				if def == nil {
					fmt.Printf("parseLinkDefn returning NIL, error %s\n",
						err.Error())
				} else {
					fmt.Println("parseLinkDefn returning definition")
				}
			}
			// END
		} // FOO
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
func (line *Line) parseImageDefinition(opt *Options, doc *Document) (
	def *Definition, err error) {

	var (
		ch                   rune
		idStart, idEnd       uint
		offset               uint
		uriStart, uriEnd     uint
		titleStart, titleEnd uint
		verbose, testing     bool
	)
	// Enter having seen an exclamation point followed by a left square
	// bracket ('![') at the beginning of a line, possibly preceded by up
	// to three spaces.  The offset is on the exclamation point.

	if opt == nil {
		err = NilOptions
	} else {
		verbose = opt.Verbose
		testing = opt.Testing
		_ = verbose
		EOL := uint(len(line.runes))
		offset = line.offset + 2 // just beyond the bracket

		// collect the id -----------------------------------------------
		for idStart = offset; offset < EOL; offset++ {
			ch = line.runes[offset]
			if ch == ']' {
				idEnd = offset // exclusive end
				offset++       // position beyond right bracket
				break
			}
		}
		// expect a colon and zero or more spaces -----------------------
		if idEnd > 0 && offset < EOL-3 {
			if line.runes[offset] == ':' {
				offset++
				// skip any spaces
				for ch = line.runes[offset]; offset < EOL && u.IsSpace(ch); ch = line.runes[offset] {

					offset++
				}
				if offset < EOL-1 && ch == '(' {
					offset++
					uriStart = offset
				}
			}
		}
		// collect the uri ----------------------------------------------
		if uriStart > 0 {
			// assume that a uri contains no spaces
			for ; offset < EOL; offset++ {
				ch = line.runes[offset]
				if u.IsSpace(ch) || ch == ')' {
					break
				}
			}
			if offset < EOL {
				uriEnd = offset
			}
		}
		// collect any title
		if uriEnd > 0 && offset < EOL {
			// skip any spaces
			for ; offset < EOL; offset++ {
				ch = line.runes[offset]
				if !u.IsSpace(ch) {
					break
				}
			}
			if offset < EOL {
				if ch == '\'' || ch == '"' {
					quote := ch
					offset++
					if offset < EOL {
						titleStart = offset
						for ch = line.runes[offset]; offset < EOL && ch != quote; ch = line.runes[offset] {

							offset++
						}
					}
					if ch == quote {
						titleEnd = offset
						offset++
					}
				}
			}
		}
		if uriEnd > 0 && offset < EOL {
			if line.runes[offset] != ')' {
				// expect a closing RPAREN
				// DEBUG
				if testing {
					fmt.Printf("expected closing paren but found '%c'\n",
						line.runes[offset])
				}
				// END
				uriEnd = 0
			} else if titleStart > 0 && titleEnd == 0 {
				// abort parse
				// DEBUG
				if testing {
					fmt.Printf("problem with title\n")
				}
				// END
				uriEnd = 0
			} else if offset != EOL-1 {
				// DEBUG
				if testing {
					fmt.Printf("offset %d but EOL is %d\n",
						offset, EOL)
				}
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
			def, err = doc.AddDefinition(id, uri, title, true) //isImg = true
			// DEBUG
			if testing {
				if def == nil {
					fmt.Printf("parseImageDefn returning NIL, error %s\n",
						err.Error())
				} else {
					fmt.Println("parseImageDefn returning definition")
				}
			}
			// END
			// DEBUG
		} else {
			if testing {
				fmt.Println("ImageDef parse failed, uriEnd is zero")
			}
			// END
		} // FOO/
	}
	return
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
} // GEEP
