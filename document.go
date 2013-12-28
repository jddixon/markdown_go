package md

// xgo/md/document.go

import (
	"fmt"
)

var _ = fmt.Print

type Document struct {
	// blocks []BlockI
	blocks []MarkdownI
	dict   map[string]*Definition
}

func NewDocument() (q *Document, err error) {

	q = &Document{
		dict: make(map[string]*Definition),
	}
	return
}

func (q *Document) addBlock(block MarkdownI) {
	q.blocks = append(q.blocks, block)
}

// A pointer to the definition is returned to signal success.
func (q *Document) addDefinition(id string, uri, title []rune) (
	def *Definition, err error) {

	if id == "" {
		err = NilDocument
	} else if len(uri) == 0 {
		err = EmptyURI
	} else {
		// XXX Note that this allows duplicate definitions
		def = &Definition{uri: uri, title: title}
		q.dict[id] = def
	}
	return
}

func (q *Document) Get() (body []rune) {
	// DEBUG
	fmt.Printf("Document.Get() sees %d blocks\n", len(q.blocks))
	// END
	for i := 0; i < len(q.blocks)-1; i++ {
		fmt.Printf("BLOCK %d\n", i)
		body = append(body, q.blocks[i].Get()...)
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
	// XXX HACK
	if body[len(body)-1] == '\n' || body[len(body)-1] == '\r' {
		body = body[:len(body)-1]
	}
	// END HACK

	// DEBUG
	fmt.Printf("Doc.Get returning '%s'\n", string(body))
	// END
	return
}
