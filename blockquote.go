package md

// xgo/md/blockquote.go

import (
	"fmt"
)

var _ = fmt.Print

// A Blockquote is a Block consisting of a number of BlockIs.
type Blockquote struct {
	Holder
}

func (q *Line) parseBlockquote(doc *Document, from int) (
	h HolderI, err error) {

	// Enter having seen "> " and with offset set to the next character
	var seq *SpanSeq
	seq, err = q.parseSpanSeq(doc, from, true) // true means trim left
	curPara := new(Para)
	curPara.seqs = append(curPara.seqs, *seq)
	if err == nil {
		bq := new(Blockquote)
		err = bq.AddChild(curPara)
		if err == nil {
			h = bq
		}
	}
	return
}

func (bq *Blockquote) Get() (runes []rune) {
	runes = append(runes, BLOCKQUOTE_OPEN...)
	for i := 0; i < bq.Size(); i++ {
		child, _ := bq.GetChild(i)
		runes = append(runes, child.Get()...)
	}
	// XXX A HACK
	runes = append(runes, '\n')
	// END HACK
	runes = append(runes, BLOCKQUOTE_CLOSE...)
	return
}
