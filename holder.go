package md

import (
	"fmt"
	"io"
	u "unicode"
)

type Holder struct {
	children []BlockI
}

func NewHolder() *Holder {
	var h = new(Holder)
	return h
}

func (h *Holder) AddChild(child BlockI) (err error) {
	if child == nil {
		err = NilChild
	} else {
		// XXX We don't prevent duplicates
		h.children = append(h.children, child)
	}
	return
}

func (h *Holder) Size() int {
	return len(h.children)
}

func (h *Holder) GetChild(n int) (child BlockI, err error) {
	if n < 0 || h.Size() <= n {
		err = ChildNdxOutOfRange
	} else {
		child = h.children[n]
	}
	return
}

func ParseHolder(holder HolderI, p *Parser,
	in chan *Line, resp chan int, stop chan bool) {

	doc := holder.(*Document)
	var (
		eofSeen          bool
		err              error
		curPara          *Para
		q                *Line
		ch0              rune
		lastBlockLineSep bool
		stopped          bool
	)
	resp <- OK // OK, setup complete

	q = <-in // WAS q = p.readLine()
	err = q.Err
	if err == io.EOF {
		eofSeen = true
	}
	resp <- ACK // MOVE ME

	// DEBUG
	fmt.Printf("ParseHolder: first line is '%s'\n", string(q.runes))
	if err == nil {
		fmt.Println("    nil error")
	} else {
		fmt.Printf("    error = %s\n", err.Error())
	}
	// END

	// pass through the document line by line
	for err == nil || err == io.EOF {
		if len(q.runes) > 0 {

			// HANDLE BLOCKS ----------------------------------------

			if err == nil || err == io.EOF {
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
						doc.AddChild(b)
						lastBlockLineSep = false
					} else {
						// default parser
						// DEBUG
						fmt.Printf("== invoking parseSpanSeq(true) ==\n")
						// END
						var seq *SpanSeq
						seq, err = q.parseSpanSeq(doc, 0, true)
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
					doc.AddChild(curPara)
					curPara = nil
					lastBlockLineSep = false
				}
				fmt.Printf("adding LineSep to document\n") // DEBUG
				if !lastBlockLineSep {
					doc.AddChild(ls)
					lastBlockLineSep = true
				}
			}
		}
		if err != nil || eofSeen {
			// DEBUG
			fmt.Println("parseHolder breaking, error or EOF seen")
			if err != nil {
				fmt.Printf("    ERROR: %s\n", err.Error())
			}
			if eofSeen {
				fmt.Println("    EOF SEEN, so breaking")
			}
			// END
			break
		}
		var ok bool
		select {
		case q, ok = <-in:
			if ok {
				err = q.Err
			} else {
				goto JUST_DIE
			}
		case stopped, ok = <-stop:
			if !ok {
				goto JUST_DIE
			} else {
				_ = stopped // not yet used - and so far unnecessary
				goto SAYOONARA
			}
		}
		// DEBUG
		if err == io.EOF {
			eofSeen = true
			fmt.Println("*** EOF SEEN ***")
		}
		// END
		resp <- ACK // MOVE ME
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
			doc.AddChild(curPara)
			curPara = nil
		}
		// DEBUG
		fmt.Printf("returning thisDoc with %d children\n", len(doc.children))
		// END
	}

SAYOONARA:
	resp <- DONE // DEADLOCK send-send
	return

JUST_DIE:
}
