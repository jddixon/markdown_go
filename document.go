package md

// xgo/md/document.go

import (
	"fmt"
)

var _ = fmt.Print

type Document struct {
	imgDict  map[string]*Definition
	linkDict map[string]*Definition
	Holder
}

func NewDocument() (q *Document, err error) {

	h, _ := NewHolder(false, uint(0)) // not Blockquote, depth 0
	q = &Document{
		imgDict:  make(map[string]*Definition),
		linkDict: make(map[string]*Definition),
		Holder:   *h,
	}
	return
}

// A pointer to the definition is returned to signal success.
func (q *Document) AddImgDefinition(id string, uri, title []rune) (
	def *Definition, err error) {

	if id == "" {
		err = NilDocument
	} else if len(uri) == 0 {
		err = EmptyURI
	} else {
		// XXX Note that this allows duplicate definitions
		def = &Definition{uri: uri, title: title}
		q.imgDict[id] = def
	}
	return
}

// A pointer to the definition is returned to signal success.
func (q *Document) AddLinkDefinition(id string, uri, title []rune) (
	def *Definition, err error) {

	if id == "" {
		err = NilDocument
	} else if len(uri) == 0 {
		err = EmptyURI
	} else {
		// XXX Note that this allows duplicate definitions
		def = &Definition{uri: uri, title: title}
		q.linkDict[id] = def
	}
	return
}

func (q *Document) Get() (body []rune) {
	var (
		inUnordered bool
		inOrdered   bool
	)
	for i := 0; i < len(q.blocks)-1; i++ {
		fmt.Printf("BLOCK %d\n", i)

		var ()

		block := q.blocks[i]
		content := block.Get()

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
		fmt.Printf("skipping final LineSep\n") // DEBUG
	default:
		// DEBUG
		fmt.Printf("outputting '%s'\n", string(lastBlock.Get()))
		// END
		body = append(body, lastBlock.Get()...)
	}
	if inOrdered {
		body = append(body, OL_CLOSE...)
	}
	if inUnordered {
		body = append(body, UL_CLOSE...)
	}
	// drop any terminating CR/LF
	for body[len(body)-1] == '\n' || body[len(body)-1] == '\r' {
		body = body[:len(body)-1]
	}

	// DEBUG
	fmt.Printf("Doc.Get returning '%s'\n", string(body))
	// END
	return
}
