package main

import (
	"flag"
	"fmt"

	"io/ioutil"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"gitlab.com/nickwanninger/geode/pkg/gen"
	"gitlab.com/nickwanninger/geode/pkg/parser"
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

// FlagConfig -
type FlagConfig struct {
	OutFile       *string
	OptimizeLevel *int
	PrintTokens   *bool
	PrintAst      *bool
	PrintASTJson  *bool
	PrintLLVMIR   *bool
	Args          []string
}

func parseFlags() *FlagConfig {
	c := &FlagConfig{}
	flag.Usage = Usage
	c.OutFile = flag.String("o", "out", "The output filename")
	c.OptimizeLevel = flag.Int("opt", 3, "add some optimization passes")
	c.PrintTokens = flag.Bool("tok", false, "Print tokens as they are parsed (for debugging)")
	c.PrintAst = flag.Bool("ast", false, "print abstract syntax tree (for debugging)")
	c.PrintASTJson = flag.Bool("json", false, "If true, the ast will be dumped to the console as json instead of raw")
	c.PrintLLVMIR = flag.Bool("s", false, "dump LLVM IR to console and don't compile")
	flag.Parse()
	c.Args = flag.Args()

	return c
}

func main() {

	spew.Config.Indent = "  "
	spew.Config.SortKeys = true
	spew.Config.SpewKeys = true

	config := parseFlags()

	if flag.NArg() == 0 {
		fmt.Println("No .act files or folders containing .act files provided.")
		Usage()
		return
	}
	// Get the filename with the resolver method. This allows a user to enter `.` and the compiler will assume they meant `./main.act`
	// it also allows the user to enter `foo` and the compiler will attempt to compile `foo.act`
	filename, ferr := resolveFileName(config.Args[0])
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
	go lexer.Lex([]byte(string(gen.RuntimeSource) + string(data)))

	tokens := lexer.Tokens

	if *config.PrintTokens {
		tokens = parser.DumpTokens(lexer.Tokens)
	}

	nodes := gen.Parse(tokens)

	if *config.PrintAst {
		nodes = gen.DumpTree(nodes, *config.PrintASTJson)
	}

	comp := gen.NewCompiler(filename, *config.OutFile)
	for node := range nodes {
		node.Codegen(comp.RootScope.SpawnChild(), comp)
	}
	// compiler.Optimize()

	if *config.PrintLLVMIR {
		fmt.Println(comp.GetLLVMIR())
		return
	}

	comp.EmitModuleObject()
	compiled := comp.Compile()
	if !compiled {
		fmt.Println("Compilation failed. Please check the logs")
	}
}
