package md

import (
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestEmphAndCode(c *C) {

	doc := new(Document)	// just a dummy

	input := []rune("abc _def_ **ghi** __jkl mno__ qrs ")
	input = append(input, []rune("`kode &a <b >c` foo")...)
	q, err := NewLine(doc, input, rune(0))
	c.Assert(err, IsNil)
	c.Assert(q, NotNil)

	eol := len(input)
	spans, err := q.parseToSpans()
	c.Assert(err, IsNil)
	c.Assert(spans, NotNil)
	c.Assert(q.offset, Equals, eol)
	c.Assert(len(spans), Equals, 9)

	s0 := string(spans[0].Get())
	c.Assert(s0, Equals, "abc ") // a text span

	s1 := string(spans[1].Get())
	c.Assert(s1, Equals, "<em>def</em>")

	s2 := string(spans[2].Get())
	c.Assert(s2, Equals, " ")

	s3 := string(spans[3].Get())
	c.Assert(s3, Equals, "<strong>ghi</strong>")

	s4 := string(spans[4].Get())
	c.Assert(s4, Equals, " ")

	s5 := string(spans[5].Get())
	c.Assert(s5, Equals, "<strong>jkl mno</strong>")

	s6 := string(spans[6].Get())
	c.Assert(s6, Equals, " qrs ")

	s7 := string(spans[7].Get())
	c.Assert(s7, Equals, "<code>kode &amp;a &lt;b &gt;c</code>")

	s8 := string(spans[8].Get())
	c.Assert(s8, Equals, " foo")

}

func (s *XLSuite) TestLinkSpan(c *C) {
	doc := new(Document)	// just a dummy

	input := []rune("abc [foo](http://example.com) ")
	input2 := []rune("def [bar](/its/somewhere \"I hope\")")
	input = append(input, input2...)
	q, err := NewLine(doc, input, CR)
	c.Assert(err, IsNil)
	c.Assert(q, NotNil)

	eol := len(input)
	spans, err := q.parseToSpans()
	c.Assert(err, IsNil)
	c.Assert(spans, NotNil)
	c.Assert(q.offset, Equals, eol)

	s0 := string(spans[0].Get())
	c.Assert(s0, Equals, "abc ")

	s1 := string(spans[1].Get())
	c.Assert(s1, Equals, "<a href=\"http://example.com\">foo</a>")

	s2 := string(spans[2].Get())
	c.Assert(s2, Equals, " def ")

	s3 := string(spans[3].Get())
	c.Assert(s3, Equals, "<a href=\"/its/somewhere\" title=\"I hope\">bar</a>")

}
