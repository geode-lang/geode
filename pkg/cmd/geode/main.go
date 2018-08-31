package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/geode-lang/geode/pkg/ast"
	"github.com/geode-lang/geode/pkg/info"
	"github.com/geode-lang/geode/pkg/util"
	"github.com/geode-lang/geode/pkg/util/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Some constants that represent the program in it's current compiled state
const (
	VERSION = "0.4.2"
	AUTHOR  = "Nick Wanninger"
)

var startTime time.Time

func main() {
	// :^)
	if runtime.GOOS == "windows" {
		log.Fatal("The Geode Compiler does not support windows at this time.")
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
		RunTests("./tests")

	case newTestCMD.FullCommand():
		CreateTestCMD()

	case cleanCMD.FullCommand():
		os.RemoveAll(buildDir)

	case versionCMD.FullCommand():
		fmt.Println(VERSION)
		os.Exit(0)
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

	program := ast.NewProgram()
	program.ParseDep("", "std:runtime")
	program.Entry = c.Input
	program.ParsePath(c.Input)
	program.TargetTripple = c.TargetTripple
	program.Congeal()

	options := ast.FunctionCompilationOptions{}
	program.CompileFunction("main", options)

	// for _, pkg := range program.Packages {
	// 	fmt.Println("============================")
	// 	for file := range pkg.Files {
	// 		fmt.Println(color.Green(file))
	// 	}
	// 	fmt.Println("============================")

	// 	buff := &bytes.Buffer{}
	// 	for _, node := range pkg.Nodes {
	// 		fmt.Fprintf(buff, "%s\n\n", node)
	// 	}

	// 	fmt.Println(buff)
	// }

	// os.Exit(0)

	// // Construct a linker object
	target := ast.BinaryTarget
	if *emitASM {
		target = ast.ASMTarget
	}

	linker := ast.NewLinker(*buildOutput)
	linker.SetTarget(target)
	linker.SetBuildDir(buildDir)
	linker.SetOutput(c.Output)
	linker.SetOptimize(*optimize)

	for _, clink := range program.CLinkages {
		linker.AddObject(clink)
	}

	if *dumpScopeTree {
		fmt.Println(program.Scope)
	}

	if *disableEmission {
		if *dumpResult {
			fmt.Println(program)
		}
		return
	}

	linker.SetDump(*dumpResult)

	linker.AddObject(program.Emit(buildDir))
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

	// The program exited with a failed code. So we need to exit with that same code.
	// This is because the run command should feel like just running the binary
	if err != nil {

		// pretty.Print(err)

		if exitError, ok := err.(*exec.ExitError); ok {

			waitStatus := exitError.Sys().(syscall.WaitStatus)
			os.Exit(waitStatus.ExitStatus())
		}
	}

	// The program exited safely, so we should too
	os.Exit(0)
}
