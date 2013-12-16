package md

// xgo/md/lineSep.go

import (
	"fmt"
)

type LineSep struct {
	runes []rune
}

// Whether this character is in our list of line separator characters.
func IsSepChar(ch rune) bool {
	for i := 0; i < len(SEP_CHAR); i++ {
		if ch == SEP_CHAR[i] {
			return true
		}
	}
	return false
}

// Make a new line separator container.
func NewLineSep(ch rune) (ls *LineSep, err error) {
	if !IsSepChar(ch) {
		err = NotALineSeparator
	}
	if err == nil {
		runes := []rune{ch}
		ls = &LineSep{runes: runes}
	}
	return
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
func (ls *LineSep) Get() []rune {
	// DEBUG
	fmt.Printf("LineSep.Get: %d separators: ", len(ls.runes))
	for i := 0; i < len(ls.runes); i++ {
		fmt.Printf("%d ", int(ls.runes[i]))
	}
	fmt.Println()
	// END
	return ls.runes
}
