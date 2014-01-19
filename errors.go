package md

import (
	e "errors"
)

var (
	ChildNdxOutOfRange      = e.New("child index out of range")
	EmptyHeaderTitle        = e.New("empty header title")
	EmptyID                 = e.New("empty id parameter")
	EmptyOrderedItem        = e.New("empty ordered item")
	EmptyUnorderedItem      = e.New("empty unordered item")
	EmptyURI                = e.New("empty uri parameter")
	HeaderNOutOfRange       = e.New("header N out of range")
	InvalidCharInID         = e.New("invalid char in ID")
	InvalidLineSeparator    = e.New("invalid line separator (not LF, CR, zero")
	NilChild                = e.New("nil child parameter")
	NilDocument             = e.New("nil document parameter")
	NilID                   = e.New("nil id parameter")
	NilOptions              = e.New("nil options parameter")
	NilParser               = e.New("nil parser parameter")
	NilWriter               = e.New("nil writer parameter")
	NotALineSeparator       = e.New("not a line separator")
	OnlyBlockquoteSupported = e.New("only blockquote supported at depth")
)
