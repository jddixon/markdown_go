package md

// xgo/md/inlineHtmlSpan.go

import (
	"fmt"
	"strings"
)

var _ = fmt.Print

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
var (
	_tagCount  = len(INLINE_TAGS)
	isNestable = make([]bool, _tagCount)
	isEmpty    = make([]bool, _tagCount)
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

var tagLen = make([]uint, len(INLINE_TAGS))
var tagMap map[string]int

func init() {
	isEmpty[IL_TAG_BR] = true
	isNestable[IL_TAG_Q] = true

	for i := 0; i < len(INLINE_TAGS); i++ {
		tagLen[i] = uint(len(INLINE_TAGS[i]))
	}
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
// XXX PROBABLY SHOULD DROP err - failure to match is not an error,
// and offset == 0 means not found.
//
func scanForTag(buf []rune, from uint) (
	offset uint, // one beyond the closing > or 0 if not found
	tagNdx int, // the tag found
	err error) {

	bufLen := uint(len(buf))
	if from >= bufLen-1 {
		// no room for closing '>'
		return
	}
	var maybe bool
	ch0 := lower(buf[from])
	ch1 := lower(buf[from+1])

	if ch0 == 'q' {
		if ch1 == '>' {
			offset = from + 2
			tagNdx = IL_TAG_Q
		}
		return
	}
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

	// the shortest pattern, "em>", needs three characters to complete
	if !maybe || from+3 >= bufLen {
		return
	}
	matched := false
	ch2 := lower(buf[from+2])
	if ch0 == 'e' {
		matched = ch1 == 'm' && ch2 == '>'
		if !matched {
			return
		}
		offset = from + 3
	} else if ch0 == 'b' {
		if ch1 == 'r' {
			// accept any of <br> or <br/> or <br />"
			if ch2 == '>' {
				matched = true
				offset = from + 3
			} else if bufLen == from+4 {
				if buf[from+3] == '/' && buf[from+4] == '>' {
					matched = true
					offset = from + 5
				}
			} else {
				if buf[from+3] == ' ' && buf[from+4] == '/' && buf[from+5] == '>' {
					matched = true
					offset = from + 6
				}
			}

		} else if ch1 == 'd' {
			if from+3 < bufLen {
				matched = ch2 == 'o' && buf[from+3] == '>'
				offset = from + 4
			}
		}
		if !matched {
			return
		}
	}
	if !matched {
		// all other possible matches need at least four characters
		if from+4 >= bufLen {
			return
		}
		ch3 := lower(buf[from+3])
		// all of these have 3-character tags
		if ch0 == 'd' || ch0 == 'i' || ch0 == 'k' || ch0 == 'v' || ch0 == 'w' {
			// tags are all 3 characters
			if ch3 != '>' {
				return
			}
			switch ch0 {
			case 'd':
				matched = ch1 == 'e' && ch2 == 'l'
			case 'i':
				matched = ch1 == 'n' && ch2 == 's'
			case 'k':
				matched = ch1 == 'b' && ch2 == 'd'
			case 'v':
				matched = ch1 == 'a' && ch2 == 'r'
			case 'w':
				matched = ch1 == 'b' && ch2 == 'r'
			default:
				fmt.Printf("INTERNAL ERROR: '%c' seen in level 3 switch\n",
					ch0)
			}
			if !matched {
				return
			}
			offset = from + 4
		}
		if !matched {
			// DEBUG
			fmt.Printf("checking + 5: ch0 is %c\n", ch0)
			// END

			// all other possible matches need at least five characters
			if from+5 > bufLen {
				return
			}
			ch4 := lower(buf[from+4])
			if ch0 == 'a' || ch0 == 'c' {
				if ch4 != '>' {
					return
				}
				if ch0 == 'a' {
					matched = ch1 == 'b' && ch2 == 'b' && ch3 == 'r'
				} else {
					matched = ((ch1 == 'i' && ch2 == 't') ||
						(ch1 == 'o' && ch2 == 'd')) && ch3 == 'e'
				}
				if matched {
					offset = from + 5
				}
			} else {
				if ch0 != 's' {
					return
				}

				// samp, span, small, strong
				if ch4 == '>' {
					matched = (ch1 == 'a' && ch2 == 'm' && ch3 == 'p') ||
						(ch1 == 'p' && ch2 == 'a' && ch3 == 'n')
					if matched {
						offset = from + 5
					}
				} else if from+6 >= bufLen {
					ch5 := lower(buf[from+5])
					if ch5 == '>' {
						matched = ch1 == 'm' && ch1 == 'a' && ch3 == 'l' && ch4 == 'l'
						if matched {
							offset = from + 7
						}
					} else if from+7 >= bufLen {
						if buf[from+6] == '>' {
							matched = ch1 == 't' && ch2 == 'r' &&
								ch3 == 'o' && ch4 == 'n' && ch5 == 'g'
							if matched {
								offset = from + 8
							}
						}
					}
				}
				if !matched {
					return
				}
			}
			// if we get here, we found a match
		}
		// XXX won't work  with <br/>, <br />
		tag := buf[from : offset-1]
		strTag := strings.ToLower(string(tag))
		tagNdx = tagMap[strTag]
		// DEBUG
		fmt.Printf("MATCH %s, index %d\n", strTag, tagNdx)
		// END
	}
	return
}
