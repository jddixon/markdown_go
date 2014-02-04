package md

// xgo/md/hrule.go

import (
	"fmt"
	u "unicode"
)

var _ = fmt.Print

// This must implement BlockI
type HRule struct {
	runes []rune // contains nothing
}

func NewHRule() (h BlockI, err error) {
	h = &HRule{}
	return
}

func (h *HRule) GetHtml() []rune {
	return H_RULE
}

// In this implementation, a Markdown horizontal rule is denoted by
// a single line beginning with one of hyphen, asterisk, or underscore
// and containing at least three of that character, possibly separated
// by an arbitrary number of spaces or hyphens.  We enter with the
// line offset pointing to the special character ('-' or '*' or '_').
//
// If the parse succeeds we return a pointer to the HRule object.
// Otherwise the offset is unchanged and b's value is nil.
func (q *Line) parseHRule(from uint) (b BlockI, err error) {

	var (
		badCharSeen bool
		eol         uint = uint(len(q.runes))
		offset      uint = from
		char        rune = q.runes[offset]
		charCount   int  = 1
	)
	if char == '-' || char == '*' || char == '_' {
		for offset++; offset < eol; offset++ {
			ch := q.runes[offset]
			if ch == char {
				charCount++
			} else if !u.IsSpace(ch) {
				badCharSeen = true
			}
		}
	}
	if charCount >= 3 && !badCharSeen {
		b = &HRule{}
	}
	return
}
