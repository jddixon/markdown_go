package md

// xgo/md/inlineHtmlSpan_test.go

import (
	"fmt"
	//xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

var _ = fmt.Print

func (s *XLSuite) doScanOK(c *C, text string, from uint) (
	elm *InlineHtmlElm, err error) {

	runes := []rune(text)
	offset, tagNdx, empty, nestable, err := scanForTag(runes, from)
	c.Assert(err, IsNil)
	elm = &InlineHtmlElm{tagNdx, empty, nestable, offset, nil}
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

	// check the tag scanner; do this for all 22 tags

	// single characer tags
	s.doCheckOneCharTag(c, IL_TAG_A)
	s.doCheckOneCharTag(c, IL_TAG_B)
	s.doCheckOneCharTag(c, IL_TAG_I)
	s.doCheckOneCharTag(c, IL_TAG_Q)
	s.doCheckOneCharTag(c, IL_TAG_S)
	s.doCheckOneCharTag(c, IL_TAG_U)

}
