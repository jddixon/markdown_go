package md

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
		if ch == '\r' || ch == '\n' {
			fmt.Printf("    got CR|LF while collecting emph\n") // DEBUG

			// we didn't collect the whole thing so just push it all back
			lx.PushBack(emphChar)
			if p.emphDoubled {
				lx.PushBack(emphChar)
			}
			lx.PushBackChars(runes)
			lx.PushBack(ch)
			// DEBUG
			fmt.Printf("    EOL in collectEmph: PUSHING BACK '%s'\n",
				string(runes))
			// END
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
