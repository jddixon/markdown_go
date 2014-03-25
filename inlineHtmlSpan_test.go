package md

// xgo/md/inlineHtmlSpan_test.go

import (
	"fmt"
	// xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

var _ = fmt.Print

func (s *XLSuite) doScanOK(c *C, text string, from uint) (
	elm *InlineHtmlElm, err error) {

	runes := []rune(text)
	offset, tagNdx, err := scanForTag(runes, from)
	c.Assert(err, IsNil)
	//// DEBUG
	//if true {
	//	fmt.Printf("%-6s returns offset %d, tagNdx %d\n",
	//		text, offset, tagNdx)
	//}
	//// END
	c.Assert(offset > 0, Equals, true)
	elm = &InlineHtmlElm{
		tagNdx: tagNdx,
		end:    offset,
	}
	return
}

// TODO: add rng argument, then add random amount of leading space
func (s *XLSuite) doCheckOneCharTag(c *C, ndx int) {
	tag := INLINE_TAGS[ndx]
	text := fmt.Sprintf("<%s>  ", tag)
	elm, err := s.doScanOK(c, text, 1)
	c.Assert(err, IsNil)
	c.Assert(elm.tagNdx, Equals, ndx)
	c.Assert(elm.end, Equals, uint(3))
}

func (s *XLSuite) doCheckOtherTag(c *C, ndx int) {
	var text string
	var expected uint
	tag := INLINE_TAGS[ndx]

	//if isEmpty[ndx] {
	//	text = fmt.Sprintf("<%s />", tag)
	//	expected = tagLen[ndx] + 4
	//} else {

	text = fmt.Sprintf("<%s>", tag)
	expected = tagLen[ndx] + 2

	// DEBUG
	// fmt.Printf("CHECKING ndx %d as '%s'\n", ndx, text)
	// END
	elm, err := s.doScanOK(c, text, 1)
	c.Assert(err, IsNil)
	c.Assert(elm.tagNdx, Equals, ndx)
	c.Assert(elm.end, Equals, expected)
}

// force me to execute first ;-)
func (s *XLSuite) TestAAAInlineHtmlSpan(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_INLINE_HTML_SPAN")
	}

	// check the lower() function
	c.Assert(lower('A'), Equals, 'a')
	c.Assert(lower('Z'), Equals, 'z')
	c.Assert(lower('a'), Equals, 'a')
	c.Assert(lower('z'), Equals, 'z')

	// check mapping between tags and constants
	c.Assert(INLINE_TAGS[IL_TAG_A], Equals, "a")
	c.Assert(INLINE_TAGS[IL_TAG_WBR], Equals, "wbr")
	c.Assert(len(INLINE_TAGS), Equals, IL_TAG_WBR+1)

	for i := IL_TAG_A; i <= IL_TAG_WBR; i++ {
		strTag := INLINE_TAGS[i]
		c.Assert(tagMap[strTag], Equals, i)
	}

	// spot-check our slice of tag lenths -- first, longest, last
	c.Assert(tagLen[0], Equals, uint(len(INLINE_TAGS[0])))
	c.Assert(tagLen[IL_TAG_STRONG], Equals,
		uint(len(INLINE_TAGS[IL_TAG_STRONG])))
	c.Assert(tagLen[IL_TAG_WBR], Equals, uint(len(INLINE_TAGS[IL_TAG_WBR])))

	// check the tag scanner; do this for all 22 tags

	// single characer tags (6)
	s.doCheckOneCharTag(c, IL_TAG_A)
	s.doCheckOneCharTag(c, IL_TAG_B)
	s.doCheckOneCharTag(c, IL_TAG_I)
	s.doCheckOneCharTag(c, IL_TAG_Q)
	s.doCheckOneCharTag(c, IL_TAG_S)
	s.doCheckOneCharTag(c, IL_TAG_U)

	s.doCheckOtherTag(c, IL_TAG_ABBR)
	s.doCheckOtherTag(c, IL_TAG_BDO)
	s.doCheckOtherTag(c, IL_TAG_BR_SIMPLE)
	s.doCheckOtherTag(c, IL_TAG_BR_SHORT)
	s.doCheckOtherTag(c, IL_TAG_BR)
	s.doCheckOtherTag(c, IL_TAG_CITE)
	s.doCheckOtherTag(c, IL_TAG_CODE)
	s.doCheckOtherTag(c, IL_TAG_DEL)
	s.doCheckOtherTag(c, IL_TAG_DFN)
	s.doCheckOtherTag(c, IL_TAG_EM)
	s.doCheckOtherTag(c, IL_TAG_INS)
	s.doCheckOtherTag(c, IL_TAG_KBD)
	s.doCheckOtherTag(c, IL_TAG_SAMP)
	s.doCheckOtherTag(c, IL_TAG_SMALL)
	s.doCheckOtherTag(c, IL_TAG_SPAN)
	s.doCheckOtherTag(c, IL_TAG_STRONG)
	s.doCheckOtherTag(c, IL_TAG_SUB)
	s.doCheckOtherTag(c, IL_TAG_VAR)
	s.doCheckOtherTag(c, IL_TAG_WBR)

}
