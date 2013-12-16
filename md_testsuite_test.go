package md

// xgo/lex/lex_input_test.go

import (
	"bytes"
	"fmt"
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

		"EOL-CR+LF",
		"EOL-CR",
		"EOL-LF",

		"paragraph-hard-return",
		"paragraph-line",
		"paragraphs-2-leading-spaces", // drop leading spaces in para
		"paragraphs-3-leading-spaces",
		"paragraphs-leading-space",
		"paragraphs-trailing-spaces",
		"paragraph-trailing-leading-spaces",
		"paragraph-trailing-tab", // trailing tab becomes spaces?
	}
)

func (s *XLSuite) doMDTest(c *C, name string) {
	fmt.Printf("TEST %s\n", name)
	path := path.Join(PATH_TO_TESTS, name)
	mdPath := path + ".md"
	outPath := path + ".out"
	fmt.Printf("%s => %s\n", mdPath, outPath)

	// convert []rune to []MarkdownI
	in, err := os.Open(mdPath)
	c.Assert(err, IsNil)
	c.Assert(in, NotNil)

	bits, err := Parse(in)
	c.Assert(err, IsNil)
	c.Assert(len(bits) > 0, Equals, true)

	// convert []MarkdownI to HTML == []rune
	output, err := HtmlWrite(bits)
	c.Assert(err, IsNil)
	// outBytes := []byte(output)

	// verify equality
	expectedOut, err := ioutil.ReadFile(outPath)
	c.Assert(err, IsNil)
	expectedRunes := bytes.Runes(expectedOut)

	// DEBUG
	fmt.Printf("EXPECTED: %s\nACTUAL: %s\n",
		string(expectedRunes), string(output))
	// END
	c.Assert(len(output), Equals, len(expectedRunes))
	for i := 0; i < len(output); i++ {
		c.Assert(output[i], Equals, expectedRunes[i])
	}
}

func (s *XLSuite) TestTestsInSuite(c *C) {

	for i := 0; i < len(TEST_NAMES); i++ {
		name := TEST_NAMES[i]
		s.doMDTest(c, name)
	}
}
