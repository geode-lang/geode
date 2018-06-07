package main

import (
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	// "github.com/nickwanninger/act/pkg/ast"
	"github.com/nickwanninger/act/pkg/parser"
	"io/ioutil"
	"os"
	"strings"
)

var (
	batch       = flag.Bool("b", false, "batch (non-interactive) mode")
	optimized   = flag.Bool("opt", true, "add some optimization passes")
	printTokens = flag.Bool("tok", false, "Print tokens")
	printAst    = flag.Bool("ast", false, "print abstract syntax tree")
	printLLVMIR = flag.Bool("s", false, "print LLVM generated code")
)

// Usage will print the usage of the program
func Usage() {
	fmt.Println("Usage: act [options] <file>")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

// if the filename passed in is a folder, look in that folder for a main.act
// if the filename is not, look for a file matching that filename, but with a .act extension
func resolveFileName(filename string) (string, error) {
	// Grab the stats of the file
	stats, err := os.Stat(filename)

	// If there was an error (file doesnt exist)
	if err != nil {
		// Try resolving the filename with .act extension
		if !strings.HasSuffix(filename, ".act") {
			return resolveFileName(filename + ".act")
		}
		// There was no file by that name, so we fail
		return "", fmt.Errorf("fatal error: No such file or directory %s", filename)
	}
	if stats.IsDir() {
		return resolveFileName(filename + "/main.act")
	}

	return filename, nil
}

func main() {

	spew.Config.Indent = "        "
	spew.Config.SortKeys = true
	flag.Usage = Usage
	flag.Parse()
	args := flag.Args()

	lexer := parser.NewLexer()

	if flag.NArg() == 0 {
		fmt.Println("No .act files provided.")
		Usage()
		return
	}

	filename, ferr := resolveFileName(args[0])
	if ferr != nil {
		fmt.Println(ferr)
		return
	}

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println(err)
	}

	go lexer.Lex(data)

	// nodes := ast.Parse(lexer.Tokens)

	// if *printAst {
	// 	nodes = ast.DumpTree(nodes)
	// }

	for true {
		t, ok := <-lexer.Tokens
		if !ok {
			break
		}

		fmt.Println(t)
	}
}
