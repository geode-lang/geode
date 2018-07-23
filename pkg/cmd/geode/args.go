package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app             = kingpin.New("geode", "Compiler for the Geode Programming Language").Version(VERSION).Author(AUTHOR)
	dumpResult      = app.Flag("dump", "Print either llvm or ASM code after compiled (llvm by default, asm if --asm is passed)").Short('S').Bool()
	buildOutput     = app.Flag("output", "Output binary name.").Short('o').Default("a.out").String()
	optimize        = app.Flag("optimize", "Enable full optimization").Short('O').Bool()
	printVerbose    = app.Flag("verbose", "Enable verbose printing").Short('v').Bool()
	disableEmission = app.Flag("disable-emission", "Disable emission and only go through the syntax checking process").Bool()
	// logLevel = app.Flag("loglevel", "Set the level of logging to show").Default("info").Enum("info", "verbose")

	buildCMD   = app.Command("build", "Build an executable.").Alias("b")
	buildInput = buildCMD.Arg("input", "Geode source file or package").Default(".").String()
	emitASM    = buildCMD.Flag("asm", "Set the target to .s asm files with intel syntax instead of a single binary.").Bool()

	runCMD   = app.Command("run", "Build and run an executable, clean up afterwards").Alias("r").Default()
	runInput = runCMD.Arg("input", "Geode source file or package").String()
	runArgs  = runCMD.Arg("args", "Arguments to be passed into the program after building").Strings()

	testCMD = app.Command("test", "Run tests").Alias("t")
	testDir = testCMD.Arg("dir", "Test Directory").Default("./tests").String()

	cleanCMD = app.Command("clean", "Remove the hidden build directory")
)
