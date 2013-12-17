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
	SEP_CHAR    = []rune{'\r', '\n'}
	FOUR_SPACES = []rune{' ', ' ', ' ', ' '}

	OPEN_EM      = []rune{'<', 'e', 'm', '>'}
	CLOSE_EM     = []rune{'<', '/', 'e', 'm', '>'}
	H_RULE       = []rune{'<', 'h', 'r', ' ', '/', '>'}
	OPEN_STRONG  = []rune{'<', 's', 't', 'r', 'o', 'n', 'g', '>'}
	CLOSE_STRONG = []rune{'<', '/', 's', 't', 'r', 'o', 'n', 'g', '>'}
)

type Parser struct {
	lexer            *gl.LexInput
	state            State
	lineSep          *LineSep
	crCount, lfCount int
	maybes           []rune
	nonSeps          []rune
	seps             []rune
	bits             []MarkdownI

	// for handling emphasis
	emphChar    rune
	emphDoubled bool
}

func NewParser(reader io.Reader) (p *Parser, err error) {
	lx, err := gl.NewNewLexInput(reader)
	if err == nil {
		p = &Parser{
			lexer: lx,
			state: START,
		}
	}
	return
}

func (p *Parser) Parse() ([]MarkdownI, error) {
	var (
		ch  rune
		err error
	)
	lx := p.lexer
	ch, err = lx.NextCh()
	for err == nil {
		if err != nil {
			break
		}
		if p.state == START {
			if len(p.nonSeps) == 0 {
				if ch == ' ' { // leading tab?
					// ignore
					goto NEXT
				} else if ch == '#' {
					p.collectHeader()
					goto NEXT
				}
			}
			if ch == '_' || ch == '*' {
				// scan ahead for matching ch; if we find it, we
				// wrap it properly and append it to p.nonSeps and
				// return true.  Otherwise we push whatever has been
				// scanned back on input, append ch to p.nonSeps,
				// and return false.
				p.collectEmph(ch)
			} else if ch == CR || ch == LF {
				if ch == CR {
					p.crCount++
				} else {
					p.lfCount++
				}
				p.seps = append(p.seps, ch)
				p.state = SEP_COLL
			} else {
				p.nonSeps = append(p.nonSeps, ch)
				p.state = NONSEP_COLL
			}
		} else if p.state == SEP_COLL {
			if ch == CR || ch == LF {
				if ch == CR {
					p.crCount++
					if p.crCount < 3 {
						p.seps = append(p.seps, ch)
					}
				} else {
					p.lfCount++
					if p.lfCount < 3 {
						p.seps = append(p.seps, ch)
					}
				}
				// p.state unchanged
			} else {
				p.lineSep, err = NewLineSep(p.seps)
				p.seps = p.seps[:0]
				p.crCount = 0
				p.lfCount = 0
				p.bits = append(p.bits, p.lineSep)
				p.lineSep = nil
				p.nonSeps = append(p.nonSeps, ch)
				p.state = NONSEP_COLL
			}
		} else if p.state == NONSEP_COLL {
			if len(p.nonSeps) == 0 && ch == ' ' { // leading tab?
				// ignore
			} else if ch == '_' || ch == '*' {
				// scan ahead for matching ch; if we find it, we
				// wrap it properly and append it to p.nonSeps and
				// return true.  Otherwise we push whatever has been
				// scanned back on input, append ch to p.nonSeps,
				// and return false.
				p.collectEmph(ch)
			} else if ch == CR || ch == LF {
				if ch == CR {
					p.crCount = 1
				} else {
					p.lfCount = 1
				}
				p.maybes = append(p.maybes, ch)
				p.state = MAYBE_COLL
			} else {
				p.nonSeps = append(p.nonSeps, ch)
				// p.state unchanged
			}
		} else if p.state == MAYBE_COLL {
			if ch == ' ' || ch == '\t' {
				// ignore it
			} else if ch == CR || ch == LF {
				p.maybes = append(p.maybes, ch)
				if ch == CR {
					p.crCount++
				} else {
					p.lfCount++
				}
				if p.crCount > 1 || p.lfCount > 1 {
					p.bits = append(p.bits, NewPara(p.nonSeps))
					p.nonSeps = p.nonSeps[:0]
					p.seps = make([]rune, len(p.maybes))
					copy(p.seps, p.maybes)
					p.maybes = p.maybes[:0]
					p.state = SEP_COLL
				}
			} else {
				// If the last nonSep is a space (or tab?) we
				// make the nonSep a para, insert a p.lineSep,
				// and start a new para.
				lastChar := p.nonSeps[len(p.nonSeps)-1]
				if lastChar == ' ' || lastChar == '\t' {
					fmt.Printf("SPACE AT END OF LINE\n") // DEBUG
					if lastChar == '\t' {
						fmt.Printf("TAB AT END OF LINE\n") // DEBUG
						p.nonSeps = p.nonSeps[:len(p.nonSeps)-1]
						p.nonSeps = append(p.nonSeps, FOUR_SPACES...)
					}
					p.bits = append(p.bits, NewPara(p.nonSeps))
					p.nonSeps = p.nonSeps[:0]
					p.lineSep, _ = NewLineSep(p.maybes)
					p.bits = append(p.bits, p.lineSep)
					p.maybes = p.maybes[:0]
				} else {
					p.nonSeps = append(p.nonSeps, p.maybes...)
					p.maybes = p.maybes[:0]
				}
				p.nonSeps = append(p.nonSeps, ch)
				p.state = NONSEP_COLL
			}
		}
	NEXT:
		ch, err = lx.NextCh()
	}
	if err == io.EOF {
		if p.state == SEP_COLL {
			p.seps = p.seps[:0] // just discard
		} else if p.state == NONSEP_COLL || p.state == MAYBE_COLL {
			lastChar := p.nonSeps[len(p.nonSeps)-1]
			if lastChar == '\t' {
				fmt.Printf("TAB AT END OF LINE\n") // DEBUG
				p.nonSeps = p.nonSeps[:len(p.nonSeps)-1]
				p.nonSeps = append(p.nonSeps, FOUR_SPACES...)
			}
			p.bits = append(p.bits, NewPara(p.nonSeps))
			p.nonSeps = p.nonSeps[:0]
		}
		err = nil
	}
	if err == nil && len(p.nonSeps) > 0 {
		p.bits = append(p.bits, NewPara(p.nonSeps))
	}
	return p.bits, err
}
