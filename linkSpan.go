package md

// xgo/md/linkSpan.go

import (
	"fmt"
)

var _ = fmt.Print

type LinkSpan struct {
	linkText []rune
	uri      []rune
	title    []rune
}

func NewLinkSpan(linkText, uri, title []rune) (t *LinkSpan) {

	link := make([]rune, len(linkText))
	copy(link, linkText)

	where := make([]rune, len(uri))
	copy(where, uri)

	t = &LinkSpan{
		linkText: link,
		uri:      where,
	}
	if title != nil && len(title) > 0 {
		tt := make([]rune, len(title))
		copy(tt, title)
		t.title = tt
	}
	return
}

func (p *LinkSpan) Get() (out []rune) {

	out = append(out, []rune("<a href=\"")...)
	out = append(out, p.uri...)
	if p.title != nil {
		out = append(out, []rune("\" title=\"")...)
		out = append(out, p.title...)
	}
	out = append(out, []rune("\">")...)
	if p.linkText != nil {
		out = append(out, p.linkText...)
	}
	out = append(out, []rune("</a>")...)
	return
}

// PARSE ------------------------------------------------------------

// Attempt to parse out a LinkSpan, returning a SpanI reference
// to it on success and nil and possibly an error on failure.  If the parse
// fails but there is no input error, leave the line offset unchanged and
// return a nil SpanI.  If the parse succeeds, return a SpanI and advance
// the offset accordingly.
//
// In Markdown serialization, a LinkSpan looks like
//     [linkText](URL "optional title")
// That is, it begins with linkText enclosed in square brackets and
// ends with a URL or path (in the file system) enclosed by parentheses.
// We make no attempt to verify that the URI is well-formed.
//
func (q *Line) parseLinkSpan() (span SpanI, err error) {

	offset := q.offset + 1
	var (
		linkTextStart        int = offset
		linkTextEnd          int
		uriStart, uriEnd     int
		titleStart, titleEnd int
		end                  int // offset of closing paren, if found
		linkText, uri, title []rune
	)

	// DEBUG
	fmt.Printf("parseLinkSpan: offset %d, text %s\n",
		offset, string(q.runes[offset:]))
	// END

	// look for the end of the linkText
	for ; offset < len(q.runes); offset++ {
		ch := q.runes[offset]
		if ch == ']' {
			linkTextEnd = offset
			// DEBUG
			fmt.Printf("linkTextEnd = %d; end is %d\n",
				offset, len(q.runes)) // DEBUG
			// END
			offset++
			break
		}
	}
	if (offset < len(q.runes)-1) && linkTextEnd > 0 {
		// optional space
		if q.runes[offset] == ' ' {
			fmt.Printf("skipping space at %d\n", offset)
			offset++
		}
		if q.runes[offset] == '(' {
			offset++
			uriStart = offset
			// fmt.Printf("uriStart = %d\n", offset) // DEBUG
		}
	}
	if uriStart > 0 {
		for offset = uriStart; offset < len(q.runes); offset++ {
			ch := q.runes[offset]
			if ch == ')' {
				end = offset
				fmt.Printf("FOUND RPAREN LinkSpan END at %d\n", end)
				if uriEnd == 0 {
					uriEnd = end
				}
				break
			}
			if ch == '"' {
				if titleStart == 0 {
					uriEnd = offset
					if q.runes[uriEnd-1] == ' ' {
						uriEnd--
					}
					titleStart = offset + 1 // inclusive
				} else {
					titleEnd = offset // exclusive
				}
			}
		}
	}
	if end > 0 {
		if titleStart > 0 && titleEnd == 0 {
			fmt.Printf("found start of title but not end\n") // DEBUG
			// just give up
			end = 0
		}
	}
	if end > 0 {
		linkText = q.runes[linkTextStart:linkTextEnd]
		uri = q.runes[uriStart:uriEnd]
		if titleStart > 0 {
			title = q.runes[titleStart:titleEnd]
		}

		span = NewLinkSpan(linkText, uri, title)
		q.offset = offset + 1
	}
	return
}
