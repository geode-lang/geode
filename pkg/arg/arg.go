package arg

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Primary globally valid commands and arguments
var (
	App                   = kingpin.New("geode", "Compiler for the Geode Programming Language").Author("Nick Wanninger")
	BuildOutput           = App.Flag("output", "Output binary name.").Short('o').Default("a.out").String()
	Optimize              = App.Flag("optimize", "Enable full optimization").Short('O').Default("0").Int()
	PrintVerbose          = App.Flag("verbose", "Enable verbose printing").Short('v').Bool()
	StopAfterCompilation  = App.Flag("no-binary", "Stop after compilation").Short('c').Bool()
	DisableEmission       = App.Flag("no-emission", "Disable emission and only run through the syntax checking process").Bool()
	DisableRuntime        = App.Flag("no-runtime", "Disable calls to the runtime. Warning: garbage collector, etc will be gone. Most standard libraries will not work.").Bool()
	DisableStringDataCopy = App.Flag("no-dynamic-strings", "Disable the dynamic string copy and replace with static/constant .data section pointers").Bool()
	LinkerArgs            = App.Flag("linker-args", "Arguments to pass clang when linking object files").String()
	EmitASM               = App.Flag("asm", "Emit the asm of the program to the current directory. (will not produce binary)").Bool()
	EmitLLVM              = App.Flag("llvm", "Emit the llvm of the program to the current directory. (will not produce binary)").Bool()
	ShowLLVM              = App.Flag("show-llvm", "Print the llvm to stdout for debugging codegen").Short('S').Bool()
	EmitObject            = App.Flag("obj", "Emit the object file of the program to the current directory. (will not produce binary)").Bool()
	DumpScopeTree         = App.Flag("dump-scope-tree", "Dump a tree representation of the scope to stdout").Bool()
	ClangFlags            = App.Flag("clang-flags", "flags to pass into the clang compiler/linker").String()
	EnableDebug           = App.Flag("debug", "(NOT WORKING) Enable debug information").Short('g').Bool()
)

// Global arguments accessable throughout the program
var (
	VersionCMD = App.Command("version", "Display the version")

	BuildCMD   = App.Command("build", "Build an executable.")
	BuildInput = BuildCMD.Arg("input", "Geode source file or package").Default(".").String()

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

// Commands related to the pkg subcommand
var (
	PkgCMD  = App.Command("pkg", "Envoke the geode git package manager")
	PkgInit = PkgCMD.Flag("init", "initialize the package manager config").Bool()
)
