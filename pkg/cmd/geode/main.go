package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
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

	spew.Config.DisableMethods = true
	startTime = time.Now()
	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	buildDir := fmt.Sprintf(".geode_build/")

	log.PrintVerbose = *printVerbose

	clangVersion, clangError := util.RunCommand("clang", "--version")
	if clangError != nil {
		log.Fatal("Unable to find a clang install in your path. Please install clang and add it to your path\n")
	}
	log.Verbose("Clang Version: %s\n", clangVersion)
	log.Verbose("Building to %s...\n", buildDir)

	switch command {
	case buildCMD.FullCommand():
		log.Timed("Compilation", func() {
			context := NewContext(*buildInput, *buildOutput)
			context.Build(buildDir)
		})

	case runCMD.FullCommand():
		out := path.Join(buildDir, "a.out")
		context := NewContext(*runInput, out)
		context.Build(buildDir)
		context.Run(*runArgs, buildDir)

	case testCMD.FullCommand():
		RunTests(*testDir)

	case cleanCMD.FullCommand():
		util.RunCommand("rm", "-rf", buildDir)
	}
	duration := time.Since(startTime)
	log.Verbose("Total time taken: %s\n", duration)
	if duration >= time.Second {
		log.Verbose("!! Compilation took more than one second. !!\n")
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
func (c *Context) Build(buildDir string) {

	src, err := lexer.NewSourcefile(c.Input)
	if err != nil {
		log.Fatal("Unable to construct a source file.\n")
	}
	err = src.ResolveFile(c.Input)
	if err != nil {
		log.Fatal("Unable to read file %s into sourcefile structure: %s\n", c.Input, err)
	}
	path := strings.Split(c.Input, "/")
	rootPackage := ast.NewPackage(path[len(path)-1], src)
	pkgs := make([]*ast.Package, 0)
	primaryTree := make([]ast.Node, 0)

	for pkg := range rootPackage.Parse() {
		log.Debug("Added package %s\n", pkg.Name)
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
	log.Timed("llvm emission", func() {
		for c := range rootPackage.Compile() {
			log.Debug("Emitting pkg %s with namespace %s\n", c.Name, c.NamespaceName)
			obj := c.Emit(buildDir)
			primaryTree = append(primaryTree, c.Nodes...)
			linker.AddObject(obj)
			for _, link := range c.CLinkages {
				log.Debug("Added c linkage %s\n", link)
				linker.AddObject(link)
			}
		}
	})

	if *emitLLVM {
		log.Printf("%s\n", rootPackage)
	}

	log.Timed("Linking", func() {
		linker.Run()
	})

}

// Run a context with a given set of arguments
func (c *Context) Run(args []string, buildDir string) {
	cmd := exec.Command(c.Output, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	// os.RemoveAll(buildDir) // Clean up the build dir after running

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
