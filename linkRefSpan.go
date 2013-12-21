package md

// xgo/md/linkRefSpan.go

import (
	"fmt"
)

var _ = fmt.Print

// In Markdown serialization, a LinkRef looks like
//     [linkText][id]
type LinkRefSpan struct {
	linkText []rune
	id       string
	p        *Parser
}

func NewLinkRefSpan(p *Parser, linkText []rune, id string) (
	t *LinkRefSpan, err error) {

	if p == nil {
		err = NilParser
	} else {
		link := make([]rune, len(linkText))
		copy(link, linkText)

		t = &LinkRefSpan{
			linkText: link,
			id:       id,
			p:        p,
		}
	}
	return
}

// XXX WE NEED A DICTIONARY TO MAKE THIS WORK

func (ls *LinkRefSpan) Get() (out []rune) {

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
	if ls.linkText != nil {
		out = append(out, ls.linkText...)
	}
	out = append(out, []rune("</a>")...)
	return
}

// PARSE ------------------------------------------------------------

// Attempt to parse out a LinkRefSpan, returning a SpanI reference
// to it on success and nil and possibly an error on failure.  If the parse
// fails but there is no input error, leave the line offset unchanged and
// return a nil SpanI.  If the parse succeeds, return a SpanI and advance
// the offset accordingly.
//
// In Markdown serialization, a LinkRefSpan looks like
//     [linkText][id]
// That is, it begins with linkText enclosed in square brackets.  This
// is optionally followed by a space.  An id in square brackets follows.
// We make no attempt to verify that the id is well-formed.
//
func (q *Line) parseLinkRefSpan(p *Parser) (span SpanI, err error) {

	if p == nil {
		err = NilParser
	} else {
		offset := q.offset + 1
		var (
			linkTextStart  int = offset
			linkTextEnd    int
			idStart, idEnd int
			end            int // offset of closing paren, if found
			linkText, id   []rune
		)

		// look for the end of the linkText -------------------------
		for ; offset < len(q.runes); offset++ {
			ch := q.runes[offset]
			if ch == ']' {
				linkTextEnd = offset
				fmt.Printf("linkTextEnd = %d\n", offset) // DEBUG
				offset++
				break
			}
		}
		if linkTextEnd > 0 {
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
			linkText = q.runes[linkTextStart:linkTextEnd]
			id = q.runes[idStart:idEnd]
			span, err = NewLinkRefSpan(p, linkText, string(id))
			q.offset = offset + 1
		}
	}
	return
}
