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
		"ampersand-uri",
		"asterisk",
		"asterisk-near-text",
		"backslash-escape",
		"blockquote-added-markup",
		"blockquote-line-2-paragraphs",
		"blockquote-line",
		"blockquote-multiline-1-space-begin",
		//"blockquote-multiline-1-space-end",	// \u00a0 replaces space
		"blockquote-multiline-2-paragraphs",
		"blockquote-multiline",
		"blockquote-nested",
		//"blockquote-nested-multiplereturn-level1",  // infinite recursion
		//"blockquote-nested-multiplereturn",	// nested blank line 
		//"blockquote-nested-return-level1",	// infinite recursion 
		//"code-1-tab",
		//"code-4-spaces-escaping",
		//"code-4-spaces",
		"em-middle-word",
		"em-star",
		"em-underscore",
		"entities-text-flow",
		"EOL-CR+LF",
		"EOL-CR",
		"EOL-LF",
		// "header-level1-equal-underlined",
		"header-level1-hash-sign-closed",
		"header-level1-hash-sign",
		"header-level1-hash-sign-trailing-1-space",
		//"header-level1-hash-sign-trailing-2-spaces", // shd force blank line
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
		"horizontal-rule-3-dashes",
		"horizontal-rule-3-dashes-spaces",
		"horizontal-rule-3-stars",
		"horizontal-rule-3-underscores",
		"horizontal-rule-7-dashes",
		//"img-idref",
		//"img-idref-title",
		"img",
		"img-title",
		"inline-code-escaping-entities",
		// "inline-code",				// tests doubled backtick
		//"inline-code-with-visible-backtick",
		// "line-break-2-spaces",		// spaces at end of line => <br />
		// "line-break-5-spaces",		// spaces at end of line => <br />
		// "link-automatic-email",
		// "link-automatic",
		"link-bracket-paranthesis",
		"link-bracket-paranthesis-title",
		"link-idref-angle-bracket",
		"link-idref-implicit",
		"link-idref-implicit-spaces",
		"link-idref",
		"link-idref-space",
		"link-idref-title",
		//"link-idref-title-next-line",
		"link-idref-title-paranthesis",
		"link-idref-title-single-quote",
		//"list-blockquote",
		//"list-code",
		//"list-multiparagraphs",
		//"list-multiparagraphs-tab",
		"ordered-list-escaped",
		"ordered-list-items",
		"ordered-list-items-random-number",
		"paragraph-hard-return",
		"paragraph-line",
		"paragraphs-2-leading-spaces",
		"paragraphs-3-leading-spaces",
		"paragraphs-leading-space",
		"paragraphs-trailing-spaces",
		"paragraph-trailing-leading-spaces",
		// "paragraph-trailing-tab", // trailing tab becomes spaces?
		"strong-middle-word",
		"strong-star",
		"strong-underscore",
		"unordered-list-items-dashsign",
		//"unordered-list-items-leading-1space",	// ?? bad test ??
		"unordered-list-items-leading-2spaces",
		"unordered-list-items-leading-3spaces",
		"unordered-list-items",
		"unordered-list-items-plussign",
		//"unordered-list-paragraphs",
		"unordered-list-unindented-content",
		//"unordered-list-with-indented-content",
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

	options := NewOptions(in, mdPath, outPath, true, true)
	p, err := NewParser(options)
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)
	c.Assert(p.opt, NotNil)

	doc, err := p.Parse()
	c.Assert(err, Equals, io.EOF)
	c.Assert(doc, NotNil)

	html := doc.Get()
	fmt.Printf("HTML: '%s'\n", string(html))

	// convert []MarkdownI to bytes REDUNDANT CODE ?
	var b bytes.Buffer
	var wPtr io.Writer = &b
	c.Assert(wPtr, NotNil)
	wr, err := NewHtmlWriter(wPtr)
	c.Assert(err, IsNil)
	c.Assert(wr, NotNil)
	count, err := wr.Write(html)
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
