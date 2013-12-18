package md

// xgo/md/ampersand.go

import (
	"fmt"
	"io"
	u "unicode"
)

// We have seen an ampersand.  If it begins an entity, then collect
// the entity and append it to p.nonSeps.  Otherwise, append the
// ampersand to p.nonSeps and push back anything read.
func (p *Parser) handleAmpersand() (collected bool, err error) {

	fmt.Printf("Entering handleAmpersand()\n")

	lx := p.lexer
	const (
		MAX_ENTITY_CHAR = 6 // between & and ;
	)
	var (
		atEOF bool
		runes []rune
	)
	ch, err := lx.NextCh()
	for err == nil {
		runes = append(runes, ch)
		if u.IsSpace(ch) || ch == '\r' || ch == '\n' {
			p.nonSeps = append(p.nonSeps, '&') // output the &
			lx.PushBackChars(runes)
			collected = false
			break
		} else if ch == ';' {
			fmt.Printf("got emph char %c while collecting\n", ch)

			// WORKING HERE

			fmt.Println("COLLECTED") // DEBUG
			collected = true
			break
		}
		// AND HERE

		if atEOF {
			break
		}
		ch, err = lx.NextCh()
		if err == io.EOF {
			err = nil
			atEOF = true
		}
	}
	return
}
