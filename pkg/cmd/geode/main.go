package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"gitlab.com/nickwanninger/geode/pkg/gen"
	"gitlab.com/nickwanninger/geode/pkg/lexer"
	"gitlab.com/nickwanninger/geode/pkg/util/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Some constants that represent the program in it's current compiled state
const (
	VERSION = "0.0.1"
	AUTHOR  = "Nick Wanninger"
)

var startTime time.Time

func main() {
	spew.Config.DisableMethods = true
	startTime = time.Now()

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch command {
	case buildCMD.FullCommand():
		filename, _ := resolveFileName(*buildInput)
		context := NewContext(filename, *buildOutput)
		context.Build()

	case runCMD.FullCommand():
		filename, _ := resolveFileName(*runInput)
		context := NewContext(filename, "/tmp/"+time.Now().String())
		context.Build()
		context.Run(*runArgs)
	}
}

// if the filename passed in is a folder, look in that folder for a main.g
// if the filename is not, look for a file matching that filename, but with a .g extension
func resolveFileName(filename string) (string, error) {
	// Grab the stats of the file
	stats, err := os.Stat(filename)

	// If there was an error (file doesnt exist)
	if err != nil {
		// Try resolving the filename with .g extension
		if !strings.HasSuffix(filename, ".g") {
			return resolveFileName(filename + ".g")
		}
		// There was no file by that name, so we fail
		return "", fmt.Errorf("fatal error: No such file or directory %s", filename)
	}
	if stats.IsDir() {
		return resolveFileName(filename + "/main.g")
	}

	return filename, nil
}

// Context contains information for this compilation
type Context struct {
	// Searchpaths []string

	Input  string
	Output string

	// moduleLookup *ast.ModuleLookup
	// modules      []*ast.Module
	// depGraph     *ast.DependencyGraph

	// modulesToRead []*ast.ModuleName
}

// NewContext constructs a new context and returns a pointer to it
func NewContext(in string, out string) *Context {

	if in == "" {
		log.Fatal("Failed to create context, no input file passed\n")
	}
	res := &Context{}
	res.Input = in
	res.Output = out
	return res
}

// Build some context into a binary file
func (c *Context) Build() {
	data, err := ioutil.ReadFile(c.Input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Build a new lexer. This contans the methods required to parse some string of data into
	lexer := lexer.NewLexer()
	// Run the lexer concurrently
	go lexer.Lex([]byte(string(gen.RuntimeSource) + string(data)))

	tokens := lexer.Tokens

	nodes := gen.Parse(tokens)

	comp := gen.NewCompiler(c.Input, c.Output)

	for node := range nodes {
		node.Codegen(comp.RootScope.SpawnChild(), comp)
	}

	if *emitLLVM {
		log.Debug("%s\n", comp.GetLLVMIR())
	}
	comp.EmitModuleObject()

	target := gen.BinaryTarget
	if *emitASM {
		target = gen.ASMTarget
	}

	compiled := comp.Compile(target)
	if !compiled {
		log.Fatal("Compilation failed. Please check the logs\n")
	}

}

// Run a context with a given set of arguments
func (c *Context) Run(args []string) {
	cmd := exec.Command(c.Output, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	// The program exited with a failed code. So we need to exit with that same code.
	// This is because the run command should feel like just running the binary
	if err != nil {
		exitCodeString := strings.Replace(err.Error(), "exit status ", "", -1)
		exitCode, _ := strconv.Atoi(exitCodeString)
		os.Exit(exitCode)
	}
	// The program exited safely, so we should too
	os.Exit(0)

}
