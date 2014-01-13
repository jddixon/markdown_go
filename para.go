package md

// xgo/md/para.go

import (
	"fmt"
)

var _ = fmt.Print

// A Para is a Block consisting of a number of SpanSeq, each of which is
// a sequence of SpanI, things that implement the SpanI interface.  In
// addition, each SpanSeq has an associated lineSeq, a slice of line
// terminators (null, CR, or LF)
type Para struct {
	seqs []SpanSeq
}

func (p *Para) Add(seq SpanSeq) (err error) {
	p.seqs = append(p.seqs, seq)
	return
}

func (p *Para) Get() (runes []rune) {
	runes = append(runes, PARA_OPEN...)
	for i := 0; i < len(p.seqs); i++ {
		spans := p.seqs[i].spans
		for j := 0; j < len(spans); j++ {
			runes = append(runes, spans[j].Get()...)
		}
		if i < len(p.seqs)-1 {
			runes = append(runes, p.seqs[i].lineSep...)
		}
	}
	runes = append(runes, PARA_CLOSE...)
	runes = append(runes, p.seqs[len(p.seqs)-1].lineSep...)
	// DEBUG
	// fmt.Printf("SPAN: '%s'\n", string(runes))
	// END
	return
}
