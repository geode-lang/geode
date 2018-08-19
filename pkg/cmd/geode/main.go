package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/geode-lang/geode/pkg/ast"
	"github.com/geode-lang/geode/pkg/info"
	"github.com/geode-lang/geode/pkg/lexer"
	"github.com/geode-lang/geode/pkg/util"
	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/geode-lang/llvm/ir"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Some constants that represent the program in it's current compiled state
const (
	VERSION = "0.2.0"
	AUTHOR  = "Nick Wanninger"
)

var startTime time.Time

func main() {
	if runtime.GOOS == "windows" {
		log.Fatal("Geode does not support windows at this time.")
	}

	startTime = time.Now()
	command := kingpin.MustParse(app.Parse(os.Args[1:]))
	home := util.HomeDir()
	buildDir := path.Join(home, ".geode/build/")

	log.PrintVerbose = *printVerbose

	clangVersion, clangError := util.RunCommand("clang", "-v")
	if clangError != nil {
		log.Fatal("Unable to find a clang install in your path. Please install clang and add it to your path\n")
	}

	clangVersionLines := strings.Split(string(clangVersion), "\n")
	targetTripple := ""

	for _, line := range clangVersionLines {
		if strings.HasPrefix(line, "Target: ") {
			targetTripple = strings.Replace(line, "Target: ", "", 1)
		}
	}

	log.Verbose("Clang Version: %s\n", clangVersion)
	log.Verbose("Building to %s...\n", buildDir)

	switch command {
	case buildCMD.FullCommand():
		log.Timed("Compilation", func() {
			context := NewContext(*buildInput, *buildOutput)
			context.TargetTripple = targetTripple
			context.Build(buildDir)
		})

	case runCMD.FullCommand():
		out := path.Join(buildDir, "a.out")
		context := NewContext(*runInput, out)
		context.TargetTripple = targetTripple
		context.Build(buildDir)
		context.Run(*runArgs, buildDir)

	case testCMD.FullCommand():
		RunTests(*testDir)

	case cleanCMD.FullCommand():
		os.RemoveAll(buildDir)

	case infoCMD.FullCommand():
		log.Timed("information gathering", func() {
			context := NewContext(*infoInput, "/tmp/geodeinfooutput")
			*disableEmission = true
			context.TargetTripple = targetTripple
			context.Build(buildDir)
			info.DumpJSON()
		})
	}
	duration := time.Since(startTime)
	log.Verbose("Total time taken: %s\n", duration)
	if duration >= time.Second {
		log.Verbose("!! Compilation took more than one second. !!\n")
	}
}

// Context contains information for this compilation
type Context struct {
	Input         string
	Output        string
	TargetTripple string
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
	src.Preprocess()
	if err != nil {
		log.Fatal("Unable to read file %s into sourcefile structure: %s\n", c.Input, err)
	}
	path := strings.Split(c.Input, "/")
	scope := ast.NewScope()

	nodeNamespaces := make(map[string]*[]ast.Node, 0)
	rootPackage := ast.NewPackage(path[len(path)-1], src, scope, nodeNamespaces)
	pkgs := make([]*ast.Package, 0)
	primaryTree := make([]ast.Node, 0)

	runtime := ast.GetRuntime(scope, nodeNamespaces)

	rootPackage.Inject(runtime)

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
	linker.SetBuildDir(buildDir)
	linker.SetOutput(c.Output)
	linker.SetOptimize(*optimize)

	module := ir.NewModule()

	// Loop over the compilers and generate to .ll files
	log.Timed("llvm emission", func() {
		for c := range rootPackage.Compile(module, c.TargetTripple) {
			primaryTree = append(primaryTree, c.Nodes...)
			if !*disableEmission {
				log.Debug("Compiled pkg %s with namespace %s\n", c.Name, c.NamespaceName)
				for _, link := range c.CLinkages {
					log.Debug("Added c linkage %s\n", link)
					linker.AddObject(link)
				}
			}
		}
	})

	if *dumpScopeTree {
		fmt.Println(scope)
	}

	for _, n := range primaryTree {
		info.AddNode(n)
	}

	if *disableEmission {
		if *dumpResult {
			fmt.Println(rootPackage)
		}
		return
	}

	linker.SetDump(*dumpResult)

	linker.AddObject(rootPackage.Emit(buildDir))
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
