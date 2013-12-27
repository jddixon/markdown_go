package md

// xgo/md/parser.go

import (
	"fmt"
	gl "github.com/jddixon/xgo/lex"
	"io"
)

var _ = fmt.Print

type Parser struct {
	lexer *gl.LexInput
	doc   *Document
}

func NewParser(reader io.Reader) (p *Parser, err error) {

	var doc *Document
	lx, err := gl.NewNewLexInput(reader)
	if err == nil {
		doc, err = NewDocument()
	}
	if err == nil {
		p = &Parser{
			lexer: lx,
			doc:   doc,
		}
	}
	return
}

func (p *Parser) readLine() (line *Line, err error) {

	var thisLine Line
	var runes []rune
	var atEOF bool

	lx := p.lexer
	ch, err := lx.NextCh()
	for err == nil {
		if ch == CR || ch == LF || ch == rune(0) {
			thisLine.lineSep = append(thisLine.lineSep, ch)
			thisLine.runes = runes
			break
		}
		runes = append(runes, ch)
		if atEOF {
			break
		}
		ch, err = lx.NextCh()
		if err == io.EOF {
			err = nil
			atEOF = true
		}
	}
	if err == nil {
		thisLine.runes = runes
		line = &thisLine
		if atEOF {
			err = io.EOF
		}
	}
	return
}

func (p *Parser) Parse() (doc *Document, err error) {
	var (
		imageDefn *Definition
		linkDefn  *Definition
		curPara   *Para
		q         *Line
		thisDoc   Document
	)
	docPtr := &thisDoc

	q, err = p.readLine()

	// DEBUG
	fmt.Printf("Parse: first line is '%s'\n", string(q.runes))
	// END

	// pass through the document line by line
	for err == nil {
		if len(q.runes) > 0 {

			// rigidly require that definitions start in the first column
			if q.runes[0] == '[' { // possible link definition
				linkDefn, err = q.parseLinkDefinition(docPtr)
			}
			if err == nil && linkDefn == nil && q.runes[0] == '!' {
				imageDefn, err = q.parseImageDefinition(docPtr)
			}
			if err == nil && linkDefn == nil && imageDefn == nil {
				var b BlockI

				_ = b

				// XXX STUB : DO GOOD THINGS

				// DEBUG
				fmt.Printf("invoking parseSpanSeq()\n")
				// END
				var seq *SpanSeq
				seq, err = q.parseSpanSeq()
				if err == nil {
					if curPara == nil {
						curPara = new(Para)
					}
					curPara.seqs = append(curPara.seqs, *seq)
				}
			}

		} else {
			// we got a blank line
			ls, err := NewLineSep(q.lineSep)
			if err == nil {
				if curPara != nil {
					docPtr.addBlock(curPara)
					curPara = nil
				}
				docPtr.addBlock(ls)
			}
		}
		if err != nil {
			break
		}

		q, err = p.readLine()
		if len(q.runes) == 0 {
			fmt.Printf("ZERO-LENGTH LINE")
			if len(q.lineSep) == 0 && q.lineSep[0] == rune(0) {
				break
			}
			fmt.Printf("  lineSep is 0x%x\n", q.lineSep[0])
		}
		// DEBUG
		fmt.Printf("Parse: next line is '%s'\n", string(q.runes))
		// END
	}
	if err == nil || err == io.EOF {
		if curPara != nil {
			docPtr.addBlock(curPara)
			curPara = nil
		}
		// DEBUG
		fmt.Printf("returning thisDoc with %d blocks\n", len(docPtr.blocks))
		// END
		doc = docPtr
	}
	return
}
