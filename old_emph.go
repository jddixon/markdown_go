package md

// xgo/md/old_emph.go

import (
	"fmt"
	"io"
)

// XXX CHANGE: Attempt to parse out an EmphSpan, returning a SpanI reference
// to it on success and nil and possibly an error on failure.  If the parse
// fails but there is no input error, push back all characters seen on the
// input.  If the parse has succeeded, the SpanI will added to the current
// Para's list of SpanIs.
//
func (p *OldParser) oldParseEmph(emphChar rune) (collected bool, err error) {

	fmt.Printf("Entering parseEmph(%c)\n", emphChar)

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
				fmt.Printf("    parseEmph; at EOF;  PUSHING BACK '%s'\n",
					string(runes))
			}
			// we didn't collect the whole thing so forward the emphChar(s)
			// and push back the rest, including the current character, so
			// long as it's not a null byte
			p.curText = append(p.curText, emphChar)
			if p.emphDoubled {
				p.curText = append(p.curText, emphChar)
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
					p.curText = append(p.curText, OPEN_STRONG...)
					p.curText = append(p.curText, runes...)
					p.curText = append(p.curText, CLOSE_STRONG...)
				}
			} else {
				fmt.Printf("closing single-emph\n")
				p.curText = append(p.curText, OPEN_EM...)
				p.curText = append(p.curText, runes...)
				p.curText = append(p.curText, CLOSE_EM...)
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
