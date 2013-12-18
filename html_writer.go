package md

// xgo/md/html_writer.go

import (
	"fmt"
	"io"
)

var _ = fmt.Print

type HtmlWriter struct {
	wr io.Writer
}

func NewHtmlWriter(w io.Writer) (hw *HtmlWriter, err error) {
	if w == nil {
		err = NilWriter
	} else {
		hw = &HtmlWriter{
			wr: w,
		}
	}
	return
}

func (hw *HtmlWriter) Write(downers []MarkdownI) (
	bytesOut int, err error) {

	for i := 0; err == nil && i < len(downers); i++ {
		runes := downers[i].Get()
		s := string(runes)
		// DEBUG
		n := len(s)
		fmt.Printf("CHUNK %d, %d runes, %d bytes: '%s'\n",
			i, len(runes), n, s)
		// END
		var count int
		data := []byte(s)
		count, err = hw.wr.Write(data)
		if err == nil {
			bytesOut += count
		}
	}
	return
}
