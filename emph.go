package md

// xgo/md/emph.go

import (
	"fmt"
	"io"
)

func (p *Parser) collectEmph(emphChar rune) (collected bool, err error) {

	fmt.Printf("Entering collectEmph(%c)\n", emphChar)

	p.emphChar = emphChar
	p.emphDoubled = false
	firstChar := true
	lx := p.lexer

	var (
		atEOF               bool
		haveSeenFirstCloser bool
		runes               []rune
	)
	ch, err := lx.NextCh()
	for err == nil {
		if firstChar {
			firstChar = false
			if ch == emphChar {
				fmt.Printf("first char %c, emphasis doubled\n", ch) // DEBUG
				p.emphDoubled = true
				goto NEXT
			}
		}
		if atEOF || ch == '\r' || ch == '\n' {
			// DEBUG
			if ch == CR || ch == LF {
				fmt.Printf("    got CR|LF while collecting emph\n") // DEBUG
			} else {
				fmt.Printf("    collectEmph; at EOF;  PUSHING BACK '%s'\n",
					string(runes))
			}
			// we didn't collect the whole thing so forward the emphChar(s)
			// and push back the rest, including the current character, so
			// long as it's not a null byte
			p.nonSeps = append(p.nonSeps, emphChar)
			if p.emphDoubled {
				p.nonSeps = append(p.nonSeps, emphChar)
			}
			lx.PushBackChars(runes)
			if ch != rune(0) {
				lx.PushBack(ch)
			}
			collected = false
			break
		} else if ch == emphChar {
			fmt.Printf("    got emph char %c while collecting\n", ch)
			if p.emphDoubled {
				if !haveSeenFirstCloser {
					fmt.Printf("first of two closers seen\n") // DEBUG
					haveSeenFirstCloser = true
					goto NEXT
				} else {
					fmt.Printf("closing double-emph\n") // DEBUG
					p.nonSeps = append(p.nonSeps, OPEN_STRONG...)
					p.nonSeps = append(p.nonSeps, runes...)
					p.nonSeps = append(p.nonSeps, CLOSE_STRONG...)
				}
			} else {
				fmt.Printf("closing single-emph\n")
				p.nonSeps = append(p.nonSeps, OPEN_EM...)
				p.nonSeps = append(p.nonSeps, runes...)
				p.nonSeps = append(p.nonSeps, CLOSE_EM...)
			}
			fmt.Println("COLLECTED") // DEBUG
			collected = true
			break
		} else {
			runes = append(runes, ch)
		}

		if atEOF {
			break
		}
	NEXT:
		ch, err = lx.NextCh()
		if err == io.EOF {
			err = nil
			atEOF = true
		}
	}
	return
}
