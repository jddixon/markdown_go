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

var (
	SEP_CHAR    = []rune{CR, LF}
	FOUR_SPACES = []rune("    ")

	OPEN_EM      = []rune("<em>")
	CLOSE_EM     = []rune("</em>")
	H_RULE       = []rune("<hr />")
	OPEN_STRONG  = []rune("<strong>")
	CLOSE_STRONG = []rune("</strong>")
)

type Parser struct {
	lexer            *gl.LexInput
	state            State
	lineSep          *LineSep
	crCount, lfCount int
	maybes           []rune
	nonSeps          []rune
	seps             []rune
	downers          []MarkdownI

	// for handling emphasis
	emphChar    rune
	emphDoubled bool

	// our little dictionary
	dict map[string]*Definition
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
		ch           rune
		err          error
		leadingSpace bool

		// header handling
		collected bool
		hashCount int
	)
	lx := p.lexer
	ch, err = lx.NextCh()
	for err == nil {
		if err != nil {
			break
		}
		if p.state == START {
			fmt.Printf("START: '%c' = %d\n", ch, ch) // DEBUG
			if len(p.nonSeps) == 0 {
				if ch == SPACE { // leading tab?
					leadingSpace = true
					goto NEXT
				} else if ch == '#' && !leadingSpace {
					collected, hashCount, err = p.collectHeader()
					if collected {
						p.state = START // XXX
					}
					_ = hashCount
					goto NEXT
				}
			}
			if ch == BACKSLASH {
				var nextChar rune
				p.nonSeps = append(p.nonSeps, BACKSLASH)
				nextChar, err = lx.PeekCh()
				// DEBUG
				fmt.Printf("Parser:START sees BACKSLASH + '%c'\n", nextChar)
				// END
				if escaped(nextChar) {
					ch, err = lx.NextCh()
					p.nonSeps = append(p.nonSeps, ch)
					// DEBUG
				} else {
					fmt.Printf("    %c does not get escaped\n", nextChar)
					// END
				}
				p.state = NONSEP_COLL
			} else if ch == '_' || ch == '*' {
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
			fmt.Printf("SEP_COLL: %d\n", ch) // DEBUG
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
				p.downers = append(p.downers, p.lineSep)
				p.lineSep = nil
				p.nonSeps = append(p.nonSeps, ch)
				p.state = NONSEP_COLL
			}
		} else if p.state == NONSEP_COLL {
			fmt.Printf("NONSEP_COLL: '%c'\n", ch) // DEBUG
			if ch == BACKSLASH {
				var nextChar rune
				p.nonSeps = append(p.nonSeps, BACKSLASH)
				nextChar, err = lx.PeekCh()
				// DEBUG
				fmt.Printf("    sees BACKSLASH + '%c'\n",
					nextChar)
				// END
				if escaped(nextChar) {
					ch, err = lx.NextCh()
					p.nonSeps = append(p.nonSeps, ch)
				}
			} else if len(p.nonSeps) == 0 && ch == SPACE { // leading tab?
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
			if ch == SPACE || ch == TAB {
				// ignore it
			} else if ch == CR || ch == LF {
				p.maybes = append(p.maybes, ch)
				if ch == CR {
					p.crCount++
				} else {
					p.lfCount++
				}
				if p.crCount > 1 || p.lfCount > 1 {
					p.downers = append(p.downers, NewPara(p.nonSeps))
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
				if lastChar == SPACE || lastChar == TAB {
					fmt.Printf("SPACE AT END OF LINE\n") // DEBUG
					if lastChar == TAB {
						fmt.Printf("TAB AT END OF LINE\n") // DEBUG
						p.nonSeps = p.nonSeps[:len(p.nonSeps)-1]
						p.nonSeps = append(p.nonSeps, FOUR_SPACES...)
					}
					p.downers = append(p.downers, NewPara(p.nonSeps))
					p.nonSeps = p.nonSeps[:0]
					p.lineSep, _ = NewLineSep(p.maybes)
					p.downers = append(p.downers, p.lineSep)
					p.maybes = p.maybes[:0]
				} else {
					p.nonSeps = append(p.nonSeps, p.maybes...)
					p.maybes = p.maybes[:0]
				}
				lx.PushBack(ch)
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
			if lastChar == TAB {
				fmt.Printf("TAB AT END OF LINE\n") // DEBUG
				p.nonSeps = p.nonSeps[:len(p.nonSeps)-1]
				p.nonSeps = append(p.nonSeps, FOUR_SPACES...)
			}
			p.downers = append(p.downers, NewPara(p.nonSeps))
			p.nonSeps = p.nonSeps[:0]
		}
		err = nil
	}
	if err == nil && len(p.nonSeps) > 0 {
		p.downers = append(p.downers, NewPara(p.nonSeps))
	}
	return p.downers, err
}
