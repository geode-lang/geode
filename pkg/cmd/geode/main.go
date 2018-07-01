package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/nickwanninger/geode/pkg/ast"
	"github.com/nickwanninger/geode/pkg/lexer"
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

	switch command {
	case buildCMD.FullCommand():
		filename, _ := resolveFileName(*buildInput)
		context := NewContext(filename, *buildOutput)
		context.Build()

	case runCMD.FullCommand():
		filename, _ := resolveFileName(*runInput)
		dir, err := ioutil.TempDir("", "geode")
		if err != nil {
			log.Fatal("Unable to produce tmp directory for `geode run` executable\n")
		}
		out := dir + "/exe"
		context := NewContext(filename, out)
		context.Build()
		context.Run(*runArgs)

	case testCMD.FullCommand():
		RunTests(*testDir)
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
	err = src.LoadFile(c.Input)
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
	}
	if *emitLLVM {
		log.Debug("%s\n", root)
	}

	gopath := os.Getenv("GOPATH")
	clibpath := gopath + "/src/github.com/nickwanninger/geode/lib/"

	var clibFiles []string

	// clibFiles = append(clibFiles, clibpath+"lib.c")

	err = filepath.Walk(clibpath, func(path string, info os.FileInfo, e error) error {
		if filepath.Ext(path) == ".c" {
			clibFiles = append(clibFiles, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal("Unable to scan for c libraries\n")
	}
	for _, path := range clibFiles {
		linker.AddObject(path)
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
