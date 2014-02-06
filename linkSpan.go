package md

// xgo/md/linkSpan.go

import (
	"fmt"
	"strings"
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

func (p *LinkSpan) String() string {
	if len(p.title) == 0 {
		return fmt.Sprintf("[%s] (%s) ", p.linkText, p.uri)
	} else {
		return fmt.Sprintf("[%s] (%s \"%s\") ",
			p.linkText, p.uri, p.title)
	}
}

func (p *LinkSpan) GetHtml() (out []rune) {

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
func (q *Line) parseLinkSpan(opt *Options) (span SpanI, err error) {

	offset := q.offset + 1
	var (
		linkTextStart        uint = offset
		linkTextEnd          uint
		uriStart, uriEnd     uint
		titleStart, titleEnd uint
		end                  uint // offset of closing paren, if found
		linkText, uri, title []rune
		verbose, testing     bool
	)
	if opt == nil {
		err = NilOptions
	} else {
		verbose = opt.Verbose
		testing = opt.Testing
		_ = verbose

		// DEBUG
		if testing {
			fmt.Printf("parseLinkSpan: offset %d, text %s\n",
				offset, string(q.runes[offset:]))
		}
		// END

		// look for the end of the linkText
		for ; offset < uint(len(q.runes)); offset++ {
			ch := q.runes[offset]
			if ch == ']' {
				linkTextEnd = offset
				// DEBUG
				if testing {
					fmt.Printf("linkTextEnd = %d; end is %d\n",
						offset, uint(len(q.runes))) // DEBUG
				}
				// END
				offset++
				break
			}
		}
		if (offset < uint(len(q.runes))-1) && linkTextEnd > 0 {
			// optional space
			if q.runes[offset] == ' ' {
				offset++
			}
			if q.runes[offset] == '(' {
				offset++
				uriStart = offset
			}
		}
		if uriStart > 0 {
			for offset = uriStart; offset < uint(len(q.runes)); offset++ {
				ch := q.runes[offset]
				if ch == ')' {
					end = offset
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
				// DEBUG
				if testing {
					fmt.Printf("found start of title but not end\n")
				}
				// END
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
		} // FOO
	}
	return
}

// An automatic link looks like <http://www.w3c.org>
// XXX This simplistic approach will treat an email address as a
// URL.  Need to look tat the content to decide how to handle.
func (q *Line) parseAutomaticLink() (span SpanI, err error) {

	offset := q.offset
	eol := uint(len(q.runes))
	found := false

	// we require that the cursor is on the <
	offset++

	// look for the end of the span
	for offset < eol {
		ch := q.runes[offset]
		if ch == '>' {
			found = true
			break
		}
		offset++
	}
	if found {
		var start, end uint
		start = q.offset + 1
		end = offset
		body := q.runes[start:end]
		// XXX Rough decision as to whether it's a link or an email address.
		strBody := string(body)
		if strings.HasPrefix(strBody, "http") {
			q.offset = offset + 1
			span = NewLinkSpan(body, body, body[:0])
		} else if strings.Contains(strBody, "@") {
			body2 := MAIL_TO
			for i := 0; i < len(body); i++ {
				var repr string
				if i%2 == 0 {
					repr = fmt.Sprintf("&#x%x", body[i]) // so hex
				} else {
					repr = fmt.Sprintf("&#%x", body[i]) // decimal
				}
				addMe := []rune(repr)
				body2 = append(body2, addMe...)
			}
			q.offset = offset + 1
			span = NewLinkSpan(body2, body2, body2[:0])
		}
		// otherwise we just ignore the attempt
	}
	return
}
