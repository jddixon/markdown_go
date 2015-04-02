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

func (h *HRule) String() string {
	return "---\n"
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
// 2015-04-01: Commonmark allows zero to three spaces at the beginning
// of the line and any number of spaces at the end of the line.
//
// If the parse succeeds we return a pointer to the HRule object.
// Otherwise the offset is unchanged and b's value is nil.
func (q *Line) parseHRule(from uint) (b BlockI, err error) {

	var (
		badCharSeen bool
		eol         uint = uint(len(q.runes))
		offset      uint = from
		char        rune
		matchCount  int
		spaceStart  uint = offset
	)
	if offset < eol {
		// ignore up to three leading spaces
		for char = q.runes[offset]; char == ' '; char = q.runes[offset] {
			if offset + 1 < eol {
				offset++
			} else {
				break
			}
		}
		if (offset < eol) && (offset < spaceStart + 4) && 
			(char == '-' || char == '*' || char == '_' ){
	
			matching := char
			matchCount++
			for offset++; offset < eol; offset++ {
				char = q.runes[offset]
				if char == matching {
					matchCount++
				} else {
					// we allow spaces mixed in with matching characters
					if !u.IsSpace(char) {
						badCharSeen = true
						break
					}
				}
			}
		}
		if matchCount >= 3 && !badCharSeen {
			b = &HRule{}
		}
	}
	return
}
