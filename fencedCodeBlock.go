package md

// xgo/md/fencedFencedCodeBlock.go

import (
	"fmt"
	"strings"
)

var _ = fmt.Print

// A FencedCodeBlock is a Block consisting of a number of lines, collected
// here as CodeLines
type FencedCodeBlock struct {
	language string     // may not be present
	lines    []CodeLine // to force conversion to entities
}

func (p *FencedCodeBlock) String() string {
	var ss []string

	ss = append(ss, string(FENCE))
	for i := 0; i < len(p.lines); i++ {
		s := p.lines[i].String() // NL-terminated
		ss = append(ss, s)
	}
	ss = append(ss, string(FENCE))
	return strings.Join(ss, "")
}

// XXX Might make more sense to copy.
func (p *FencedCodeBlock) Add(line *CodeLine) {
	p.lines = append(p.lines, *line)
	return
}

func (p *FencedCodeBlock) Clear() {
	p.lines = p.lines[:0]
}

func (p *FencedCodeBlock) GetHtml() (runes []rune) {
	runes = append(runes, CODE_OPEN...)
	for i := 0; i < len(p.lines); i++ {
		text := p.lines[i].GetHtml()
		runes = append(runes, text...)
		runes = append(runes, LF)
	}
	runes = append(runes, CODE_CLOSE...)
	return
}

func (p *FencedCodeBlock) Language() string {
	return p.language
}
func (p *FencedCodeBlock) Size() uint {
	return uint(len(p.lines))
}
