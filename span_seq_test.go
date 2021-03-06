package md

// xgo/md/span_seq_test.go

import (
	"fmt"
	. "gopkg.in/check.v1"
)

// Test various kinds of emphasis spans with intermixed text.
func (s *XLSuite) TestParaEmphAndCode(c *C) {

	if VERBOSITY > 0 {
		fmt.Println("TEST_PARA_EMPH_AND_CODE")
	}
	opt := NewOptions(nil, "", "", false, false)
	doc, _ := NewDocument(opt) // just a dummy
	NULL_EOL := []rune{0}

	input := []rune("abc _def_ **ghi** __jkl mno__ qrs ")
	input = append(input, []rune("`kode &a <b >c` foo")...)
	q := NewLine(input, NULL_EOL)
	c.Assert(q.Err, IsNil)
	c.Assert(q, NotNil)

	eol := uint(len(input))
	seq, err := q.parseSpanSeq(opt, doc, 0, true)
	c.Assert(err, IsNil)
	c.Assert(seq, NotNil)
	c.Assert(q.offset, Equals, eol)

	spans := seq.spans
	c.Assert(len(spans), Equals, 9)

	s0 := string(spans[0].GetHtml())
	c.Assert(s0, Equals, "abc ") // a text span

	s1 := string(spans[1].GetHtml())
	c.Assert(s1, Equals, "<em>def</em>")

	s2 := string(spans[2].GetHtml())
	c.Assert(s2, Equals, " ")

	s3 := string(spans[3].GetHtml())
	c.Assert(s3, Equals, "<strong>ghi</strong>")

	s4 := string(spans[4].GetHtml())
	c.Assert(s4, Equals, " ")

	s5 := string(spans[5].GetHtml())
	c.Assert(s5, Equals, "<strong>jkl mno</strong>")

	s6 := string(spans[6].GetHtml())
	c.Assert(s6, Equals, " qrs ")

	s7 := string(spans[7].GetHtml())
	c.Assert(s7, Equals, "<code>kode &amp;a &lt;b &gt;c</code>")

	s8 := string(spans[8].GetHtml())
	c.Assert(s8, Equals, " foo")
}

// test link span with and without title
func (s *XLSuite) TestParaLinkSpan(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_PARA_LINKE_SPAN")
	}
	opt := NewOptions(nil, "", "", false, false)
	doc, _ := NewDocument(opt) // just a dummy
	EOL := []rune{CR}

	input := []rune("abc [foo](http://example.com) ")
	input2 := []rune("def [bar](/its/somewhere \"I hope\")")
	input = append(input, input2...)
	q := NewLine(input, EOL)
	c.Assert(q.Err, IsNil)
	c.Assert(q, NotNil)

	eol := uint(len(input))
	seq, err := q.parseSpanSeq(opt, doc, 0, true)
	c.Assert(err, IsNil)
	c.Assert(seq, NotNil)
	c.Assert(q.offset, Equals, eol)

	spans := seq.spans
	c.Assert(spans, NotNil)
	c.Assert(q.offset, Equals, eol)

	s0 := string(spans[0].GetHtml())
	c.Assert(s0, Equals, "abc ")

	s1 := string(spans[1].GetHtml())
	c.Assert(s1, Equals, "<a href=\"http://example.com\">foo</a>")

	s2 := string(spans[2].GetHtml())
	c.Assert(s2, Equals, " def ")

	s3 := string(spans[3].GetHtml())
	c.Assert(s3, Equals, "<a href=\"/its/somewhere\" title=\"I hope\">bar</a>")

}

// test image span with and without title
func (s *XLSuite) TestParaImageSpan(c *C) {

	if VERBOSITY > 0 {
		fmt.Println("TEST_PARA_IMAGE_SPAN")
	}
	opt := NewOptions(nil, "", "", false, false)
	doc, _ := NewDocument(opt) // just a dummy
	EOL := []rune{CR}

	// we expect to left-trim the abc
	input := []rune("   abc ![foo](/images/example.jpg) ")
	input2 := []rune("def ![bar](/its/somewhere.png \"I hope\")")
	input = append(input, input2...)
	q := NewLine(input, EOL)
	c.Assert(q.Err, IsNil)
	c.Assert(q, NotNil)

	eol := uint(len(input))
	seq, err := q.parseSpanSeq(opt, doc, 0, true)
	c.Assert(err, IsNil)
	c.Assert(seq, NotNil)
	c.Assert(q.offset, Equals, eol)

	spans := seq.spans
	c.Assert(spans, NotNil)
	c.Assert(q.offset, Equals, eol)

	s0 := string(spans[0].GetHtml())
	c.Assert(s0, Equals, "abc ")

	s1 := string(spans[1].GetHtml())
	c.Assert(s1, Equals, "<img src=\"/images/example.jpg\" alt=\"foo\" />")

	s2 := string(spans[2].GetHtml())
	c.Assert(s2, Equals, " def ")

	s3 := string(spans[3].GetHtml())
	c.Assert(s3, Equals, "<img src=\"/its/somewhere.png\" alt=\"bar\" title=\"I hope\" />")

}
