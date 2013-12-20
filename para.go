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
