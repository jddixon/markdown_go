package md

// xgo/md/options.go

import (
	"io"
)

type Options struct {
	Reader  io.Reader
	InFile  string
	OutFile string
	Testing bool
	Verbose bool
}

func NewOptions(rd io.Reader, inFile, outFile string, testing, verbose bool) (
	o *Options) {

	o = &Options{
		Reader:  rd,
		InFile:  inFile,
		OutFile: outFile,
		Testing: testing,
		Verbose: verbose,
	}
	return o
}
