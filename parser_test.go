package md

// xgo/md/parser_test.go

import (
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	"io"
	. "launchpad.net/gocheck"
	"strings"
)

func (s *XLSuite) TestReadLine(c *C) {

	text1 := "this is line 1"
	text2 := "and this the second line"
	text3 := "unterminated stuff"
	text := text1 + "\n" + text2 + "\r" + text3

	var rd io.Reader = strings.NewReader(text)
	p, err := NewParser(rd)
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	line1, err := p.readLine()
	c.Assert(err, IsNil)
	c.Assert(line1, NotNil)
	c.Assert(string(line1.runes), Equals, text1)
	c.Assert(line1.offset, Equals, 0)
	c.Assert(len(line1.lineSep), Equals, 1)
	c.Assert(line1.lineSep[0], Equals, LF)

	line2, err := p.readLine()
	c.Assert(err, IsNil)
	c.Assert(line2, NotNil)
	c.Assert(string(line2.runes), Equals, text2)
	c.Assert(line2.offset, Equals, 0)
	c.Assert(len(line2.lineSep), Equals, 1)
	c.Assert(line2.lineSep[0], Equals, CR)

	// verify that we get (a) an io.EOF and (b) a terminating null byte
	line3, err := p.readLine()
	c.Assert(err, Equals, io.EOF)
	c.Assert(line3, NotNil)
	c.Assert(string(line3.runes), Equals, text3)
	c.Assert(line3.offset, Equals, 0)
	c.Assert(len(line3.lineSep), Equals, 1)
	c.Assert(line3.lineSep[0], Equals, rune(0))

}

func (s *XLSuite) TestLinkDefinition(c *C) {

	// 1: link definition with optional title
	id1 := "george"
	uri1 := "http://example.com"
	title1 := "title one"
	def1 := fmt.Sprintf("[%s]: %s \"%s\"", id1, uri1, title1)

	// 2: link definition without title
	id2 := "foo"
	uri2 := "http://bar.com"
	def2 := fmt.Sprintf("[%s]: %s", id2, uri2)

	text := def1 + "\n" + def2

	var rd io.Reader = strings.NewReader(text)
	p, err := NewParser(rd)
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	line1, err := p.readLine()
	c.Assert(err, IsNil) // gets io.EOF
	c.Assert(string(line1.runes), Equals, def1)

	defn1, err := line1.parseLinkDefinition(p)
	c.Assert(err, IsNil)
	c.Assert(defn1, NotNil)
	c.Assert(defn1.GetURI(), Equals, uri1)
	c.Assert(defn1.GetTitle(), Equals, title1)

	line2, err := p.readLine()
	c.Assert(err, Equals, io.EOF)
	c.Assert(string(line2.runes), Equals, def2)

	defn2, err := line2.parseLinkDefinition(p)
	c.Assert(err, IsNil)
	c.Assert(defn2, NotNil)
	c.Assert(defn2.GetURI(), Equals, uri2)
	c.Assert(defn2.GetTitle(), Equals, "")
}

func (s *XLSuite) TestImageDefinition(c *C) {

	// 1: image definition with optional title
	id1 := "pic1"
	uri1 := "/images/pic1.png"
	title1 := "img title one"
	def1 := fmt.Sprintf("![%s]: (%s \"%s\")", id1, uri1, title1)

	// 2: image definition without title
	id2 := "secondPic"
	uri2 := "/pictures/bar.jpg"
	def2 := fmt.Sprintf("![%s]:(%s)", id2, uri2)

	text := def1 + "\n" + def2

	var rd io.Reader = strings.NewReader(text)
	p, err := NewParser(rd)
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	line1, err := p.readLine()
	c.Assert(err, IsNil) // gets io.EOF
	c.Assert(string(line1.runes), Equals, def1)

	defn1, err := line1.parseImageDefinition(p)
	c.Assert(err, IsNil)
	c.Assert(defn1, NotNil)
	c.Assert(defn1.GetURI(), Equals, uri1)
	c.Assert(defn1.GetTitle(), Equals, title1)

	line2, err := p.readLine()
	c.Assert(err, Equals, io.EOF)
	c.Assert(string(line2.runes), Equals, def2)

	defn2, err := line2.parseImageDefinition(p)
	c.Assert(err, IsNil)
	c.Assert(defn2, NotNil)
	c.Assert(defn2.GetURI(), Equals, uri2)
	c.Assert(defn2.GetTitle(), Equals, "")
}

func (s *XLSuite) TestHeader(c *C) {
	rng := xr.MakeSimpleRNG()

	titles := make([]string, 6)
	for i := 0; i < 6; i++ {
		titles[i] = rng.NextFileName(16)
	}
	lines := make([]string, 6)
	for i := 0; i < 6; i++ {
		var hashes string
		for j := 0; j < i+1; j++ {
			hashes += "#"
		}
		lines[i] = hashes + titles[i]
		// half get trailing hashes
		if i%2 == 0 {
			lines[i] += hashes
		}
	}
	input := strings.Join(lines, "\n")

	var rd io.Reader = strings.NewReader(input)
	p, err := NewParser(rd)
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	for i := 0; i < 6; i++ {
		line, err := p.readLine()
		if i < 5 {
			c.Assert(err, IsNil)
		} else {
			c.Assert(err, Equals, io.EOF)
		}
		c.Assert(line, NotNil)
		b, err := line.parseHeader()
		c.Assert(err, IsNil)
		c.Assert(b, NotNil)
		h := b.(*Header)
		c.Assert(h.n, Equals, i+1)
		c.Assert(string(h.runes), Equals, string(titles[i]))
	}
}
