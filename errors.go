package md

import (
	e "errors"
)

var (
	EmptyHeaderTitle  = e.New("empty header title")
	EmptyID           = e.New("empty id parameter")
	HeaderNOutOfRange = e.New("header N out of range")
	InvalidCharInID   = e.New("invalid char in ID")
	NilID             = e.New("nil id parameter")
	NilWriter         = e.New("nil writer parameter")
	NotALineSeparator = e.New("not a line separator")
)
