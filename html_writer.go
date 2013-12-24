package md

// xgo/md/html_writer.go

import (
	"bytes"
	"fmt"
	"io"
	u "unicode"
)

var _ = fmt.Print

type HtmlWriter struct {
	wr io.Writer
}

// An HtmlWriter converts markdown objects ("downers") to serial form.
func NewHtmlWriter(w io.Writer) (hw *HtmlWriter, err error) {
	if w == nil {
		err = NilWriter
	} else {
		hw = &HtmlWriter{
			wr: w,
		}
	}
	return
}

// Output a slice of Markdown objects ("downers") as a sequence of bytes.
//
// XXX This implementation is simplistic.  The downers contain
// unfiltered runes, which may include characters such as ampersands ('&')
// which need to be converted to character entities.
func (hw *HtmlWriter) Write(downers []MarkdownI) (
	bytesOut int, err error) {

	var count2 int
	for i := 0; err == nil && i < len(downers); i++ {
		runes := downers[i].Get()

		var (
			b     bytes.Buffer
			count int
		)
		p := &b
		n := len(runes)
		for j := 0; j < n; j++ {
			r := runes[j]
			if r == BACKSLASH {
				if j < n-1 {
					nextChar := runes[j+1]
					if escaped(nextChar) {
						count, _ = p.WriteRune(nextChar)
						j++
					} else {
						// just output the backslash
						count, _ = p.WriteRune(r)
					}
				} else {
					// backslash is last character
					count, _ = p.WriteRune(r)
				}
			} else if r == '&' {
				if j == n-1 || u.IsSpace(runes[j+1]) {
					count, _ = p.WriteString("&amp;")
				} else {
					// handle entities
					semiOffset := -1
					for k := j + 1; k < n; k++ {
						if k > j+MAX_ENTITY_CHAR {
							break
						}
						ch := runes[k]
						if ch == ';' {
							semiOffset = k
							break
						}
					}
					// weird errors are possible here, should a semicolon
					// appear within MAX_ENTITY_CHAR of the ampersand
					if semiOffset > 0 {
						// found end of entity, so just output the ampersand
						count, _ = p.WriteRune(r)
					} else {
						// no semicolon, so just:
						count, _ = p.WriteString("&amp;")
					}
				}
			} else {
				count, _ = p.WriteRune(r)
			}
			bytesOut += count
		}
		count2, err = hw.wr.Write(p.Bytes())
		_ = count2
	}
	return
}
