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
		line *Line
	)
	line, err = p.readLine()
	for err == nil {

		// DO GOOD THINGS

		line, err = p.readLine()
	}
	// XXX STUB

	_ = line // DEBUG
	return
}
