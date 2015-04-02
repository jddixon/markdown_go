package md

// xgo/md/holder.go

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
		if testing {
			fmt.Printf("depth %d: have dangling h.curPara '%s'\n",
				h.depth, string(h.curPara.GetHtml()))
		}
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

// Parse a holder, given a pointer to the main processor and the first
// line of text.  Returns a status code and a possibly non-empty line
// pointer.  If the line is not empty, it contains the first line that
// the holder could not process.
func (h *Holder) ParseHolder(p *Parser, q *Line) (out *Line, status int) {

	// get error from the first line ----------------------
	var eofSeen bool
	err := q.Err
	if err == io.EOF {
		eofSeen = true
	}

	// default value of out -------------------------------
	emptyRunes := make([]rune, 0)
	nullRune := make([]rune, 1) // a single rune with value zero
	out = NewLine(emptyRunes, nullRune)
	out.Err = io.EOF

	// set up other local variables -----------------------
	doc := p.GetDocument()
	var (
		codeBlock        = new(CodeBlock)
		fencedCodeBlock  *FencedCodeBlock
		fatalError       bool
		lineProcessed    bool
		ch0              rune
		lastBlockLineSep bool
		testing          = p.opt.Testing
		verbose          = p.opt.Verbose

		// used to control child holder (for Blockquote)
		child     *Blockquote
		lostChild BlockI
	)
	_ = verbose // still not used

	// -- top -------------------------------------------------------

	if p.opt.Testing {
		fmt.Printf("entering ParseHolder depth %d: first line is '%s'\n",
			h.depth, string(q.runes))
		if err != nil {
			fmt.Printf("    error = %s\n", err.Error())
		}
	}

	// pass through the document line by line
	for err == nil || err == io.EOF {
		var (
			b           BlockI
			blankLine   bool
			forceNL     bool
			from        uint
			statusChild int
		)
		lineProcessed = false
		b = nil
		/////////////////////////////////////////////////////////////
		// XXX THIS IS CAUSING ERRORS -- around line 339 eol has a 
		// value differnt from len(q.runes) XXX
		/////////////////////////////////////////////////////////////
		lineLen := uint(len(q.runes)) // XXX REDUNDANT
		eol := uint(len(q.runes))     // XXX identical to lineLen)
		if lineLen == 0 {
			blankLine = true
		}

		if !lineProcessed && h.depth == 0 && !blankLine {
			// HANDLE DEFINITIONS -----------------------------------
			var (
				imageDefn *Definition
				linkDefn  *Definition
			)
			// rigidly require that definitions start in the first column
			if q.runes[0] == '[' { // possible link definition
				linkDefn, err = q.parseLinkDefinition(p.opt, doc)
			}
			if err == nil && linkDefn == nil && q.runes[0] == '!' {
				imageDefn, err = q.parseImageDefinition(p.opt, doc)
			}
			if imageDefn != nil || linkDefn != nil {
				lineProcessed = true
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
						if testing {
							fmt.Printf("  BLANK LINE\n")
						}
					}
				}
				// the first case arises when > is last character on line
				// XXX QUESTIONABLE LOGIC
				if !blankLine && q.runes[from] == '>' {
					// toChild = make(chan *Line)
					// fromChild = make(chan int)
					// stopChild = make(chan bool)
					child, _ = NewBlockquote(h.opt, h.depth+1)
					if testing {
						fmt.Printf("*** CREATED BLOCKQUOTE, DEPTH %d ***\n",
							h.depth+1)
					}
					h.dumpAnyPara(true, testing)
					q, statusChild = child.ParseHolder(p, q)

					lineProcessed = (statusChild == ACK) ||
						(statusChild == (DONE | LAST_LINE_PROCESSED))

					if testing {
						fmt.Printf("new child status = 0x%x : ", statusChild)
						if lineProcessed {
							fmt.Println("child has processed line")
						} else {
							fmt.Println("child has NOT processed line")
						}
					}
					// child may have set q.err
					err = q.Err
					if err != nil || (statusChild&LAST_LINE_PROCESSED != 0) {
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

					// XXX A HACK: recalculate eol --------
					eol = uint(len(q.runes))
					lineLen = eol
					if lineLen == 0 {
						blankLine = true
					}
					// XXX END HACK -----------------------

					spanLen := eol - from
					
					// DEBUG -- what happens if from > eol ??
					//actualEOL := uint(len(q.runes))
					//fmt.Printf("eol %d actualEOL %d lineLen %d from %d spanLen %d\n",
					//	eol, actualEOL, lineLen, from, spanLen)
					// END
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

									//runish := h.curPara.GetHtml()
									//title := strings.TrimSpace(string(runish))

									crude := h.curPara.String()
									title := strings.TrimSpace(crude)
									h.curPara = nil
									b, err = NewHeader(1, []rune(title))
								}
							}
						}
						// HEADERS ----------------------------------
						if b == nil &&
							(err == nil || err == io.EOF) && ch0 == '#' {
							b, forceNL, err = q.parseHeader(from + 1)

						}

						// HORIZONTAL RULES ------------------------
						if b == nil && (err == nil || err == io.EOF) {
							b, err = q.parseHRule(from)
						}

						// XXX STUB : TRY OTHER PARSERS

						// ORDERED LISTS --------------------------

						// XXX We require a space after these starting chars
						if b == nil && (err == nil || err == io.EOF) {
							var myFrom uint
							for myFrom = from; (myFrom < from+3) && (myFrom < lineLen); myFrom++ {

								if !u.IsSpace(q.runes[myFrom]) {
									break
								}
							}
							if (lineLen > 2) && (myFrom < lineLen-2) {
								// DEBUG
								//fmt.Printf("myFrom %d lineLen %d actual EOL %d\n",
								//	myFrom, lineLen, uint(len(q.runes)))
								// END
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

						// XXX We require a space after these starting chars
						if b == nil && (err == nil || err == io.EOF) {
							var myFrom uint
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
		}

		// prepare for next iteration ---------------------
		if err != nil || eofSeen {
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
			status = DONE
			if lineProcessed {
				status |= LAST_LINE_PROCESSED
			}
			break
		}

		// -- in ----------------------------------------------------
		q = p.readLine()
		err = q.Err

		if testing {
			fmt.Printf("IN: line = '%s'\n", string(q.runes))
			if err != nil {
				fmt.Printf("    err = %s\n", err.Error())
			}
		}

		if err == io.EOF {
			eofSeen = true
		}
		// BREAK-FORCING CONDITIONS -----------------------

		if (err != nil && err != io.EOF) || fatalError || q == nil {
			break
		}
		// ------------------------------------------------
		if len(q.runes) == 0 {
			if testing {
				fmt.Println("ZERO-LENGTH LINE")
			}
			if eofSeen {
				break
			}
			// XXX previous condition
			if len(q.lineSep) == 0 && q.lineSep[0] == rune(0) {
				break
			}
		}
		if testing {
			fmt.Printf("ParseHolder %d, bottom for loop: next line is '%s'\n",
				h.depth, string(q.runes))
		}
	} // END FOR LOOP -----------------------------------------------

	if !fatalError {
		if err == nil || err == io.EOF {
			if codeBlock.Size() > 0 {
				h.AddBlock(codeBlock)
				codeBlock = new(CodeBlock) // pedantry
			}
			h.dumpAnyPara(false, testing)
			if lostChild != nil {
				if testing {
					fmt.Printf(
						"*** DEPTH %d APPENDING LOSTCHILD BLOCKQUOTE ***\n",
						h.depth)
					fmt.Printf("    err is %v\n", err)
					fmt.Printf("    APPENDED %s\n",
						string(lostChild.GetHtml()))
				}
				h.AddBlock(lostChild)
				lastBlockLineSep = false
			}
			if testing {
				fmt.Printf("parseHolder depth %d returning; holder has %d blocks\n",
					h.depth, len(h.blocks))
				for i := 0; i < len(h.blocks); i++ {
					fmt.Printf("    BLOCK %d:%d:\n'%s'\n",
						h.depth, i, string(h.blocks[i].GetHtml()))
				}
			}
		}
		if testing {
			fmt.Printf("saying goodbye, depth %d ... \n", h.depth)
		}
		status = DONE
		if lineProcessed {
			status |= LAST_LINE_PROCESSED
		}
		if testing {
			fmt.Printf("    goodbye said, depth %d\n", h.depth)
		}
	}

	return
}
