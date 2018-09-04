package arg

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Primary globally valid commands and arguments
var (
	App                   = kingpin.New("geode", "Compiler for the Geode Programming Language").Author("Nick Wanninger")
	BuildOutput           = App.Flag("output", "Output binary name.").Short('o').Default("a.out").String()
	Optimize              = App.Flag("optimize", "Enable full optimization").Short('O').Int()
	PrintVerbose          = App.Flag("verbose", "Enable verbose printing").Short('v').Bool()
	StopAfterCompilation  = App.Flag("no-binary", "Stop after compilation").Short('c').Bool()
	DisableEmission       = App.Flag("no-emission", "Disable emission and only run through the syntax checking process").Bool()
	DisableStringDataCopy = App.Flag("no-dynamic-strings", "Disable the dynamic string copy and replace with static/constant .data section pointers").Bool()
	LinkerArgs            = App.Flag("linker-args", "Arguments to pass clang when linking object files").String()
)

// Global arguments accessable throughout the program
var (
	VersionCMD = App.Command("version", "Display the version")

	BuildCMD      = App.Command("build", "Build an executable.")
	BuildInput    = BuildCMD.Arg("input", "Geode source file or package").Default(".").String()
	EmitASM       = BuildCMD.Flag("asm", "Set the target to .s asm files with intel syntax instead of a single binary.").Bool()
	EmitLLVM      = BuildCMD.Flag("llvm", "Set the target to a single .ll file in the current directory").Bool()
	DumpScopeTree = BuildCMD.Flag("dump-scope-tree", "Dump a tree representation of the scope to stdout").Bool()

	RunCMD   = App.Command("run", "Build and run an executable, clean up afterwards").Default()
	RunInput = RunCMD.Arg("input", "Geode source file or package").String()
	RunArgs  = RunCMD.Arg("args", "Arguments to be passed into the program after building").Strings()

	TestCMD = App.Command("test", "Run tests in the ./tests/ directory")

	NewTestCMD  = App.Command("new-test", "Create a new test")
	NewTestName = NewTestCMD.Arg("name", "the name of the test").Required().String()

	CleanCMD = App.Command("clean", "Remove the hidden build directory")

	InfoCMD   = App.Command("info", "Get information about a program (does not compile, just lexes and parses)")
	InfoInput = InfoCMD.Arg("input", "Geode source file or package").String()
)

// Parse returns the kingpin command returned by kingpin.MustParse
func Parse() string {
	return kingpin.MustParse(App.Parse(os.Args[1:]))
}
