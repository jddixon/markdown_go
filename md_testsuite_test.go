package md

// xgo/lex/md_testsuite_test.go

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"path"
)

var _ = fmt.Print

const (
	PATH_TO_TESTS = "markdown-testsuite/tests"
)

// A subset of the markdown test suite.  Where test names are commented
// out, they currently fail.  The converse is not necessarily true.

var (
	TEST_NAMES = []string{
		"2-paragraphs-hard-return",
		"2-paragraphs-hard-return-spaces",
		"2-paragraphs-line",
		"2-paragraphs-line-returns",
		"2-paragraphs-line-spaces",
		"2-paragraphs-line-tab",
		"ampersand-text-flow",

		// This tests the interpretation of [ampersand] reference
		// "ampersand-uri",

		"asterisk",
		"asterisk-near-text",
		"backslash-escape",

		"em-middle-word",
		"em-star",
		"em-underscore",
		"entities-text-flow",
		// "EOL-CR+LF",
		// "EOL-CR",
		// "EOL-LF",

		"header-level1-hash-sign-closed",
		"header-level1-hash-sign",
		"header-level1-hash-sign-trailing-1-space",
		// "header-level1-hash-sign-trailing-2-spaces",

		"header-level2-hash-sign-closed",
		"header-level2-hash-sign",
		"header-level3-hash-sign-closed",
		"header-level3-hash-sign",
		"header-level4-hash-sign-closed",
		"header-level4-hash-sign",
		"header-level5-hash-sign-closed",
		"header-level5-hash-sign",
		"header-level6-hash-sign-closed",
		"header-level6-hash-sign",

		// "link-idref",				, BUT NOT YET

		"paragraph-hard-return",
		"paragraph-line",
		"paragraphs-2-leading-spaces", // drop leading spaces in para
		"paragraphs-3-leading-spaces",
		"paragraphs-leading-space",
		"paragraphs-trailing-spaces",
		"paragraph-trailing-leading-spaces",
		// "paragraph-trailing-tab", // trailing tab becomes spaces?
		"strong-middle-word",
		"strong-star",
		"strong-underscore",
	}
)

func (s *XLSuite) doMDTest(c *C, name string) {
	fmt.Printf("TEST %s\n", name)
	path := path.Join(PATH_TO_TESTS, name)
	mdPath := path + ".md"
	outPath := path + ".out"

	// convert []rune to []MarkdownI
	in, err := os.Open(mdPath)
	c.Assert(err, IsNil)
	c.Assert(in, NotNil)

	p, err := NewParser(in)
	c.Assert(err, IsNil)

	doc, err := p.Parse()
	c.Assert(err, Equals, io.EOF)
	c.Assert(doc, NotNil)

	// DEBUG ?
	html := doc.Get()
	fmt.Printf("HTML: '%s'\n", string(html))
	// END

	// convert []MarkdownI to bytes
	var b bytes.Buffer
	var wPtr io.Writer = &b
	c.Assert(wPtr, NotNil)
	wr, err := NewHtmlWriter(wPtr)
	c.Assert(err, IsNil)
	c.Assert(wr, NotNil)
	count, err := wr.Write(doc.blocks)
	c.Assert(err, IsNil)
	_ = count
	actualOut := string(b.Bytes())

	bytesFromDisk, err := ioutil.ReadFile(outPath)
	c.Assert(err, IsNil)
	expectedOut := string(bytesFromDisk)

	// c.Assert(len(actualOut), Equals, len(expectedOut))
	c.Assert(actualOut, Equals, expectedOut)
}

func (s *XLSuite) TestTestsInSuite(c *C) {

	for i := 0; i < len(TEST_NAMES); i++ {
		name := TEST_NAMES[i]
		s.doMDTest(c, name)
	}
}
