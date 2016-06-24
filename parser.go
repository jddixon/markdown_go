package md

// xgo/md/parser.go

import (
	"fmt"
	gl "github.com/jddixon/xgo_go/lex"
	"io"
	u "unicode"
)

var _ = fmt.Print

const (
	OK int = 1 << iota
	ACK
	DONE
	LAST_LINE_PROCESSED
)

type Parser struct {
	lexer *gl.LexInput
	doc   *Document
	opt   *Options
}

func NewParser(opt *Options) (p *Parser, err error) {

	if opt == nil {
		err = NilOptions
	} else {
		var doc *Document
		reader := opt.Reader
		lx, err := gl.NewNewLexInput(reader)
		if err == nil {
			doc, err = NewDocument(opt)
		}
		if err == nil {
			p = &Parser{
				lexer: lx,
				doc:   doc,
				opt:   opt,
			}
		}
	}
	return
}

func (p *Parser) GetDocument() *Document {
	return p.doc
}

func (p *Parser) readLine() (line *Line) {

	var (
		allSpaces bool = true // if a line is all spaces, we ignore them
		atEOF     bool
		err       error
		runes     []rune
		thisLine  Line
	)

	lx := p.lexer
	ch, err := lx.NextCh()
	if err == io.EOF {
		err = nil
		atEOF = true
	}
	for err == nil {
		if ch == CR || ch == LF || ch == rune(0) {
			thisLine.lineSep = append(thisLine.lineSep, ch)
			if ch == CR {
				var ch1 rune
				ch1, err = lx.PeekCh()
				if err == io.EOF {
					err = nil
				}
				if err == nil && ch1 == LF {
					ch1, _ = lx.NextCh()
					thisLine.lineSep = append(thisLine.lineSep, ch1)
				}
			}
			if !allSpaces {
				thisLine.runes = runes
			}
			break // eol has been seen
		}
		if !u.IsSpace(ch) {
			allSpaces = false
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
		line = &thisLine
		if atEOF {
			err = io.EOF
		}
	} else {
		line = new(Line)
	}
	line.Err = err
	return
}

func (p *Parser) Parse() (doc *Document, err error) {

	var (
		status int
	)
	doc = p.doc
	q := p.readLine()
	out, status := doc.ParseHolder(p, q)
	err = out.Err
	if p.opt.Testing {
		fmt.Printf("Parser: LINE: '%s'\n", string(out.runes))
		if err != nil {
			fmt.Printf("    error = '%s'\n", err.Error())
		}
	}

	_ = status // UNUSED
	return
}
