package md

// xgo/md/Parser.go

import (
	"fmt"
	gl "github.com/jddixon/xgo/lex"
	"io"
)

var _ = fmt.Print

type State int

const (
	START State = iota
	NONSEP_COLL
	MAYBE_COLL
	SEP_COLL
)
const (
	CR = '\r'
	LF = '\n'
)

var (
	SEP_CHAR = []rune{'\r', '\n'}
)

func Parse(reader io.Reader) (bits []MarkdownI, err error) {

	lx, err := gl.NewNewLexInput(reader)
	if err == nil {
		var (
			state            State
			ch               rune
			lineSep          *LineSep
			para             *Para
			crCount, lfCount int
			maybes           []rune
			nonSeps          []rune
			seps             []rune
		)
		for err == nil {
			ch, err = lx.NextCh()
			if err != nil {
				break
			}
			if state == START {
				if ch == CR || ch == LF {
					seps = append(seps, ch)
					state = SEP_COLL
				} else {
					nonSeps = append(nonSeps, ch)
					state = NONSEP_COLL
				}
			} else if state == SEP_COLL {
				if ch == CR || ch == LF {
					seps = append(seps, ch)
					// state unchanged
				} else {
					lineSep, err = NewLineSep(seps)
					seps = seps[:0]
					bits = append(bits, lineSep)
					lineSep = nil
					nonSeps = append(nonSeps, ch)
					state = NONSEP_COLL
				}
			} else if state == NONSEP_COLL {
				if ch == CR || ch == LF {
					maybes = append(maybes, ch)
					state = MAYBE_COLL
				} else {
					nonSeps = append(nonSeps, ch)
					// state unchanged
				}
			} else if state == MAYBE_COLL {
				if ch == CR || ch == LF {
					maybes = append(maybes, ch)
					if ch == CR {
						crCount++
					} else {
						lfCount++
					}
					if crCount > 1 || lfCount > 1 {
						bits = append(bits, NewPara(nonSeps))
						para = nil
						seps = make([]rune, len(maybes))
						copy(seps, maybes)
						maybes = maybes[:0]
						state = SEP_COLL
					}
				} else {
					nonSeps = append(nonSeps, maybes...)
					maybes = maybes[:0]
				}
			}
		}
		if err == io.EOF {
			// XXX STUB, ACTION DEPENDS UPON STATE
			if state == SEP_COLL {
				seps = seps[:0] // just discard
			} else if state == NONSEP_COLL || state == MAYBE_COLL {
				bits = append(bits, NewPara(nonSeps))
			}
			err = nil
		}
		// ALMOST CERTAINLY WRONG
		if err == nil && len(nonSeps) > 0 {
			bits = append(bits, NewPara(nonSeps))
		}
	}
	return
}
