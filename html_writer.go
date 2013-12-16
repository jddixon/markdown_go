package md

// xgo/md/html_writer.go

import (
	"fmt"
)

var _ = fmt.Print

func HtmlWrite(bits []MarkdownI) (out []rune, err error) {

	for i := 0; i < len(bits); i++ {
		// DEBUG
		runes := bits[i].Get()
		s := string(runes)
		n := len(s)
		fmt.Printf("CHUNK %d, %d runes, %d bytes: '%s'\n",
			i, len(runes), n, s)
		// END
		out = append(out, bits[i].Get()...)
	}
	return
}
