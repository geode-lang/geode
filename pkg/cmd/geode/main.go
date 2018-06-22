package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jawher/mow.cli"
	"gitlab.com/nickwanninger/geode/pkg/gen"
	"gitlab.com/nickwanninger/geode/pkg/lexer"
)

// Some constants that represent the program in it's current compiled state
const (
	VERSION = "0.0.1"
	AUTHOR  = "Nick Wanninger"
)

// if the filename passed in is a folder, look in that folder for a main.g
// if the filename is not, look for a file matching that filename, but with a .g extension
func resolveFileName(filename string) (string, error) {
	// Grab the stats of the file
	stats, err := os.Stat(filename)

	// If there was an error (file doesnt exist)
	if err != nil {
		// Try resolving the filename with .g extension
		if !strings.HasSuffix(filename, ".gd") {
			return resolveFileName(filename + ".gd")
		}
		// There was no file by that name, so we fail
		return "", fmt.Errorf("fatal error: No such file or directory %s", filename)
	}
	if stats.IsDir() {
		return resolveFileName(filename + "/main.g")
	}

	return filename, nil
}

var (
	printLLVM *bool
)

func main() {

	app := cli.App("geode", "A programming language by Nick Wanninger")
	app.Version("version", VERSION)

	app.Spec = "[-S]"
	printLLVM = app.BoolOpt("S", false, "Print the LLVM IR")
	mainSpec := "[-o] SOURCE"

	// Declare our first command, which is invocable with "uman list"
	app.Command("build", "Compile a geode source file", func(cmd *cli.Cmd) {

		cmd.Spec = mainSpec
		source := cmd.StringArg("SOURCE", "", "Source file to compile")
		output := cmd.StringOpt("o output", "main", "Binary output name")
		// Run this function when the command is invoked
		cmd.Action = func() {
			build(*source, *output)
		}
	})

	//
	app.Command("run", "Run a geode source file", func(cmd *cli.Cmd) {

		cmd.Spec = "SOURCE [ARGS...]"
		source := cmd.StringArg("SOURCE", "", "Source file to compile")
		args := cmd.StringsArg("ARGS", nil, "Arguments to pass into geode program")
		// Run this function when the command is invoked
		cmd.Action = func() {
			run(*source, *args)
		}
	})

	app.Run(os.Args)
}

func run(filename string, args []string) {
	outFile := "/tmp/geodeprogram"
	if build(filename, outFile) {
		cmd := exec.Command(outFile, args...)
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
	} else {
		fmt.Printf("Failed to run %q because the build failed", filename)
	}
}

func build(filename string, output string) bool {
	if filename == "" {
		fmt.Println("No input files passed.")
	}

	filename, _ = resolveFileName(filename)

	data, err := ioutil.ReadFile(filename)
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

	comp := gen.NewCompiler(filename, output)

	for node := range nodes {
		node.Codegen(comp.RootScope.SpawnChild(), comp)
	}

	if *printLLVM {
		fmt.Println("=====")
		fmt.Println(comp.GetLLVMIR())
		fmt.Println("=====")
	}
	// fmt.Println(comp.GetLLVMIR())
	comp.EmitModuleObject()
	compiled := comp.Compile()
	if !compiled {
		fmt.Println("Compilation failed. Please check the logs")
		return false
	}
	return true
}
