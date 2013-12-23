package md

// xgo/md/document.go

import ()

type Document struct {
	blocks []BlockI
	dict   map[string]*Definition
}

func NewDocument() (q *Document, err error) {

	q = &Document{}
	return
}
