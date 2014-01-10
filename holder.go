package md

import (
	"fmt"
	"io"
	u "unicode"
)

// A holder is a syntactic structure, a collection of BlockIs, things with
// a BlockI interface.
//
// Remember that a top level holder has definitions and is called a Document.
type Holder struct {
	isBlockquote bool
	depth        uint
	blocks       []BlockI
}

func NewHolder(isBq bool, depth uint) (h *Holder, err error) {
	if depth > 0 && !isBq {
		err = OnlyBlockquoteSupported
	} else {
		h = &Holder{
			isBlockquote: isBq,
			depth:        depth,
		}
	}
	return
}

func (h *Holder) AddBlock(block BlockI) (err error) {
	if block == nil {
		err = NilChild
	} else {
		// XXX We don't prevent duplicates
		h.blocks = append(h.blocks, block)
	}
	return
}

func (h *Holder) Size() int {
	return len(h.blocks)
}

func (h *Holder) GetBlock(n int) (block BlockI, err error) {
	if n < 0 || h.Size() <= n {
		err = ChildNdxOutOfRange
	} else {
		block = h.blocks[n]
	}
	return
}

// Return an offset 1 beyond the number of chevrons ('>') expected
// for this depth.  At depth N, we skip N.
func SkipChevrons(q *Line, depth uint) (from uint) {

	var count uint
	eol := uint(len(q.runes))
	for offset := uint(0); offset < eol; offset++ {
		if q.runes[offset] == '>' {
			count++
			if count >= depth {
				from = offset + 1
				break
			}
		}
	}
	return
}

func (h *Holder) ParseHolder(p *Parser,
	in chan *Line, resp chan int, stop chan bool) {

	doc := p.GetDocument()
	var (
		eofSeen          bool
		err              error
		curPara          *Para
		q                *Line
		ch0              rune
		lastBlockLineSep bool
		stopped          bool

		// used to control child holder (for Blockquote)
		haveChild bool
		child     *Blockquote
		toChild   chan *Line
		fromChild chan int
		stopChild chan bool
	)
	resp <- OK // OK, setup complete

	q = <-in // WAS q = p.readLine()
	err = q.Err
	if err == io.EOF {
		eofSeen = true
	}
	resp <- ACK // MOVE ME

	// DEBUG
	fmt.Printf("ParseHolder depth %d: first line is '%s'\n",
		h.depth, string(q.runes))
	if err == nil {
		fmt.Println("    nil error")
	} else {
		fmt.Printf("    error = %s\n", err.Error())
	}
	// END

	// pass through the document line by line
	for err == nil || err == io.EOF {
		var from uint
		if haveChild {
			toChild <- q
			statusChild := <-fromChild
			// child may have set q.err
			err = q.Err
			if err != nil || (statusChild|LAST_LINE_PROCESSED != 0) {
				haveChild = false
				if err == nil || err == io.EOF {
					fmt.Println("*** APPENDING BLOCKQUOTE: A ***")
					h.blocks = append(h.blocks, child)
				}
				// child = nil
			}
			goto GET_NEXT
		}
		if len(q.runes) > 0 {
			if h.depth > 0 {
				from = SkipChevrons(q, h.depth)
			}
			if q.runes[from] == '>' {
				toChild = make(chan *Line)
				fromChild = make(chan int)
				stopChild = make(chan bool)
				child, _ = NewBlockquote(h.depth + 1)
				fmt.Printf("*** CREATED BLOCKQUOTE, DEPTH %d ***\n",
					h.depth+1)
				go child.ParseHolder(p, toChild, fromChild, stopChild)
				haveChild = true
				statusChild := <-fromChild // setup complete

				toChild <- q
				statusChild = <-fromChild
				// child may have set q.err
				err = q.Err
				if err != nil || (statusChild|LAST_LINE_PROCESSED != 0) {
					haveChild = false
					if err == nil || err == io.EOF {
						fmt.Println("*** APPENDING BLOCKQUOTE: B ***")
						h.blocks = append(h.blocks, child)
					}
					child = nil
				}
				goto GET_NEXT
			}
			// HANDLE BLOCKS ----------------------------------------

			if err == nil || err == io.EOF {
				var b BlockI
				ch0 = q.runes[from]
				eol := uint(len(q.runes))

				// HEADERS --------------------------------
				if ch0 == '#' {
					b, err = q.parseHeader(from + 1)
				}

				// HORIZONTAL RULES ----------------------
				if b == nil && (err == nil || err == io.EOF) &&
					(ch0 == '-' || ch0 == '*' || ch0 == '_') {
					b, err = q.parseHRule(from)
				}

				// XXX STUB : TRY OTHER PARSERS

				// BLOCKQUOTE -----------------------------
				//if b == nil && (err == nil || err == io.EOF) && ch0 == '>' {
				//	b, err = q.parseBlockquote(doc, 1)
				//}
				// ORDERED LISTS --------------------------

				// XXX We require a space after these starting characters
				if b == nil && (err == nil || err == io.EOF) {
					myFrom := from
					for ; myFrom < from+3 && myFrom < eol; myFrom++ {
						if !u.IsSpace(q.runes[myFrom]) {
							break
						}
					}
					if myFrom < eol-2 {

						// we are positioned on a non-space character
						ch0 := q.runes[myFrom]
						ch1 := q.runes[myFrom+1]
						ch2 := q.runes[myFrom+2]
						if u.IsDigit(ch0) && ch1 == '.' && ch2 == ' ' {
							b, err = q.parseOrdered(myFrom + 2)

						}
					}
				}

				// UNORDERED LISTS ------------------------

				// XXX We require a space after these starting characters
				if b == nil && (err == nil || err == io.EOF) {
					myFrom := from
					for myFrom = 0; myFrom < 3 && myFrom < eol; myFrom++ {
						if !u.IsSpace(q.runes[myFrom]) {
							break
						}
					}
					if myFrom < eol-1 {
						// we are positioned on a non-space character
						ch0 := q.runes[myFrom]
						ch1 := q.runes[myFrom+1]
						if (ch0 == '*' || ch0 == '+' || ch0 == '-') &&
							ch1 == ' ' {

							b, err = q.parseUnordered(myFrom + 2)
						}
					}
				}

				// DEFAULT: PARA --------------------------
				// If we have parsed the line, we hang the block off
				// the document.  Otherwise, we treat whatever we have
				// as a sequence of spans and make a Para out of it.
				if err == nil || err == io.EOF {
					if b != nil {
						h.AddBlock(b)
						lastBlockLineSep = false
					} else {
						// default parser
						// DEBUG
						fmt.Printf("== invoking parseSpanSeq(true) ==\n")
						// END
						var seq *SpanSeq
						seq, err = q.parseSpanSeq(doc, from, true)
						if err == nil || err == io.EOF {
							if curPara == nil {
								curPara = new(Para)
							}
							fmt.Printf("* adding seq to curPara\n") // DEBUG
							curPara.seqs = append(curPara.seqs, *seq)
							fmt.Printf("  curPara depth %d  has %d seqs\n",
								h.depth, len(curPara.seqs))
						}
					}
				}
			}

		} else {
			// we got a blank line
			ls, err := NewLineSep(q.lineSep)
			if err == nil {
				if curPara != nil {
					h.AddBlock(curPara)
					curPara = nil
					lastBlockLineSep = false
				}
				fmt.Printf("adding LineSep to holder\n") // DEBUG
				if !lastBlockLineSep {
					h.AddBlock(ls)
					lastBlockLineSep = true
				}
			}
		}

		// prepare for next iteration ---------------------
	GET_NEXT:
		if err != nil || eofSeen {
			// DEBUG
			fmt.Printf("parseHolder depth %d breaking, error or EOF seen\n",
				h.depth)
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
		if err == io.EOF {
			eofSeen = true
			// DEBUG
			fmt.Println("*** EOF SEEN ***")
			// END
		}
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
	} // END FOR LOOP -----------------------------------------------

	if err == nil || err == io.EOF {
		if haveChild {
			fmt.Println("*** APPENDING BLOCKQUOTE: C ***")
			h.blocks = append(h.blocks, child)
		}
		if curPara != nil {
			fmt.Printf("depth %d: have dangling curPara\n", h.depth) // DEBUG
			h.AddBlock(curPara)
			curPara = nil
		}
		// DEBUG
		fmt.Printf("parseHolder depth %d returning; holder has %d blocks\n",
			h.depth, len(h.blocks))
		// END
	}

SAYOONARA:
	resp <- DONE | LAST_LINE_PROCESSED

JUST_DIE:
	return
}
