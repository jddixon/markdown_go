package md

import (
	. "launchpad.net/gocheck"
)

// Test various kinds of emphasis spans with intermixed text.
func (s *XLSuite) TestParaEmphAndCode(c *C) {

	doc := new(Document) // just a dummy
	NULL_EOL := []rune{0}

	input := []rune("abc _def_ **ghi** __jkl mno__ qrs ")
	input = append(input, []rune("`kode &a <b >c` foo")...)
	q, err := NewLine(doc, input, NULL_EOL)
	c.Assert(err, IsNil)
	c.Assert(q, NotNil)

	eol := len(input)
	b, err := q.parsePara()
	c.Assert(err, IsNil)
	c.Assert(b, NotNil)
	c.Assert(q.offset, Equals, eol)
	pa := b.(*Para)

	spans := pa.spans
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

// test link span with and without title
func (s *XLSuite) TestParaLinkSpan(c *C) {
	doc := new(Document) // just a dummy
	EOL := []rune{CR}

	input := []rune("abc [foo](http://example.com) ")
	input2 := []rune("def [bar](/its/somewhere \"I hope\")")
	input = append(input, input2...)
	q, err := NewLine(doc, input, EOL)
	c.Assert(err, IsNil)
	c.Assert(q, NotNil)

	eol := len(input)
	b, err := q.parsePara()
	c.Assert(err, IsNil)
	c.Assert(b, NotNil)
	c.Assert(q.offset, Equals, eol)
	pa := b.(*Para)

	spans := pa.spans
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

// test image span with and without title
func (s *XLSuite) TestParaImageSpan(c *C) {
	doc := new(Document) // just a dummy
	EOL := []rune{CR}

	input := []rune("abc ![foo](/images/example.jpg) ")
	input2 := []rune("def ![bar](/its/somewhere.png \"I hope\")")
	input = append(input, input2...)
	q, err := NewLine(doc, input, EOL)
	c.Assert(err, IsNil)
	c.Assert(q, NotNil)

	eol := len(input)
	b, err := q.parsePara()
	c.Assert(err, IsNil)
	c.Assert(b, NotNil)
	c.Assert(q.offset, Equals, eol)
	pa := b.(*Para)

	spans := pa.spans
	c.Assert(spans, NotNil)
	c.Assert(q.offset, Equals, eol)

	s0 := string(spans[0].Get())
	c.Assert(s0, Equals, "abc ")

	s1 := string(spans[1].Get())
	c.Assert(s1, Equals, "<a href=\"/images/example.jpg\">foo</a>")

	s2 := string(spans[2].Get())
	c.Assert(s2, Equals, " def ")

	s3 := string(spans[3].Get())
	c.Assert(s3, Equals, "<a href=\"/its/somewhere.png\" title=\"I hope\">bar</a>")

}
