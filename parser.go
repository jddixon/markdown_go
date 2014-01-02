package md

// xgo/md/parser.go

import (
	"fmt"
	gl "github.com/jddixon/xgo/lex"
	"io"
	u "unicode"
)

var _ = fmt.Print

const (
	OK int = iota
	ACK
	DONE
)

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
	// DEBUG
	if err != nil {
		fmt.Printf("Parser.readLine(): err = %s\n", err.Error())
	}
	// END
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
		eofSeen    bool
		lastWasDef bool
	)
	doc = p.doc
	out := make(chan *Line)
	resp := make(chan int)
	stop := make(chan bool)

	go ParseHolder(doc, p, out, resp, stop)
	status := <-resp

	q := p.readLine()
	err = q.Err
	if err == io.EOF {
		eofSeen = true
	}
	// DEBUG
	fmt.Printf("Parser: LINE: '%s'\n", string(q.runes))
	if err == nil {
		fmt.Println("    NIL error")
	} else {
		fmt.Printf("    error = '%s'\n", err.Error())
	}
	// END
	for err == nil || err == io.EOF {
		var (
			imageDefn *Definition
			linkDefn  *Definition
		)
		if len(q.runes) > 0 {
			// HANDLE DEFINITIONS -----------------------------------

			// rigidly require that definitions start in the first column
			if q.runes[0] == '[' { // possible link definition
				linkDefn, err = q.parseLinkDefinition(doc)
			}
			if err == nil && linkDefn == nil && q.runes[0] == '!' {
				imageDefn, err = q.parseImageDefinition(doc)
			}
		}
		if imageDefn == nil && linkDefn == nil {
			lastWasDef = false
			out <- q // send-send DEADLOCK
			status = <-resp
		} else {
			lastWasDef = true
		}
		if err == io.EOF || eofSeen {
			break
		}
		q = p.readLine()
		err = q.Err
		if err == io.EOF {
			eofSeen = true
		}
	}
	if lastWasDef {
		stop <- true
	}
	status = <-resp
	_ = status // UNUSED
	if err == nil && eofSeen {
		err = io.EOF
	}
	return
}
