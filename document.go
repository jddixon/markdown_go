package md

// xgo/md/document.go

import (
	"fmt"
)

var _ = fmt.Print

type Document struct {
	refDict map[string]*Definition
	Holder
}

func NewDocument(opt *Options) (q *Document, err error) {

	h, err := NewHolder(opt, false, uint(0)) // not Blockquote, depth 0
	if err == nil {
		q = &Document{
			refDict: make(map[string]*Definition),
			Holder:  *h,
		}
	}
	return
}

// A pointer to the definition is returned to signal success.
func (q *Document) AddDefinition(id string, uri, title []rune, isImg bool) (
	def *Definition, err error) {

	if id == "" {
		err = NilDocument
	} else if len(uri) == 0 {
		err = EmptyURI
	} else {
		// XXX Note that this allows duplicate definitions
		def = &Definition{uri: uri, title: title, isImg: isImg}
		q.refDict[id] = def
	}
	return
}

func (q *Document) GetHtml() (body []rune) {
	if len(q.blocks) > 0 {
		var (
			inUnordered bool
			inOrdered   bool
			testing     = q.opt.Testing
		)
		// DEBUG
		if testing {
			fmt.Printf("Document.GetHtml(): have %d blocks\n", len(q.blocks))
		}
		// END
		for i := 0; i < len(q.blocks)-1; i++ {

			block := q.blocks[i]
			content := block.GetHtml()

			switch block.(type) {
			case *Ordered:
				if inUnordered {
					inUnordered = false
					body = append(body, UL_CLOSE...)
				}
				if !inOrdered {
					inOrdered = true
					body = append(body, OL_OPEN...)
				}
			case *Unordered:
				if inOrdered {
					inOrdered = false
					body = append(body, OL_CLOSE...)
				}
				if !inUnordered {
					inUnordered = true
					body = append(body, UL_OPEN...)
				}
			default:
				if inUnordered {
					body = append(body, UL_CLOSE...)
					inUnordered = false
				}
				if inOrdered {
					body = append(body, OL_CLOSE...)
					inOrdered = false
				}
			}
			body = append(body, content...)

		}

		// output last block IF it is not a LineSep
		lastBlock := q.blocks[len(q.blocks)-1]
		switch lastBlock.(type) {
		case *LineSep:
			// do nothing
			if testing {
				fmt.Printf("skipping final LineSep\n") // DEBUG
			}
		default:
			// DEBUG
			if testing {
				fmt.Printf("outputting '%s'\n", string(lastBlock.GetHtml()))
			}
			// END
			body = append(body, lastBlock.GetHtml()...)
		}
		if inOrdered {
			body = append(body, OL_CLOSE...)
		}
		if inUnordered {
			body = append(body, UL_CLOSE...)
		}
		// drop any terminating CR/LF
		// 2015-04-01 dropped this action (commonmark needs terminating
		// newline
		//for body[len(body)-1] == '\n' || body[len(body)-1] == '\r' {
		//	body = body[:len(body)-1]
		//
		// 2015-04-01 for commonmark
		if (len(body) == 0) || (body[len(body)-1] != '\n') {
			body = append(body, '\n')
		}
	}
	// XXX Blockquote is preceded by LF
	if len(body) > 0 && body[0] == LF {
		body = body[1:]
	}
	return
}
