package main

import (
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/nickwanninger/act/pkg/ast"
	"github.com/nickwanninger/act/pkg/parser"
	// "github.com/nickwanninger/act/pkg/types"
	"io/ioutil"
	"os"
	"strings"
)

var (
	outFile       = flag.String("o", "a.out", "The output filename")
	optimizeLevel = flag.Int("opt", 3, "add some optimization passes")
	printTokens   = flag.Bool("tok", false, "Print tokens as they are parsed (for debugging)")
	printAst      = flag.Bool("ast", false, "print abstract syntax tree (for debugging)")
	printASTJson  = flag.Bool("json", false, "If true, the ast will be dumped to the console as json instead of raw")
	printLLVMIR   = flag.Bool("s", false, "print LLVM generated code")
)

// Usage will print the usage of the program
func Usage() {
	fmt.Println("Usage: act [options] <file>")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

//
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

	spew.Config.Indent = "  "
	spew.Config.SortKeys = true
	spew.Config.SpewKeys = true

	flag.Usage = Usage
	flag.Parse()
	// Pull the other arguments from the list of args
	// these come from after the arguments parsed abov
	// that allow the user to configure the compiler
	args := flag.Args()

	if flag.NArg() == 0 {
		fmt.Println("No .act files or folders containing .act files provided.")
		Usage()
		return
	}
	// Get the filename with the resolver method. This allows a user to enter `.` and the compiler will assume they meant `./main.act`
	// it also allows the user to enter `foo` and the compiler will attempt to compile `foo.act`
	filename, ferr := resolveFileName(args[0])
	if ferr != nil {
		fmt.Println(ferr)
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Build a new lexer. This contans the methods required to parse some string of data into
	lexer := parser.NewLexer()
	// Run the lexer concurrently
	go lexer.Lex(data)

	tokens := lexer.Tokens

	if *printTokens {
		tokens = parser.DumpTokens(lexer.Tokens)
	}

	nodes := ast.Parse(tokens)

	if *printAst {
		nodes = ast.DumpTree(nodes, *printASTJson)
	}
	for node := range nodes {
		node.Codegen(ast.GetRootScope().SpawnChild())
	}
	ast.Optimize()
	ast.EmitModuleObject()
	fmt.Println(ast.GetLLVMIR())
}
