package md

// xgo/md/imageRefSpan.go

import (
	"fmt"
)

var _ = fmt.Print

// In Markdown serialization, a ImageRef looks like
//     [altText][id]
type ImageRefSpan struct {
	altText []rune
	id      string
	doc     *Document
}

func NewImageRefSpan(doc *Document, altText []rune, id string) (
	t *ImageRefSpan, err error) {

	if doc == nil {
		err = NilDocument
	} else {
		image := make([]rune, len(altText))
		copy(image, altText)

		t = &ImageRefSpan{
			altText: image,
			id:      id,
			doc:     doc,
		}
	}
	return
}

func (ls *ImageRefSpan) GetHtml() (out []rune) {

	def := ls.doc.refDict[ls.id]
	uri := def.uri
	title := def.title

	out = append(out, []rune("<img src=\"")...)
	out = append(out, uri...)
	out = append(out, '"')
	if ls.altText != nil {
		out = append(out, []rune(" alt=\"")...)
		out = append(out, ls.altText...)
		out = append(out, '"')
	}
	if title != nil {
		out = append(out, []rune(" title=\"")...)
		out = append(out, title...)
		out = append(out, '"')
	}
	out = append(out, []rune(" />")...)
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
func (q *Line) parseImageRefSpan(doc *Document) (span SpanI, err error) {

	if doc == nil {
		err = NilParser
	} else {
		offset := q.offset + 2 // Enter having seen ![
		var (
			altTextStart   uint = offset
			altTextEnd     uint
			idStart, idEnd uint
			end            uint // offset of closing paren, if found
			altText, id    []rune
		)

		// look for the end of the altText -------------------------
		for ; offset < uint(len(q.runes)); offset++ {
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
			if offset < uint(len(q.runes))-1 && q.runes[offset+1] == ' ' {
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
			for offset = idStart; offset < uint(len(q.runes)); offset++ {
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
			span, err = NewImageRefSpan(doc, altText, string(id))
			q.offset = offset + 1
		}
	}
	return
}
