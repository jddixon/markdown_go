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

func NewBlockquote(opt *Options, depth uint) (bq *Blockquote, err error) {

	h, err := NewHolder(opt, true, depth)
	if err == nil {
		bq = &Blockquote{Holder: *h}
	}
	return
}

func (bq *Blockquote) GetHtml() (runes []rune) {
	// runes = append(runes, LF)
	runes = append(runes, BLOCKQUOTE_OPEN...)
	for i := 0; i < bq.Size(); i++ {
		block, _ := bq.GetBlock(i)
		runes = append(runes, block.GetHtml()...)
	}
	runes = append(runes, BLOCKQUOTE_CLOSE...)
	return
}
