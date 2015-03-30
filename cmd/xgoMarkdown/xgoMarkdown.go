package main

// xgo/cmd/xgoMarkdown/xgoMarkdown.go

import (
	"flag"
	"fmt"
	gm "github.com/jddixon/xgo/md"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Usage() {
	fmt.Printf("Usage: %s [OPTIONS] inDir [inDir ...\n", os.Args[0])
	fmt.Printf("where the options are:\n")
	flag.PrintDefaults()
}

const ()

// The main purpose of this code is to collect these command line parameters
// and then use them to create an Options block.  This is then passed on to
// the template processor for execution.
//
var (
	// these need to be referenced as pointers
	inDir    = flag.String("i", "./", "input directory")
	justShow = flag.Bool("j", false, "display option settings and exit")
	outDir   = flag.String("o", "./", "output directory")
	testing  = flag.Bool("T", false, "this is a test run")
	verbose  = flag.Bool("v", false, "be talkative")
)

func main() {
	var (
		err         error
		nameWithExt []string // input file names with extensions
	)

	flag.Usage = Usage
	flag.Parse()
	fileNames := flag.Args()

	// FIXUPS ///////////////////////////////////////////////////////

	// XXX inDir must exist
	_, err = os.Stat(*inDir)

	// XXX if outDir does not exist, create it
	if err == nil {
		if _, err = os.Stat(*outDir); os.IsNotExist(err) {
			err = os.Mkdir(*outDir, 0755)
		}
	}
	if err == nil {
    	// SANITY CHECKS ////////////////////////////////////////////
    	if len(fileNames) == 0 {
    		err = NothingToDo
    	} else {
    		for i := 0; (err == nil) && (i < len(fileNames)); i++ {
    			name := fileNames[i]
				f := filepath.Join(*inDir, name)
				if _, err = os.Stat(f); os.IsNotExist(err) {
					if !strings.HasSuffix(f, ".md") {
						f = f + ".md"
						_, err = os.Stat(f)
					}
    			}
				if err != nil {
					err = SrcFileDoesNotExist
					fmt.Printf("%s does not exist\n", f)
				}
    			nameWithExt = append(nameWithExt, f)
    		}
    	}
    	// DISPLAY STUFF ////////////////////////////////////////////
    	if *verbose || *justShow {
    		fmt.Printf("inDir        = %v\n", *inDir)
    		fmt.Printf("justShow     = %v\n", *justShow)
    		fmt.Printf("outDir       = %s\n", *outDir)
    		fmt.Printf("testing      = %v\n", *testing)
    		fmt.Printf("verbose      = %v\n", *verbose)
    		if len(nameWithExt) > 0 {
    			fmt.Println("INFILES:")
    			for i := 0; i < len(nameWithExt); i++ {
    				fmt.Printf("%3d: %s\n", i, nameWithExt[i])
    			}
    		}
    	}
	}
	if err != nil {
		fmt.Printf("\nerror = %s\n", err.Error())
	}
	if err != nil || *justShow {
		return
	}

	var (
		doc *gm.Document
		in  io.Reader
		p   *gm.Parser
	)
	for i := 0; i < len(nameWithExt); i++ {
		inFile := nameWithExt[i]
		base := fileNames[i]
		outFile := filepath.Join(*outDir, base+".html")
		in, err = os.Open(inFile)
		options := gm.NewOptions(in, inFile, outFile, *testing, *verbose)
		if err == nil {
			p, err = gm.NewParser(options)
			if err == nil {
				doc, err = p.Parse()
				if err == io.EOF {
					err = nil
				}
				if err == nil {
					html := []byte(string(doc.GetHtml()))
					err = ioutil.WriteFile(outFile, html, 0666)
				}
			}
		}
		if err != nil {
			break
		}
	}

	if err != nil {
		fmt.Printf("\nerror processing input file(s): %s\n", err.Error())
	}
	return
}
