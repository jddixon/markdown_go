package md

// xgo/md/para.go

import (
	"fmt"
	"strings"
)

var _ = fmt.Print

// A Para is a Block consisting of a number of SpanSeq, each of which is
// a sequence of SpanI, things that implement the SpanI interface.  In
// addition, each SpanSeq has an associated lineSeq, a slice of line
// terminators (null, CR, or LF)
type Para struct {
	seqs []SpanSeq
}

func (p *Para) String() string {
	var ss []string
	for i := 0; i < len(p.seqs); i++ {
		ss = append(ss, p.seqs[i].String())
	}
	return strings.Join(ss, "")
}

func (p *Para) Add(seq SpanSeq) (err error) {
	p.seqs = append(p.seqs, seq)
	return
}

func (p *Para) GetHtml() (runes []rune) {
	runes = append(runes, PARA_OPEN...)
	for i := 0; i < len(p.seqs); i++ {
		spans := p.seqs[i].spans
		for j := 0; j < len(spans); j++ {
			runes = append(runes, spans[j].GetHtml()...)
		}
		if i < len(p.seqs)-1 {
			runes = append(runes, p.seqs[i].lineSep...)
		}
	}
	ndxLast := len(runes) - 1
	for ndxLast >= 0 && (runes[ndxLast] == CR || runes[ndxLast] == LF) {
		runes = runes[:ndxLast]
		ndxLast = len(runes) - 1
	}

	runes = append(runes, PARA_CLOSE...)
	lastLineSep := p.seqs[len(p.seqs)-1].lineSep
	if len(lastLineSep) > 0 {
		runes = append(runes, lastLineSep...)
	} else {
		runes = append(runes, LF)
	}
	// DEBUG
	// fmt.Printf("SPAN: '%s'\n", string(runes))
	// END
	return
}
