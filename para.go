package md

// xgo/md/para.go

import (
// "fmt"
)

var (
	PARA_OPEN  = []rune("<p>")
	PARA_CLOSE = []rune("</p>")
)

type Para struct {
	spans []SpanI
}

// A Para is a Block consisting of a number of SpanI, things that
// implement the SpanI interface.
func NewPara(span SpanI) (pa *Para) {
	pa = &Para{}
	if span != nil {
		pa.spans = append(pa.spans, span)
	}
	return
}

func (p *Para) Add(span SpanI) (err error) {
	p.spans = append(p.spans, span)
	return
}

func (p *Para) Get() (runes []rune) {
	runes = append(runes, PARA_OPEN...)
	for i := 0; i < len(p.spans); i++ {
		runes = append(runes, p.spans[i].Get()...)
	}
	runes = append(runes, PARA_CLOSE...)
	return
}

// Parse a line known to begin a new Para.
func (q *Line) parsePara() (b BlockI, err error) {
	var pa *Para
	spans, err := q.parseToSpans()
	if err == nil {
		pa = new(Para)
		for i := 0; i < len(spans); i++ {
			err = pa.Add(spans[i])
			if err != nil {
				break
			}
		}
		if err == nil {
			b = pa
		}
	}
	return
}
