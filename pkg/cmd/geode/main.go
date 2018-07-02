package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/nickwanninger/geode/pkg/ast"
	"github.com/nickwanninger/geode/pkg/lexer"
	"github.com/nickwanninger/geode/pkg/util"
	"github.com/nickwanninger/geode/pkg/util/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Some constants that represent the program in it's current compiled state
const (
	VERSION = "0.0.2"
	AUTHOR  = "Nick Wanninger"
)

var startTime time.Time

func main() {

	_, clangError := util.RunCommand("clang", "--version")
	if clangError != nil {
		log.Fatal("Unable to find a clang install in your path. Please install clang and add it to your path\n")
	}
	spew.Config.DisableMethods = true
	startTime = time.Now()

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch command {
	case buildCMD.FullCommand():
		context := NewContext(*buildInput, *buildOutput)
		context.Build()

	case runCMD.FullCommand():
		dir, err := ioutil.TempDir("", "geode")
		if err != nil {
			log.Fatal("Unable to produce tmp directory for `geode run` executable\n")
		}
		out := dir + "/exe"
		context := NewContext(*runInput, out)
		context.Build()
		context.Run(*runArgs)

	case testCMD.FullCommand():
		RunTests(*testDir)
	}

}

// Context contains information for this compilation
type Context struct {
	Input  string
	Output string
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

	src, err := lexer.NewSourcefile(c.Input)
	if err != nil {
		log.Fatal("Unable to construct a source file.\n")
	}
	err = src.ResolveFile(c.Input)
	if err != nil {
		log.Fatal("Unable to read file %s into sourcefile structure: %s\n", c.Input, err)
	}
	path := strings.Split(c.Input, "/")
	root := ast.NewPackage(path[len(path)-1], src)
	pkgs := make([]*ast.Package, 0)
	for pkg := range root.Parse() {
		pkgs = append(pkgs, pkg)
	}

	// Construct a linker object
	target := ast.BinaryTarget
	if *emitASM {
		target = ast.ASMTarget
	}

	linker := ast.NewLinker(*buildOutput)
	linker.SetTarget(target)
	linker.SetOutput(c.Output)
	linker.SetOptimize(*optimize)

	// Loop over the compilers and generate to .ll files
	for c := range root.Compile() {
		obj := c.Emit()
		linker.AddObject(obj)
		for _, link := range c.CLinkages {
			linker.AddObject(link)
		}
	}
	if *emitLLVM {
		log.Debug("%s\n", root)
	}

	if err != nil {
		log.Fatal("Unable to scan for c libraries\n")
	}

	linker.Run()

	linker.Cleanup()
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
