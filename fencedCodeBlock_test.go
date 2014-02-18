package md

// xgo/md/fencedCodeBlock_test.go

import (
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	"io"
	. "launchpad.net/gocheck"
	"strings"
)

var _ = fmt.Print

func (s *XLSuite) doTestFencedBlock(c *C, char rune, rng *xr.PRNG) {
	prefix := []rune(rng.NextFileName(16))
	lineCount := 3 + rng.Intn(3)
	lines := make([][]rune, lineCount)
	postfix := []rune(rng.NextFileName(16))

	for i := 0; i < lineCount; i++ {
		var line []rune
		wordCount := 3 + rng.Intn(3)
		for j := 0; j < wordCount; j++ {
			word := []rune(rng.NextFileName(8))
			line = append(line, word...)
			if j < wordCount-1 {
				line = append(line, ' ')
			}
		}
		lines[i] = line
	}
	postCount := 3 + rng.Intn(3)
	var fencePosts []rune
	for i := 0; i < postCount; i++ {
		fencePosts = append(fencePosts, char)
	}
	fencePosts = append(fencePosts, '\n')

	// Construct sample input ---------------------------------------
	var text []rune
	text = append(text, prefix...)
	text = append(text, '\n')
	text = append(text, fencePosts...)
	for i := 0; i < lineCount; i++ {
		text = append(text, lines[i]...)
		text = append(text, '\n')
	}
	text = append(text, fencePosts...)
	text = append(text, postfix...)
	text = append(text, '\n')

	input := string(text)
	var rd io.Reader = strings.NewReader(input)
	options := NewOptions(rd, "", "", false, false)
	p, err := NewParser(options)
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// DEBUG
	fmt.Printf("SAMPLE INPUT:\n%s", input)
	// END

	doc, err := p.Parse()
	c.Assert(err, Equals, io.EOF)
	c.Assert(doc, NotNil)

	//line := p.readLine()
	//err = line.Err
	//c.Assert(err, Equals, nil)
	//c.Assert(line, NotNil)

	//b, err := line.parseHRule(uint(0))
	//c.Assert(err, IsNil)
	//c.Assert(b, NotNil)
	//h := b.(*HRule)

	//// test serialization -----------------------------
	//ser := string(h.GetHtml()) // GEEP

}
func (s *XLSuite) TestFencedBlock(c *C) {
	rng := xr.MakeSimpleRNG()

	s.doTestFencedBlock(c, '~', rng)
	s.doTestFencedBlock(c, '`', rng)
}
