package main

// markdown_go/cmd/commonTests/commonTests.go

import (
	"encoding/json"
	"flag"
	"fmt"
	gm "github.com/jddixon/markdown_go"
	"io"
	"io/ioutil"
	"os"
//	"path/filepath"
	"strings"
)

func Usage() {
	fmt.Printf("Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("where the options are:\n")
	flag.PrintDefaults()
}
var (
	// These need to be referenced as pointers.
	fVerbose = flag.Bool("v", false, "be talkative")
)

var (
	// must be a subdirectory
	pathToData		= "commonmark/test-dumped.json"
)

type TestPair struct {
	EndLine		float64
	Example		float64
	Html		string
	Markdown	string
	StartLine	float64
	Section		string
}

func testPair(pair *TestPair, verbose bool) {

	if pair == nil {
		fmt.Printf("testPair: nil TestPair\n")
	} else {
		exampleNbr := int(pair.Example)
		fmt.Printf("Example %3d ========================================\n", 
			exampleNbr)

		var rd io.Reader = strings.NewReader(pair.Markdown)
		opt := gm.NewOptions(rd, "", "", false, false)
		p, err := gm.NewParser(opt)
		if err == nil {	
			doc, err := p.Parse()
			if err == io.EOF {
				err = nil
			} 
			if doc == nil {
				fmt.Printf("example %d: parse failed; err = %v\n", 
					exampleNbr, err)
			}
			if (doc != nil) && (err == nil) {
				output := string(doc.GetHtml())
				if output != pair.Html {
					fmt.Printf("    markdown: %s",	pair.Markdown)
					fmt.Printf("    expected: %s",	pair.Html)
					fmt.Printf("    actual:   %s",	output)
					fmt.Println()
				}
			}
		}	
	}
}
func testCommons(verbose bool) {

	data, err := ioutil.ReadFile(pathToData)
	if err == nil {
		fmt.Printf("read %d bytes from the disk\n", len(data))
		var pairs []TestPair
		err = json.Unmarshal(data, &pairs)
		if err == nil {
			fmt.Printf("unmarshaled %d tests\n\n", len(pairs))
			for i := 0; i < len(pairs); i++ {
				testPair(&pairs[i] , verbose)
			}
		}
	}
	if err != nil {
		fmt.Printf("  ERROR: %v\n", err)
	}

}

func main() {
	var err error
	flag.Usage = Usage
	flag.Parse()
	verbose := *fVerbose
	
	_, err = os.Stat(pathToData)
	if err == nil {
		testCommons(verbose)
	} else {
		fmt.Printf("can't locate data file '%s'\n", pathToData)
	}
}
