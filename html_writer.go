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

// Output a slice of Markdown runes as a sequence of bytes.
//
// XXX This implementation is simplistic.  The slice contains
// unfiltered runes, which may include characters such as ampersands ('&')
// which need to be converted to character entities.
func (hw *HtmlWriter) Write(runes []rune) (
	bytesOut int, err error) {

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
					if k > j+MAX_ENTITY_CHARS {
						break
					}
					ch := runes[k]
					if ch == ';' {
						semiOffset = k
						break
					}
				}
				// weird errors are possible here, should a semicolon
				// appear within MAX_ENTITY_CHARS of the ampersand
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
	_, err = hw.wr.Write(p.Bytes()) // XXX count ignored
	return
}
