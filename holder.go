package md

import (
	"fmt"
	"io"
	"strings"
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
	curPara      *Para
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

func (h *Holder) String() string {
	var ss []string
	for i := 0; i < len(h.blocks); i++ {
		ss = append(ss, h.blocks[i].String())
	}
	return strings.Join(ss, "\n")
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
// beyond the chevron, skip that too.  The actual number of
// chevrons found is returned.
func SkipChevrons(q *Line, depth uint) (count, from uint) {

	var offset uint
	eol := uint(len(q.runes))
	for offset = uint(0); offset < eol; offset++ {
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

func (h *Holder) makeSimpleLineSep() (nl *LineSep) {
	newLine := []rune{'\n'}
	nl, _ = NewLineSep(newLine)
	return
}
func (h *Holder) dumpAnyPara(addNewLine, testing bool) {
	if h.curPara != nil {
		// DEBUG
		if testing {
			fmt.Printf("depth %d: have dangling h.curPara '%s'\n",
				h.depth, string(h.curPara.GetHtml()))
		}
		// END
		h.AddBlock(h.curPara)
		if addNewLine {
			lineSep := h.makeSimpleLineSep()
			h.AddBlock(lineSep)
		}
		h.curPara = nil
	}
}

// Return true if a fence is found.  If an optional language name is
// specified, return that as well.  The language name may be any
// alphanumeric string whose first character is not a digit.
func (q *Line) foundFence(from uint) (found bool, lang string) {

	var fenceChar rune

	eol := uint(len(q.runes))
	spanLen := eol - from

	if spanLen >= 3 {
		if q.runes[from+0] == '~' && q.runes[from+1] == '~' &&
			q.runes[from+2] == '~' {
			fenceChar = '~'
			found = true
		} else if q.runes[from+0] == '`' && q.runes[from+1] == '`' &&
			q.runes[from+2] == '`' {
			fenceChar = '`'
			found = true
		}
		// XXX This isn't right: ~~~XYZ will match
		if found {
			var offset uint
			// skip any more fenceposts
			for offset = from + 3; offset < eol; offset++ {
				char := q.runes[offset]
				if char != fenceChar {
					break
				}
			}
			// skip any spaces
			for offset < eol {
				char := q.runes[offset]
				if !u.IsSpace(char) {
					break
				}
				offset++
			}
			// XXX simplistic
			if offset < eol {
				rest := string(q.runes[offset:])
				lang = strings.TrimSpace(rest)
			}
		}
	}
	return
}
func (h *Holder) ParseHolder(p *Parser,
	in chan *Line, resp chan int, stop chan bool) {

	doc := p.GetDocument()
	var (
		codeBlock        = new(CodeBlock)
		fencedCodeBlock  *FencedCodeBlock
		lostChild        BlockI
		eofSeen          bool
		err              error
		fatalError       bool
		iAmDone          bool
		lineProcessed    bool
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
	_ = verbose // still not used

	resp <- OK // OK, setup complete
	// -- ok --------------------------------------------------------

	q = <-in // WAS q = p.readLine()
	err = q.Err
	if err == io.EOF {
		eofSeen = true
	}

	// DEBUG
	if p.opt.Testing {
		fmt.Printf("entering ParseHolder depth %d: first line is '%s'\n",
			h.depth, string(q.runes))
		if err != nil {
			fmt.Printf("    error = %s\n", err.Error())
		}
	}
	// END

	sayGoodbye := true

	// pass through the document line by line
	for (err == nil || err == io.EOF) && !iAmDone && !stopped {
		var (
			b           BlockI
			blankLine   bool
			forceNL     bool
			from        uint
			statusChild int
		)
		lineProcessed = false
		b = nil
		lineLen := uint(len(q.runes)) // XXX REDUNDANT
		eol := uint(len(q.runes))
		if lineLen == 0 {
			blankLine = true
		}
		if haveChild {
			if lineLen == 0 && err == io.EOF {
				stopChild <- true
				statusChild = <-fromChild
				haveChild = false

				// DEBUG
				if p.opt.Testing {
					fmt.Printf("*** DEPTH %d APPENDING BLOCKQUOTE, BLANK LINE, EOF:  ***\n",
						h.depth)
					fmt.Printf("    statusChild is 0x%x\n", statusChild)
					fmt.Printf("    APPENDED %s\n",
						string(child.GetHtml()))
				}
				// END

				// h.blocks = append(h.blocks, child)
				lostChild = child
				lineProcessed = true // we are at EOF
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
				lineProcessed = (statusChild == ACK) ||
					(statusChild == (DONE | LAST_LINE_PROCESSED))
				// DEBUG
				if testing {
					fmt.Printf("child status = 0x%x : ", statusChild)
					if lineProcessed {
						fmt.Println("child has processed line")
					} else {
						fmt.Println("child has NOT processed line")
					}
				}
				// END
				// child may have set q.err
				err = q.Err
				if err != nil || (statusChild&DONE != 0) {
					haveChild = false
					if err == nil || err == io.EOF {
						lostChild = child
					}
					// child = nil
				} // FOO
			}
		}
		if !lineProcessed {
			if lineLen > 0 {
				if h.depth > 0 {
					var count uint
					count, from = SkipChevrons(q, h.depth)
					if testing {
						fmt.Printf("depth %d, length %d, SkipChevrons finds %d, sets from to %d\n",
							h.depth, count, lineLen, from)
					}
					if count < h.depth {
						lineProcessed = false
						break
					}
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
					h.dumpAnyPara(true, testing)
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
					lineProcessed = (statusChild == ACK) ||
						(statusChild == (DONE | LAST_LINE_PROCESSED))

					// DEBUG
					if testing {
						fmt.Printf("new child status = 0x%x : ", statusChild)
						if lineProcessed {
							fmt.Println("child has processed line")
						} else {
							fmt.Println("child has NOT processed line")
						}
					}
					// END
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
				// INDENTED CODE BLOCK ==============================
				if !lineProcessed && (err == nil || err == io.EOF) {
					// if we are in a code block and this isn't code, dump
					// the code block
					spanLen := eol - from
					dumpCode := false
					if codeBlock.Size() > 0 { // we are in a code block
						if blankLine {
							dumpCode = true
						} else {
							ch0 = q.runes[from]
							if ch0 == '\t' {
								span := NewCodeLine(q.runes[from+1 : eol])
								codeBlock.Add(span)
								lineProcessed = true
							} else if spanLen < 4 {
								dumpCode = true
							} else if ch0 == ' ' && q.runes[from+1] == ' ' &&
								q.runes[from+2] == ' ' &&
								q.runes[from+3] == ' ' {

								span := NewCodeLine(q.runes[from+4 : eol])
								codeBlock.Add(span)
								lineProcessed = true
							} else {
								dumpCode = true
							}
						}
					} else { // we are not in a code block
						if !blankLine {
							ch0 = q.runes[from]
							if ch0 == '\t' {
								h.dumpAnyPara(true, testing)
								span := NewCodeLine(q.runes[from+1 : eol])
								codeBlock.Add(span)
								lineProcessed = true
							} else if spanLen >= 4 && ch0 == ' ' &&
								q.runes[from+1] == ' ' &&
								q.runes[from+2] == ' ' &&
								q.runes[from+3] == ' ' {

								h.dumpAnyPara(true, testing)
								span := NewCodeLine(q.runes[from+4 : eol])
								codeBlock.Add(span)
								lineProcessed = true
							}
						}
					}
					if dumpCode {
						h.AddBlock(codeBlock)
						codeBlock = new(CodeBlock)
					}
				}

				// FENCED CODE BLOCK ================================

				if !lineProcessed && (err == nil || err == io.EOF) {
					// if we are in a code block and this isn't code, dump
					// the code block
					dumpCode := false
					if fencedCodeBlock != nil {
						// we are in a fenced code block
						lineProcessed = true
						endingFence, _ := q.foundFence(from)
						if endingFence {
							dumpCode = true
						} else {
							span := NewCodeLine(q.runes[from:eol])
							fencedCodeBlock.Add(span)
						}

					} else { // we are not yet in a code block
						if !blankLine {
							startingFence, lang := q.foundFence(from)
							_ = lang
							if startingFence {
								lineProcessed = true
								fencedCodeBlock = new(FencedCodeBlock)
							}
						}
					}
					if dumpCode {
						h.AddBlock(fencedCodeBlock)
						fencedCodeBlock = nil
					}
				}
				if !lineProcessed {
					// HANDLE BLOCKS ----------------------------------------
					// Within this block, if b is not nil, we have found a
					// block and shouldn't look for another.

					if !blankLine && (err == nil || err == io.EOF) {
						ch0 = q.runes[from]

						// UNDERLINED HEADER ------------------------
						if ch0 == '=' {
							foundUnderline := true
							for i := uint(0); i < eol; i++ {
								if q.runes[i] != '=' {
									foundUnderline = false
									break
								}
							}
							if foundUnderline {
								if h.curPara != nil {
									// XXX A KLUDGE.  We crudely assume that
									// if any text at all has been collected
									// we can just use it as the title.

									title := strings.TrimSpace(h.curPara.String())
									h.curPara = nil
									b, err = NewHeader(1, []rune(title))
								}
							}
						}
						// HEADERS ----------------------------------
						if b == nil && (err == nil || err == io.EOF) && ch0 == '#' {
							b, forceNL, err = q.parseHeader(from + 1)
						}

						// HORIZONTAL RULES ------------------------
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

					}
				}

			} else {
				blankLine = true
			}
		}
		// DEFAULT: PARA --------------------------
		// If we have parsed the line, we hang the block off
		// the document.  Otherwise, we treat whatever we have
		// as a sequence of spans and make a Para out of it.
		if fencedCodeBlock == nil && (err == nil || err == io.EOF) {
			if b != nil {
				h.AddBlock(b)
				if forceNL {
					b = h.makeSimpleLineSep()
					h.AddBlock(b)
					lastBlockLineSep = true
				} else {
					lastBlockLineSep = false
				}
			} else if !blankLine && !lineProcessed { // XXX CHANGE 2014-01-20
				// default parser
				var seq *SpanSeq
				seq, err = q.parseSpanSeq(p.opt,
					doc, from, true)
				if err == nil || err == io.EOF {
					if h.curPara == nil {
						h.curPara = new(Para)
					}
					if testing {
						fmt.Printf("* adding seq to h.curPara\n")
					}
					h.curPara.seqs = append(h.curPara.seqs, *seq)
					if testing {
						fmt.Printf("  h.curPara depth %d  has %d seqs\n",
							h.depth, len(h.curPara.seqs))
					}
				}
				lineProcessed = true
			}
		} // end DEFAULT: PARA
		if fencedCodeBlock == nil && blankLine && !lineProcessed {
			// we got a blank line
			ls, err := NewLineSep(q.lineSep)
			if err == nil {
				if h.curPara != nil {
					h.AddBlock(h.curPara)
					h.curPara = nil
					lastBlockLineSep = false
				}
				if !lastBlockLineSep {
					h.AddBlock(ls)
					lastBlockLineSep = true
				}
			}
		} // FOO

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

			// THIS IS DANGEROUS
			//myStatus := DONE
			//if lineProcessed:
			//	myStatus |= LASTLINE_PROCESSED
			//resp <- myStatus

			sayGoodbye = true
			break
		}

		// -- in ----------------------------------------------------
		resp <- ACK
		var ok bool
		select {
		case q, ok = <-in:
			if ok {
				err = q.Err
			} else {
				// DEBUG
				fmt.Println("select, in: !ok so fatalError")
				// END
				fatalError = true
			}
		case stopped, ok = <-stop:
			// DEBUG
			if testing {
				fmt.Printf("HOLDER %d HAS BEEN STOPPED\n", h.depth)
			}
			// END
			if !ok {
				// DEBUG
				if testing {
					fmt.Println("select, stop: !ok so fatalError")
				}
				// END
				fatalError = true
			} else {
			}
			if !fatalError && haveChild {
				stopChild <- true
				statusChild = <-fromChild

				// DEBUG
				if p.opt.Testing {
					fmt.Printf("*** DEPTH %d STOPPED: QUEUING CHILD ***\n",
						h.depth)
					fmt.Printf("    statusChild is 0x%x\n", statusChild)
					fmt.Printf("    BLOCKQUOTE: %s\n",
						string(child.GetHtml()))
				}
				// END

				lostChild = child // that is, append it
				haveChild = false
			}
			break
		}
		if err == io.EOF {
			eofSeen = true
		}
		// BREAK-FORCING CONDITIONS -----------------------
		if stopped {
			sayGoodbye = true
			break
		}
		if (err != nil && err != io.EOF) || fatalError || q == nil {
			break
		}
		// ------------------------------------------------
		if len(q.runes) == 0 {
			if testing {
				fmt.Println("ZERO-LENGTH LINE")
			}
			if len(q.lineSep) == 0 && q.lineSep[0] == rune(0) {
				break
			}
		}
		// DEBUG
		if testing {
			fmt.Printf("ParseHolder %d, bottom for loop: next line is '%s'\n",
				h.depth, string(q.runes))
		}
		// END
	} // END FOR LOOP -----------------------------------------------

	if !fatalError {
		if err == nil || err == io.EOF {
			if codeBlock.Size() > 0 {
				h.AddBlock(codeBlock)
				codeBlock = new(CodeBlock) // pedantry
			}
			// XXX should never happen
			if haveChild {
				if testing {
					fmt.Println("*** APPENDING BLOCKQUOTE OUTSIDE LOOP ***")
				}
				h.blocks = append(h.blocks, child)
			}
			if lostChild != nil {
				// DEBUG
				if testing {
					fmt.Printf(
						"*** DEPTH %d APPENDING LOSTCHILD BLOCKQUOTE ***\n",
						h.depth)
					fmt.Printf("    err is %v\n", err)
					fmt.Printf("    APPENDED %s\n",
						string(lostChild.GetHtml()))
				}
				// END	// GEEP
				h.AddBlock(lostChild)
				lastBlockLineSep = false
			}
			h.dumpAnyPara(false, testing)
			// DEBUG
			if testing {
				fmt.Printf("parseHolder depth %d returning; holder has %d blocks\n",
					h.depth, len(h.blocks))
				for i := 0; i < len(h.blocks); i++ {
					fmt.Printf("    BLOCK %d:%d:\n'%s'\n",
						h.depth, i, string(h.blocks[i].GetHtml()))
				}
			}
			// END
		}
		if sayGoodbye {
			if testing {
				fmt.Printf("saying goodbye, depth %d ... \n", h.depth)
			}
			myStatus := DONE
			if lineProcessed {
				myStatus = myStatus | LAST_LINE_PROCESSED
			}
			resp <- myStatus
			if testing {
				fmt.Printf("    goodbye said, depth %d\n", h.depth)
			}
		}
	}

	return
}
