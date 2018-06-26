package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app         = kingpin.New("geode", "Compiler for the Geode Programming Language").Version(VERSION).Author(AUTHOR)
	emitLLVM    = app.Flag("emit-llvm", "Leave LLVM IR in the filesystem").Short('S').Bool()
	buildOutput = app.Flag("output", "Output binary name.").Short('o').Default("a.out").String()
	// logLevel = app.Flag("loglevel", "Set the level of logging to show").Default("info").Enum("info", "verbose")

	buildCMD   = app.Command("build", "Build an executable.")
	buildInput = buildCMD.Arg("input", "Geode source file or package").Required().String()
	emitASM    = buildCMD.Flag("asm", "Build to asm instead of a binary").Bool()

	runCMD   = app.Command("run", "Build and run an executable, clean up afterwards")
	runInput = runCMD.Arg("input", "Geode source file or package").String()
	runArgs  = runCMD.Arg("args", "Arguments to be passed into the program after building").Strings()
)
