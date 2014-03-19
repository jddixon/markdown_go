package md

// xgo/md/inlineHtmlSpan.go

// So far, just a list of names.  Of these, only <q> may be nested in
// an element of its own type.

// This is not an acceptable Go const
var (
	INLINE_TAGS = [...]string{
		"a",
		"abbr",
		"b",
		"bdo",
		"br", // need not be closed
		"cite",
		"code",
		"del",
		"dfn",
		"em",
		"i",
		"ins",
		"kbd",
		"q",
		"s",
		"samp",
		"small",
		"span",
		"strong",
		"sub",
		"u",
		"var",
		"wbr",
	}
)

const (
	IL_TAG_A = iota
	IL_TAG_ABBR
	IL_TAG_B
	IL_TAG_BDO
	IL_TAG_BR
	IL_TAG_CITE
	IL_TAG_CODE
	IL_TAG_DEL
	IL_TAG_DFN
	IL_TAG_EM
	IL_TAG_I
	IL_TAG_INS
	IL_TAG_KBD
	IL_TAG_Q
	IL_TAG_S
	IL_TAG_SAMP
	IL_TAG_SMALL
	IL_TAG_SPAN
	IL_TAG_STRONG
	IL_TAG_SUB
	IL_TAG_U
	IL_TAG_VAR
	IL_TAG_WBR
)

var tagMap map[string]int

func init() {
	tagMap = make(map[string]int)
	tagMap["a"] = IL_TAG_A
	tagMap["abbr"] = IL_TAG_ABBR
	tagMap["b"] = IL_TAG_B
	tagMap["bdo"] = IL_TAG_BDO
	tagMap["br"] = IL_TAG_BR
	tagMap["cite"] = IL_TAG_CITE
	tagMap["code"] = IL_TAG_CODE
	tagMap["del"] = IL_TAG_DEL
	tagMap["dfn"] = IL_TAG_DFN
	tagMap["em"] = IL_TAG_EM
	tagMap["i"] = IL_TAG_I
	tagMap["ins"] = IL_TAG_INS
	tagMap["kbd"] = IL_TAG_KBD
	tagMap["q"] = IL_TAG_Q
	tagMap["s"] = IL_TAG_S
	tagMap["samp"] = IL_TAG_SAMP
	tagMap["small"] = IL_TAG_SMALL
	tagMap["span"] = IL_TAG_SPAN
	tagMap["strong"] = IL_TAG_STRONG
	tagMap["sub"] = IL_TAG_SUB
	tagMap["u"] = IL_TAG_U
	tagMap["var"] = IL_TAG_VAR
	tagMap["wbr"] = IL_TAG_WBR
}

type InlineHtmlElm struct {
	tagNdx   int
	empty    bool // never has any enclosed text, like <br/>
	nestable bool // can be nested in an element of its own type; <q>
	end      uint // offset of first char beyond start tag or element
	body     *SpanSeq
}

func lower(char rune) (ch rune) {
	ch = char
	if 'A' <= char && char <= 'Z' {
		ch += 0x20
	}
	return
}

// Enter with 'from' the offset into a slice of runes 'buf'.  We assume
// that < has been seen and from is sitting on the first character of
// a candidate tag.  If a well-formed tag is found, return its index
// and 'offset' just beyond the closing >.  If offset is zero, no
// inline HTML tag was found.  Otherwise, also return the nestable
// and empty attributes of the element.  XXX It makes more sense to
// do that through table lookup.
//
func scanForTag(buf []rune, from uint) (
	offset uint, // one beyond the closing > or 0 if not found
	tagNdx int, // the tag found
	empty, nestable bool) {

	bufLen := uint(len(buf))
	if from >= bufLen-1 {
		// no room for closing '>'
		return
	}
	var maybe bool
	ch0 := lower(buf[from])
	ch1 := lower(buf[from+1])
	switch ch0 {
	// these can either stand alone or start other tags
	case 'a':
		fallthrough
	case 'b':
		fallthrough
	case 'i':
		fallthrough
	case 's':
		fallthrough
	case 'u':
		if ch1 == '>' {
			offset = from + 2
			tagNdx = tagMap[string([]rune{ch0})]
			return
		} else {
			maybe = true
		}
	// this can only be a single-character tag
	case 'q':
		if ch1 == '>' {
			offset = from + 2
			tagNdx = IL_TAG_Q
			nestable = true
			return
		}
		return
	// these cannot stand alone but can start other tags
	case 'c':
		fallthrough
	case 'd':
		fallthrough
	case 'e':
		fallthrough
	case 'k':
		fallthrough
	case 'v':
		fallthrough
	case 'w':
		maybe = true
	// otherwise it can't start a tag, so we'll forget it
	default:
		return
	}

	if !maybe {
		return
	}

	// XXX STUB XXX

	return
}