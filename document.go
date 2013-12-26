package md

// xgo/md/document.go

import (
//"fmt"
)

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
	for i := 0; i < len(q.blocks); i++ {
		body = append(body, q.blocks[i].Get()...)
	}
	return
}
