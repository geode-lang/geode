package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	// "gopkg.in/alecthomas/kingpin.v2"
	"github.com/jawher/mow.cli"
	"gitlab.com/nickwanninger/geode/pkg/gen"
	"gitlab.com/nickwanninger/geode/pkg/lexer"
)

// Usage will print the usage of the program
func Usage() {
	fmt.Println("Usage: geode <command> [options] <file>")
	fmt.Println("Options:")
	// app.UsageWriter(os.Stdout)
}

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

func main() {

	app := cli.App("geode", "A programming language")

	// Declare our first command, which is invocable with "uman list"
	app.Command("build", "Compile a geode source file", func(cmd *cli.Cmd) {

		cmd.Spec = "[-o] SOURCE"
		source := cmd.StringArg("SOURCE", "", "Source file to compile")
		output := cmd.StringOpt("o output", "main", "Binary output name")

		// Run this function when the command is invoked
		cmd.Action = func() {
			build(*source, *output)
		}
	})

	app.Run(os.Args)
	// // kingpin.MustParse(app.Parse(os.Args[1:]))
	// command := kingpin.MustParse(app.Parse(os.Args[1:]))

	// switch command {
	// case buildCommand.FullCommand():
	// 	filename, _ := resolveFileName(*buildInput)
	// 	build(filename, *buildOutput)
	// case buildCommand.FullCommand():
	// 	filename, _ := resolveFileName(*buildInput)
	// 	run(filename)
	// }

}

// func run(filename string) {
// 	if build(filename, "/tmp/geodeprogram") {
// 		cmd := exec.Command("/tmp/geodeprogram")
// 		cmd.Start()
// 		fmt.Println(cmd.Wait())
// 	}
// }

func build(filename string, output string) bool {
	if filename == "" {
		fmt.Println("No input files passed.")
	}

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

	// if *buildPrintLLVMIR {
	// 	fmt.Println(comp.GetLLVMIR())
	// 	return false
	// }

	comp.EmitModuleObject()
	compiled := comp.Compile()
	if !compiled {
		fmt.Println("Compilation failed. Please check the logs")
		return false
	}
	return true
}
