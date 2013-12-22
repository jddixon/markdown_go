package md

// xgo/md/imageRefSpan.go

import (
	"fmt"
)

// XXX --------------------------------------------------------------
// THIS IS NOT QUITE RIGHT -- JUST HACKED FROM linkRefSpan.go
// XXX --------------------------------------------------------------

var _ = fmt.Print

// In Markdown serialization, a ImageRef looks like
//     [altText][id]
type ImageRefSpan struct {
	altText []rune
	id       string
	p        *Parser
}

func NewImageRefSpan(p *Parser, altText []rune, id string) (
	t *ImageRefSpan, err error) {

	if p == nil {
		err = NilParser
	} else {
		image := make([]rune, len(altText))
		copy(image, altText)

		t = &ImageRefSpan{
			altText: image,
			id:       id,
			p:        p,
		}
	}
	return
}

// XXX WE NEED A DICTIONARY TO MAKE THIS WORK

func (ls *ImageRefSpan) Get() (out []rune) {

	def := ls.p.dict[ls.id]
	uri := def.uri
	title := def.title

	out = append(out, []rune("<a href=\"")...)
	out = append(out, uri...)
	if title != nil {
		out = append(out, []rune("\" title=\"")...)
		out = append(out, title...)
	}
	out = append(out, []rune("\">")...)
	if ls.altText != nil {
		out = append(out, ls.altText...)
	}
	out = append(out, []rune("</a>")...)
	return
}

// PARSE ------------------------------------------------------------

// Attempt to parse out a ImageRefSpan, returning a SpanI reference
// to it on success and nil and possibly an error on failure.  If the parse
// fails but there is no input error, leave the line offset unchanged and
// return a nil SpanI.  If the parse succeeds, return a SpanI and advance
// the offset accordingly.
//
// In Markdown serialization, a ImageRefSpan looks like
//     [altText][id]
// That is, it begins with altText enclosed in square brackets.  This
// is optionally followed by a space.  An id in square brackets follows.
// We make no attempt to verify that the id is well-formed.
//
func (q *Line) parseImageRefSpan(p *Parser) (span SpanI, err error) {

	if p == nil {
		err = NilParser
	} else {
		offset := q.offset + 2	// Enter having seen ![
		var (
			altTextStart  int = offset
			altTextEnd    int
			idStart, idEnd int
			end            int // offset of closing paren, if found
			altText, id   []rune
		)

		// look for the end of the altText -------------------------
		for ; offset < len(q.runes); offset++ {
			ch := q.runes[offset]
			if ch == ']' {
				altTextEnd = offset
				fmt.Printf("altTextEnd = %d\n", offset) // DEBUG
				offset++
				break
			}
		}
		if altTextEnd > 0 {
			// optional space
			if offset < len(q.runes)-1 && q.runes[offset+1] == ' ' {
				offset++
			}
			if q.runes[offset] == '[' {
				offset++
				idStart = offset
				fmt.Printf("idStart = %d\n", offset) // DEBUG
			}
		}
		// find the end of the ID -----------------------------------
		if idStart > 0 {
			for offset = idStart; offset < len(q.runes); offset++ {
				ch := q.runes[offset]
				if ch == ']' {
					end = offset
					if idEnd == 0 {
						idEnd = end
					}
					break
				}
			}
		}
		if end > 0 {
			altText = q.runes[altTextStart:altTextEnd]
			id = q.runes[idStart:idEnd]
			span, err = NewImageRefSpan(p, altText, string(id))
			q.offset = offset + 1
		}
	}
	return
}
