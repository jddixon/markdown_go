package md

// xgo/md/linkRefSpan.go

// LINK REF SPAN ----------------------------------------------------

// XXX WE AREN'T DOING THIS YET

// In Markdown serialization, a LinkRef looks like
//     [linkText][id]
type LinkRefSpan struct {
	linkText []rune
	id       string
}
