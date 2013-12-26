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
		q         *Line
	)
	q, err = p.readLine()

	// DEBUG
	fmt.Printf("Parse: first line is '%s'\n", string(q.runes))
	// END

	// pass through the document line by line
	for err == nil {

		// rigidly require that definitions start in the first column
		if q.runes[0] == '[' { // possible link definition
			linkDefn, err = q.parseLinkDefinition(doc)
		} else if err == nil && linkDefn == nil {
			imageDefn, err = q.parseImageDefinition(doc)
		} else if err == nil && imageDefn == nil {

			// XXX STUB : DO GOOD THINGS

		}
		if err != nil {
			break
		}

		q, err = p.readLine()
		// DEBUG
		fmt.Printf("Parse: next line is '%s'\n", string(q.runes))
		// END
	}
	return
}
