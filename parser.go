package md

// xgo/md/parser.go

import (
	"fmt"
	gl "github.com/jddixon/xgo/lex"
	"io"
	u "unicode"
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

	var (
		allSpaces bool = true // if a line is all spaces, we ignore them
		atEOF     bool
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
				var ch2 rune
				ch2, err = lx.PeekCh()
				if err == io.EOF {
					err = nil
				}
				if err == nil && ch2 == LF {
					ch2, _ = lx.NextCh()
					thisLine.lineSep = append(thisLine.lineSep, ch2)
				}
			}
			if !allSpaces {
				// DEBUG
				fmt.Printf("LINE: '%s'\n", string(runes))
				// END
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
	}
	return
}

func (p *Parser) Parse() (doc *Document, err error) {
	var (
		imageDefn        *Definition
		linkDefn         *Definition
		curPara          *Para
		q                *Line
		ch1              rune
		lastBlockLineSep bool
	)
	docPtr := p.doc

	q, err = p.readLine()

	// DEBUG
	fmt.Printf("Parse: first line is '%s'\n", string(q.runes))
	// END

	// pass through the document line by line
	for err == nil || err == io.EOF {
		if len(q.runes) > 0 {

			// rigidly require that definitions start in the first column
			if q.runes[0] == '[' { // possible link definition
				linkDefn, err = q.parseLinkDefinition(docPtr)
			}
			if err == nil && linkDefn == nil && q.runes[0] == '!' {
				imageDefn, err = q.parseImageDefinition(docPtr)
			}
			if (err == nil || err == io.EOF) && linkDefn == nil && imageDefn == nil {
				var b BlockI
				ch1 = q.runes[0]

				if ch1 == '#' {
					b, err = q.parseHeader()
				}
				if b == nil && (err == nil || err == io.EOF) &&
					(ch1 == '-' || ch1 == '*' || ch1 == '_') {
					b, err = q.parseHRule()
				}

				// XXX STUB : TRY OTHER PARSERS

				if err == nil || err == io.EOF {
					if b != nil {
						docPtr.addBlock(b)
						lastBlockLineSep = false
					} else {
						// default parser
						// DEBUG
						fmt.Printf("== invoking parseSpanSeq(true) ==\n")
						// END
						var seq *SpanSeq
						seq, err = q.parseSpanSeq(docPtr, true)
						if err == nil || err == io.EOF {
							if curPara == nil {
								curPara = new(Para)
							}
							fmt.Printf("* adding seq to curPara\n") // DEBUG
							curPara.seqs = append(curPara.seqs, *seq)
							fmt.Printf("  curPara has %d seqs\n",
								len(curPara.seqs))
						}
					}
				}
			}

		} else {
			// we got a blank line
			if !lastBlockLineSep {
				ls, err := NewLineSep(q.lineSep)
				if err == nil {
					if curPara != nil {
						docPtr.addBlock(curPara)
						curPara = nil
					}
					fmt.Printf("adding LineSep to document\n") // DEBUG
					docPtr.addBlock(ls)
				}
			}
			lastBlockLineSep = true
		}
		if err != nil {
			break
		}
		q, err = p.readLine()
		if (err != nil && err != io.EOF) || q == nil {
			break
		}
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
			fmt.Println("have dangling curPara") // DEBUG
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
