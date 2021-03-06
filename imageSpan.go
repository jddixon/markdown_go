package md

// xgo/md/imageSpan.go

import (
	"fmt"
)

var _ = fmt.Print

type ImageSpan struct {
	linkText []rune
	path     []rune
	title    []rune
}

func NewImageSpan(linkText, path, title []rune) (t *ImageSpan) {

	link := make([]rune, len(linkText))
	copy(link, linkText)

	where := make([]rune, len(path))
	copy(where, path)

	t = &ImageSpan{
		linkText: link,
		path:     where,
	}
	if title != nil && len(title) > 0 {
		tt := make([]rune, len(title))
		copy(tt, title)
		t.title = tt
	}
	return
}

func (p *ImageSpan) String() string {
	if len(p.title) == 0 {
		return fmt.Sprintf("![%s] (%s) ", p.linkText, p.path)
	} else {
		return fmt.Sprintf("![%s] (%s \"%s\") ",
			p.linkText, p.path, p.title)
	}
}

func (p *ImageSpan) GetHtml() (out []rune) {

	out = append(out, []rune("<img src=\"")...)
	out = append(out, p.path...)
	out = append(out, []rune("\"")...)
	if p.linkText != nil {
		out = append(out, []rune(" alt=\"")...)
		out = append(out, p.linkText...)
		out = append(out, '"')
	}
	if p.title != nil {
		out = append(out, []rune(" title=\"")...)
		out = append(out, p.title...)
		out = append(out, '"')
	}
	out = append(out, []rune(" />")...)
	return
}

// PARSE ------------------------------------------------------------

// Attempt to parse out a ImageSpan, returning a SpanI reference
// to it on success and nil and possibly an error on failure.  If the parse
// fails but there is no input error, leave the line offset unchanged and
// return a nil SpanI.  If the parse succeeds, return a SpanI and advance
// the offset accordingly.
//
// In Markdown serialization, a ImageSpan looks like
//     [linkText](PATH "optional title")
// That is, it begins with linkText enclosed in square brackets and
// ends with a PATH or path (in the file system) enclosed by parentheses.
// We make no attempt to verify that the PATH is well-formed.
//
func (q *Line) parseImageSpan(opt *Options) (span SpanI, err error) {

	offset := q.offset + 2 // enter having seen ![
	var (
		linkTextStart         uint = offset
		linkTextEnd           uint
		pathStart, pathEnd    uint
		titleStart, titleEnd  uint
		end                   uint // offset of closing paren, if found
		linkText, path, title []rune
		verbose, testing      bool
	)
	if opt == nil {
		err = NilOptions
	} else {
		EOL := uint(len(q.runes))
		verbose = opt.Verbose
		testing = opt.Testing
		_ = verbose

		// look for the end of the linkText
		for ; offset < EOL; offset++ {
			ch := q.runes[offset]
			if ch == ']' {
				linkTextEnd = offset
				// DEBUG
				if testing {
					fmt.Printf("linkTextEnd = %d\n", offset)
				}
				// END
				offset++
				break
			}
		}
		if linkTextEnd > 0 {
			// optional space
			if offset < EOL-1 && q.runes[offset+1] == ' ' {
				offset++
			}
			if (offset < EOL) && (q.runes[offset] == '(') {
				offset++
				pathStart = offset
				// DEBUG
				if testing {
					fmt.Printf("pathStart = %d\n", offset)
				}
				// END
			}
		}
		if pathStart > 0 {
			for offset = pathStart; offset < EOL; offset++ {
				ch := q.runes[offset]
				if ch == ')' {
					end = offset
					if pathEnd == 0 {
						pathEnd = end
					}
					break
				}
				if ch == '"' {
					if titleStart == 0 {
						pathEnd = offset
						if q.runes[pathEnd-1] == ' ' {
							pathEnd--
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
				// just give up
				end = 0
			}
		}
		if end > 0 {
			linkText = q.runes[linkTextStart:linkTextEnd]
			path = q.runes[pathStart:pathEnd]
			if titleStart > 0 {
				title = q.runes[titleStart:titleEnd]
			}

			span = NewImageSpan(linkText, path, title)
			q.offset = offset + 1
		} // FOO
	}
	return
}
