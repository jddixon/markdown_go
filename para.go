package md

// xgo/md/para.go

import (
	"fmt"
)

var (
	PARA_START = []rune("<p>")
	PARA_END   = []rune("</p>")
)

type Para struct {
	runes []rune
}

// This definition allows any and all characters into a p element.
func NewPara(chars []rune) (ls *Para) {
	p := []rune("<p>")
	p = append(p, chars...)
	p = append(p, PARA_END...)
	ls = &Para{runes: p}
	// DEBUG
	fmt.Printf("NewPara: '%s' => '%s'\n",
		string(chars), string(ls.runes))
	// END
	return
}

// This definition allows any and all characters into a p element.
func (ls *Para) Add(ch rune) (err error) {
	if err == nil {
		ls.runes = append(ls.runes, ch)
	}
	return
}

func (ls *Para) Get() []rune {
	return ls.runes
}
