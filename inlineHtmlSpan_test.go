package md

// xgo/md/inlineHtmlSpan_test.go

import (
	"fmt"
	//xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

var _ = fmt.Print

// force me to execute first ;-)
func (s *XLSuite) TestAAAInlineHtmlSpan(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_INLINE_HTML_SPAN")
	}

	// beware of typos
	c.Assert(INLINE_TAGS[IL_TAG_A], Equals, "a")
	c.Assert(INLINE_TAGS[IL_TAG_WBR], Equals, "wbr")
	c.Assert(len(INLINE_TAGS), Equals, IL_TAG_WBR+1)

	c.Assert(lower('A'), Equals, 'a')
	c.Assert(lower('Z'), Equals, 'z')
	c.Assert(lower('a'), Equals, 'a')
	c.Assert(lower('z'), Equals, 'z')

	for i := IL_TAG_A; i <= IL_TAG_WBR; i++ {
		strTag := INLINE_TAGS[i]
		c.Assert(tagMap[strTag], Equals, i)
	}
}
