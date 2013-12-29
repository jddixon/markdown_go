package md

// xgo/md/blockquote.go

import (
	"fmt"
)

var _ = fmt.Print

// XXX THIS IS JUST A HACK FROM PARA, WILL CERTAINLY NEED REVISION

// A Blockquote is a Block consisting of a number of SpanSeq, each of which is
// a sequence of SpanI, things that implement the SpanI interface.  In
// addition, each SpanSeq has an associated lineSeq, a slice of line
// terminators (null, CR, or LF)
type Blockquote struct {
	seqs []SpanSeq
}

func (bq *Blockquote) Add(seq SpanSeq) (err error) {
	bq.seqs = append(bq.seqs, seq)
	return
}

func (bq *Blockquote) Get() (runes []rune) {
	runes = append(runes, PARA_OPEN...)
	for i := 0; i < len(bq.seqs); i++ {
		spans := bq.seqs[i].spans
		for j := 0; j < len(spans); j++ {
			runes = append(runes, spans[j].Get()...)
		}
		if i < len(bq.seqs)-1 {
			runes = append(runes, bq.seqs[i].lineSep...)
		}
	}
	runes = append(runes, PARA_CLOSE...)
	runes = append(runes, bq.seqs[len(bq.seqs)-1].lineSep...)
	return
}

// // Parse a line known to begin a new Blockquote.
// func (q *Line) parseBlockquote() (b BlockI, err error) {
// 	var bq *Blockquote
// 	spans, err := q.parseSpanSeq()
// 	if err == nil {
// 		// DEBUG
// 		fmt.Printf("parseBlockquote() finds %d spans\n", len(spans))
// 		// END
// 		pa = new(Blockquote)
// 		for i := 0; i < len(spans); i++ {
// 			err = pa.Add(spans[i])
// 			if err != nil {
// 				break
// 			}
// 		}
// 		if err == nil {
// 			// DEBUG
// 			fmt.Printf("  parseBlockquote() returning Blockquote with %d spans\n",
// 				len(pa.spans))
// 			// END
// 			b = pa
// 		}
// 	}
// 	return
// }
