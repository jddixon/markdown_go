package md

// xgo/md/linkReference.go

import (
//"fmt"
//u "unicode"
)

// In the Markdown source, a reference looks like
//     [link text] [id]
// The space between the two pairs of brackets is optional.  If the second
// pair of brackets is empty, then the id is "implicit", meaning that it
// is the same as the link text.
//
type LinkReference struct {
	id       string
	linkText []rune
	doc      *Document
}

// Return the link referred to, first looking up id in the parser
// dictionary.  This will follow the pattern
//   <a href="URI" title="TITLE">LINK_TEXT</a>
// If the Defintion contains no title, the second attribute will be
// omitted.
//
func (ref *LinkReference) Get() (runes []rune) {

	def := ref.doc.linkDict[ref.id]
	uri := def.uri
	title := def.title

	runes = append(runes, []rune("<a href=\"")...)
	runes = append(runes, uri...)
	runes = append(runes, '"')
	if len(title) > 0 {
		runes = append(runes, []rune(" title=\"")...)
		runes = append(runes, title...)
		runes = append(runes, '"')
	}
	runes = append(runes, '>')
	runes = append(runes, ref.linkText...)
	runes = append(runes, []rune("</a>")...)

	return
}
