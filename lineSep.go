package md

// xgo/md/lineSep.go

import (
	"fmt"
)

var _ = fmt.Print

type LineSep struct {
	runes []rune
}

// Whether this character is in our list of line separator characters.
func IsSepChar(ch rune) bool {
	for i := 0; i < len(SEP_RUNES); i++ {
		if ch == SEP_RUNES[i] {
			return true
		}
	}
	return false
}

// Make a new line separator container.
func NewLineSep(chars []rune) (ls *LineSep, err error) {
	for i := 0; i < len(chars); i++ {
		ch := chars[i]
		if !IsSepChar(ch) {
			err = NotALineSeparator
			break
		}
	}
	if err == nil {
		runes := make([]rune, len(chars))
		copy(runes, chars)
		ls = &LineSep{runes: runes}
	}
	return
}

func (ls *LineSep) String() string {
	return string(ls.runes)
}

// Add a character if you know it's a line separator.
func (ls *LineSep) add(ch rune) {
	ls.runes = append(ls.runes, ch)
}

// Add a character if it should be a line separator.
func (ls *LineSep) Add(ch rune) (err error) {
	if !IsSepChar(ch) {
		err = NotALineSeparator
	}
	if err == nil {
		ls.add(ch)
	}
	return
}

// Get the line separators collected so far.
// XXX This may not accord with Markdown practice.
func (ls *LineSep) GetHtml() []rune {
	// DEBUG
	//fmt.Printf("LineSep.GetHtml() returning %d runes; first is 0x%02x\n",
	//	len(ls.runes), ls.runes[0])
	// END
	return ls.runes
}
