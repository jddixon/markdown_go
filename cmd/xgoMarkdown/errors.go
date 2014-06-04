package main

import e "errors"

var (
	NothingToDo         = e.New("nothing to do")
	SrcFileDoesNotExist = e.New("source file does not exist")
)
