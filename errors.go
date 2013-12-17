package md

import (
	e "errors"
)

var (
	EmptyHeaderTitle  = e.New("empty header title")
	HeaderNOutOfRange = e.New("header N out of range")
	NotALineSeparator = e.New("not a line separator")
)
