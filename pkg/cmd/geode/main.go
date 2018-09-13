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

	"github.com/geode-lang/geode/pkg/arg"
	"github.com/geode-lang/geode/pkg/ast"
	"github.com/geode-lang/geode/pkg/info"
	"github.com/geode-lang/geode/pkg/util"
	"github.com/geode-lang/geode/pkg/util/log"
)

// Some constants that represent the program in it's current compiled state
const (
	VERSION = "0.6.5"
	AUTHOR  = "Nick Wanninger"
)

var startTime time.Time

func main() {
	// :^)
	if runtime.GOOS == "windows" {
		log.Fatal("The Geode Compiler does not support Windows at this time.")
	}

	// ast.TestNewParser()

	startTime = time.Now()
	command := arg.Parse()
	home := util.HomeDir()
	buildDir := path.Join(home, ".geode/build/")

	log.PrintVerbose = *arg.PrintVerbose

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
	case arg.BuildCMD.FullCommand():
		log.Timed("Compilation", func() {
			context := NewContext(*arg.BuildInput, *arg.BuildOutput)
			context.TargetTripple = targetTripple
			context.Build(buildDir)
		})

	case arg.RunCMD.FullCommand():
		out := path.Join(buildDir, "a.out")
		context := NewContext(*arg.RunInput, out)
		context.TargetTripple = targetTripple
		context.Build(buildDir)
		context.Run(*arg.RunArgs, buildDir)

	case arg.TestCMD.FullCommand():
		RunTests("./tests")

	case arg.NewTestCMD.FullCommand():
		CreateTestCMD()

	case arg.CleanCMD.FullCommand():
		os.RemoveAll(buildDir)

	case arg.VersionCMD.FullCommand():
		fmt.Println(VERSION)
		os.Exit(0)

	// case arg.PkgCMD.FullCommand():
	// 	geodepkg.HandleCommand()
	// 	os.Exit(0)

	case arg.InfoCMD.FullCommand():
		log.Timed("information gathering", func() {
			context := NewContext(*arg.InfoInput, "/tmp/geodeinfooutput")
			*arg.DisableEmission = true
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

	if !*arg.DisableRuntime {
		program.ParseDep("", "std:runtime")
	}

	program.Entry = c.Input

	if _, err := os.Stat(c.Input); os.IsNotExist(err) {
		fmt.Printf("The file %q could not be found.\n", c.Input)
		os.Exit(-1)
	}

	program.ParsePath(c.Input)
	program.TargetTripple = c.TargetTripple

	_, err := program.Congeal()
	if err != nil {
		log.Fatal("%s\n", err)
	}

	options := ast.FunctionCompilationOptions{}
	main, err := program.GetFunction("main", options)
	if err != nil {
		log.Fatal("%s\n", err)
	}
	if main == nil {
		log.Fatal("No function `main` found in compilation.\n")
	}

	if *arg.ShowLLVM {
		fmt.Println(program)
	}

	// // Construct a linker object
	target := ast.BinaryTarget
	if *arg.EmitASM {
		target = ast.ASMTarget
	}

	linker := ast.NewLinker(*arg.BuildOutput)
	linker.SetTarget(target)
	linker.SetBuildDir(buildDir)
	linker.SetOutput(c.Output)
	linker.SetOptimize(*arg.Optimize)

	for _, clink := range program.CLinkages {
		linker.AddObject(clink)
	}

	if *arg.DumpScopeTree {
		fmt.Println(program.Scope)
	}

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
