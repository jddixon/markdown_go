package md

// markdown_go/commonTests_test.go

import (
	"fmt"
	"io"
	"io/ioutil"
	"encoding/json"
//	"path/filepath"
	"strings"

	. "gopkg.in/check.v1"
)

type TestPair struct {
	EndLine		float64
	Example		float64
	Html		string
	Markdown	string
	StartLine	float64
	Section		string
}

func (s *XLSuite) dumpTestPair(c *C, pair *TestPair) {
	example := int(pair.Example)

	fmt.Printf("example %3d\n", example)
	fmt.Printf("    markdown: %s",	pair.Markdown)
	fmt.Printf("    html:     %s",	pair.Html)

	var rd io.Reader = strings.NewReader(pair.Markdown)
	opt := NewOptions(rd, "", "", false, false)
	p, err := NewParser(opt)
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	doc, err := p.Parse()
	c.Assert(err, Equals, io.EOF)
	c.Assert(doc, NotNil)
	
	// XXX WORKING HERE XXX The string() is an error; maybe it should 
	// wrap doc.GetHtml()

	c.Assert(doc.GetHtml(), Equals, string(pair.Html))

	// THIS IS JUST TRASH.  Wrap a Reader around the input string (see
	// options.go).  Use that to create a Parser.

	//doc, _ := NewDocument(opt)		// just a placeholder
	//input  := []rune(p.Markdown)
	//q      := NewLine(input, NULL_EOL)

	//eol := uint(len(input))
	//seq, err := q.parseSpanSeq(opt, doc, 0, true)
	//c.Assert(err, IsNil)
	//c.Assert(seq, NotNil)
	//c.Assert(q.offset, Equals, eol)

	//spans := seq.spans
	//c.Assert(len(spans), Equals, 1)

	//actual := string(spans[0].GetHtml())
	//// DEBUG
	//fmt.Printf("DOCUMENT:  %s\n",	doc.GetHtml())
	//fmt.Printf("EXPECTING: %s",		p.Html)
	//fmt.Printf("ACTUAL:    %s",		actual)
	//// END
	//c.Check(actual, Equals, p.Html)
}
func (s *XLSuite) TestAAACommons(c *C) {

	if VERBOSITY > 0{
		fmt.Println("TEST_AAA_COMMONS")
	}

	pathToData		:= "commonmark/test-dumped.json"

	data, err := ioutil.ReadFile(pathToData)
	if err == nil {
		fmt.Printf("read %d bytes from the disk\n", len(data))
		var pairs []TestPair
		err = json.Unmarshal(data, &pairs)
		if err == nil {
			fmt.Printf("unmarshaled %d tests\n", len(pairs))
			for i := 0; i < 8; i++ {
				s.dumpTestPair(c,  &pairs[i] )
			}
		}
	}
	if err != nil {
		fmt.Printf("  ERROR: %v\n", err)
	}

}
