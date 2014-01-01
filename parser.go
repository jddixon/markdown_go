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
	} else {
		line = new(Line)
	}
	line.Err = err

	return
}

func (p *Parser) Parse() (doc *Document, err error) {

	var (
		imageDefn        *Definition
		linkDefn         *Definition
		curPara          *Para
		q                *Line
		ch0              rune
		lastBlockLineSep bool
	)
	docPtr := p.doc

	q = p.readLine()
	err = q.Err

	// DEBUG
	fmt.Printf("Parse: first line is '%s'\n", string(q.runes))
	// END

	// pass through the document line by line
	for err == nil || err == io.EOF {
		if len(q.runes) > 0 {

			// HANDLE DEFINITIONS -----------------------------------

			// rigidly require that definitions start in the first column
			if q.runes[0] == '[' { // possible link definition
				linkDefn, err = q.parseLinkDefinition(docPtr)
			}
			if err == nil && linkDefn == nil && q.runes[0] == '!' {
				imageDefn, err = q.parseImageDefinition(docPtr)
			}
			// HANDLE BLOCKS ----------------------------------------

			if (err == nil || err == io.EOF) && linkDefn == nil && imageDefn == nil {
				var b BlockI
				ch0 = q.runes[0]
				eol := len(q.runes)

				// HEADERS --------------------------------
				if ch0 == '#' {
					b, err = q.parseHeader()
				}

				// HORIZONTAL RULES ----------------------
				if b == nil && (err == nil || err == io.EOF) &&
					(ch0 == '-' || ch0 == '*' || ch0 == '_') {
					b, err = q.parseHRule()
				}

				// XXX STUB : TRY OTHER PARSERS

				// BLOCKQUOTE -----------------------------
				if b == nil && (err == nil || err == io.EOF) && ch0 == '>' {
					b, err = q.parseBlockquote(doc, 1)
				}
				// ORDERED LISTS --------------------------

				// XXX We require a space after these starting characters
				if b == nil && (err == nil || err == io.EOF) {
					var from int
					for from = 0; from < 3 && from < eol; from++ {
						if !u.IsSpace(q.runes[from]) {
							break
						}
					}
					if from < eol-2 {

						// we are positioned on a non-space character
						ch0 := q.runes[from]
						ch1 := q.runes[from+1]
						ch2 := q.runes[from+2]
						if u.IsDigit(ch0) && ch1 == '.' && ch2 == ' ' {
							b, err = q.parseOrdered(from + 2)

						}
					}
				}

				// UNORDERED LISTS ------------------------

				// XXX We require a space after these starting characters
				if b == nil && (err == nil || err == io.EOF) {
					var from int
					for from = 0; from < 3 && from < eol; from++ {
						if !u.IsSpace(q.runes[from]) {
							break
						}
					}
					if from < eol-1 {
						// we are positioned on a non-space character
						ch0 := q.runes[from]
						ch1 := q.runes[from+1]
						if (ch0 == '*' || ch0 == '+' || ch0 == '-') && ch1 == ' ' {
							b, err = q.parseUnordered(from + 2)
						}
					}
				} // GEEP

				// DEFAULT: PARA --------------------------
				// If we have parsed the line, we hang the block off
				// the document.  Otherwise, we treat whatever we have
				// as a sequence of spans and make a Para out of it.
				if err == nil || err == io.EOF {
					if b != nil {
						docPtr.AddChild(b)
						lastBlockLineSep = false
					} else {
						// default parser
						// DEBUG
						fmt.Printf("== invoking parseSpanSeq(true) ==\n")
						// END
						var seq *SpanSeq
						seq, err = q.parseSpanSeq(docPtr, 0, true)
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
			ls, err := NewLineSep(q.lineSep)
			if err == nil {
				if curPara != nil {
					docPtr.AddChild(curPara)
					curPara = nil
					lastBlockLineSep = false
				}
				fmt.Printf("adding LineSep to document\n") // DEBUG
				if !lastBlockLineSep {
					docPtr.AddChild(ls)
					lastBlockLineSep = true
				}
			}
		}
		if err != nil {
			break
		}
		q = p.readLine()
		err = q.Err
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
			docPtr.AddChild(curPara)
			curPara = nil
		}
		// DEBUG
		fmt.Printf("returning thisDoc with %d children\n", len(docPtr.children))
		// END
		doc = docPtr
	}
	return
}
