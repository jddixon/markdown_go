package md

import (
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestLine(c *C) {

	input := []rune("abc _def_ **ghi** __jkl mno__ qrs")
	q, err := NewLine(input, rune(0))
	c.Assert(err, IsNil)
	c.Assert(q, NotNil)

	eol := len(input)
	spans, err := q.parseToSpans()
	c.Assert(err, IsNil)
	c.Assert(spans, NotNil)
	c.Assert(q.offset, Equals, eol)
	//	c.Assert(len(spans), Equals, 7)

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
	c.Assert(s6, Equals, " qrs")

}
