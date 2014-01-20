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
	opt          *Options
	isBlockquote bool
	depth        uint
	blocks       []BlockI
}

func NewHolder(opt *Options, isBq bool, depth uint) (h *Holder, err error) {

	if opt == nil {
		err = NilOptions
	} else if depth > 0 && !isBq {
		err = OnlyBlockquoteSupported
	} else {
		h = &Holder{
			opt:          opt,
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
// for this depth.  At depth N, we skip N.  If there is a space
// beyond the chevron, skip that too.
func SkipChevrons(q *Line, depth uint) (from uint) {

	var count uint
	eol := uint(len(q.runes))
	for offset := uint(0); offset < eol; offset++ {
		if q.runes[offset] == '>' {
			count++
			if count >= depth {
				from = offset + 1
				if from < eol && u.IsSpace(q.runes[from]) {
					from++
				}
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
		fatalError       bool
		iAmDone          bool
		curPara          *Para
		q                *Line
		ch0              rune
		lastBlockLineSep bool
		stopped          bool
		testing          = p.opt.Testing
		verbose          = p.opt.Verbose

		// used to control child holder (for Blockquote)
		haveChild bool
		child     *Blockquote
		toChild   chan *Line
		fromChild chan int
		stopChild chan bool
	)
	// DEBUG
	_, _ = testing, verbose
	// END
	resp <- OK // OK, setup complete
	// -- ok --------------------------------------------------------

	q = <-in // WAS q = p.readLine()
	err = q.Err
	if err == io.EOF {
		eofSeen = true
	}

	// DEBUG
	if p.testing {
		fmt.Printf("ParseHolder depth %d: first line is '%s'\n",
			h.depth, string(q.runes))
		if err == nil {
			fmt.Println("    nil error")
		} else {
			fmt.Printf("    error = %s\n", err.Error())
		}
	}
	// END

	sayGoodBye := true

	// pass through the document line by line
	for (err == nil || err == io.EOF) && !iAmDone {
		var (
			blankLine     bool
			lineProcessed bool
			from          uint
			statusChild   int
		)
		lineLen := uint(len(q.runes))
		if haveChild {
			if lineLen == 0 && err == io.EOF {
				stopChild <- true
				statusChild = <-fromChild
				haveChild = false
				// DEBUG
				if p.testing {
					fmt.Printf("*** DEPTH %d APPENDING BLOCKQUOTE, BLANK LINE, EOF:  ***\n",
						h.depth)
					fmt.Printf("    statusChild is 0x%x\n", statusChild)
					fmt.Printf("    APPENDED %s\n",
						string(child.Get()))
				}
				// END
				h.blocks = append(h.blocks, child)
				// child = nil

			} else {
				// just copy the line through to the child
				// DEBUG
				if testing {
					fmt.Printf("COPYING TO CHILD: %s\n", string(q.runes))
				}
				// END
				toChild <- q
				statusChild = <-fromChild
				lineProcessed = statusChild == ACK ||
					(statusChild == (DONE | LAST_LINE_PROCESSED))
				// child may have set q.err
				err = q.Err
				if err != nil || (statusChild&DONE != 0) {
					haveChild = false
					if err == nil || err == io.EOF {
						// DEBUG
						if testing {
							fmt.Printf("*** DEPTH %d APPENDING BLOCKQUOTE: AFTER '%s' ***\n",
								h.depth, string(q.runes))
							fmt.Printf("    err is %v\n", err)
							fmt.Printf("    statusChild is 0x%x\n", statusChild)
							fmt.Printf("    APPENDED %s\n",
								string(child.Get()))
						}
						// END
						h.blocks = append(h.blocks, child)

					}
					// child = nil
				} // FOO
			}
		}
		if !lineProcessed {
			if lineLen > 0 {
				if h.depth > 0 {
					from = SkipChevrons(q, h.depth)
					// DEBUG
					if testing {
						fmt.Printf("depth %d, length %d, SkipChevrons sets from to %d\n",
							h.depth, lineLen, from)
					}
					// END
					if from >= lineLen {
						blankLine = true
					}
				}
				// the first case arises when > is last character on line
				// XXX QUESTIONABLE LOGIC
				if !blankLine && q.runes[from] == '>' {
					toChild = make(chan *Line)
					fromChild = make(chan int)
					stopChild = make(chan bool)
					child, _ = NewBlockquote(h.opt, h.depth+1)
					if testing {
						fmt.Printf("*** CREATED BLOCKQUOTE, DEPTH %d ***\n",
							h.depth+1)
					}
					go child.ParseHolder(p, toChild, fromChild, stopChild)
					haveChild = true
					statusChild := <-fromChild // setup complete

					// DEBUG
					if testing {
						fmt.Printf("COPYING TO NEW CHILD: %s\n",
							string(q.runes))
					}
					// END
					toChild <- q
					statusChild = <-fromChild
					lineProcessed = statusChild == ACK ||
						(statusChild == (DONE | LAST_LINE_PROCESSED))

					// child may have set q.err
					err = q.Err
					if err != nil || (statusChild&LAST_LINE_PROCESSED != 0) {
						haveChild = false
						if err == nil || err == io.EOF {
							if testing {
								fmt.Println("*** APPENDING BLOCKQUOTE: B ***")
							}
							h.blocks = append(h.blocks, child)
						}
						child = nil
					}
				}
				if !lineProcessed {
					// HANDLE BLOCKS ----------------------------------------
					if !blankLine && (err == nil || err == io.EOF) {
						var b BlockI
						ch0 = q.runes[from]

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

						// ORDERED LISTS --------------------------

						// XXX We require a space after these starting characters
						if b == nil && (err == nil || err == io.EOF) {
							myFrom := from
							for ; myFrom < from+3 && myFrom < lineLen; myFrom++ {
								if !u.IsSpace(q.runes[myFrom]) {
									break
								}
							}
							if myFrom < lineLen-2 {

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
							for myFrom = 0; myFrom < 3 && myFrom < lineLen; myFrom++ {
								if !u.IsSpace(q.runes[myFrom]) {
									break
								}
							}
							if myFrom < lineLen-1 {
								// we are positioned on a non-space character
								ch0 := q.runes[myFrom]
								ch1 := q.runes[myFrom+1]
								if (ch0 == '*' || ch0 == '+' || ch0 == '-') &&
									ch1 == ' ' {

									b, err = q.parseUnordered(myFrom + 2)
								}
							}
						}

						// CODE -----------------------------------

						// XXX STUB

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
								var seq *SpanSeq
								seq, err = q.parseSpanSeq(p.opt,
									doc, from, true)
								if err == nil || err == io.EOF {
									if curPara == nil {
										curPara = new(Para)
									}
									if testing {
										fmt.Printf("* adding seq to curPara\n")
									}
									curPara.seqs = append(curPara.seqs, *seq)
									if testing {
										fmt.Printf("  curPara depth %d  has %d seqs\n",
											h.depth, len(curPara.seqs))
									}
								}
							}
						}
					}
				}

			} else {
				blankLine = true
			}
			if blankLine {
				// we got a blank line
				ls, err := NewLineSep(q.lineSep)
				if err == nil {
					if curPara != nil {
						h.AddBlock(curPara)
						curPara = nil
						lastBlockLineSep = false
					}
					if !lastBlockLineSep {
						h.AddBlock(ls)
						lastBlockLineSep = true
					}
				}
			}
		}
		// prepare for next iteration ---------------------
		if err != nil || eofSeen || iAmDone {
			// DEBUG
			if testing {
				fmt.Printf("parseHolder depth %d breaking, error or EOF seen\n",
					h.depth)
				if err != nil {
					fmt.Printf("    ERROR: %s\n", err.Error())
				}
				if eofSeen {
					fmt.Println("    EOF SEEN, so breaking")
				}
			}
			// END

			// XXX LAST_LINE_PROCESSED SHOULD BE CONDITIONAL
			resp <- DONE | LAST_LINE_PROCESSED
			break
		}

		resp <- ACK
		// -- in ----------------------------------------------------
		var ok bool
		select {
		case q, ok = <-in:
			if ok {
				err = q.Err
			} else {
				fatalError = true
				break
			}
		case stopped, ok = <-stop:
			if !ok {
				fatalError = true
				break
			} else {
				_ = stopped        // not yet used - and so far unnecessary
				sayGoodBye = false // XXX WRONG
				break
			}
		}
		if err == io.EOF {
			eofSeen = true
		}
		// -- ack ---------------------------------------------------
		if (err != nil && err != io.EOF) || q == nil {
			break
		}
		if len(q.runes) == 0 {
			if testing {
				fmt.Printf("ZERO-LENGTH LINE")
			}
			if len(q.lineSep) == 0 && q.lineSep[0] == rune(0) {
				break
			}
		}
		// DEBUG
		if testing {
			fmt.Printf("Parse: next line is '%s'\n", string(q.runes))
		}
		// END
	} // END FOR LOOP -----------------------------------------------

	if !fatalError {
		if err == nil || err == io.EOF {
			if haveChild {
				if testing {
					fmt.Println("*** APPENDING BLOCKQUOTE: C ***")
				}
				h.blocks = append(h.blocks, child)
			}
			if curPara != nil {
				// DEBUG
				if testing {
					fmt.Printf("depth %d: have dangling curPara\n", h.depth)
				}
				// END
				h.AddBlock(curPara)
				curPara = nil
			}
			// DEBUG
			if testing {
				fmt.Printf("parseHolder depth %d returning; holder has %d blocks\n",
					h.depth, len(h.blocks))
				for i := 0; i < len(h.blocks); i++ {
					fmt.Printf("BLOCK %d:%d: '%s'\n",
						h.depth, i, string(h.blocks[i].Get()))
				}
			}
			// END
		}
		if sayGoodBye {
			// XXX WRONG!
			resp <- DONE | LAST_LINE_PROCESSED
		}
	}

	return
}
