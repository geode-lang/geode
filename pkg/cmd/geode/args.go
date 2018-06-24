package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("geode", "Compiler for the Geode Programming Language").Version(VERSION).Author(AUTHOR)

	// logLevel = app.Flag("loglevel", "Set the level of logging to show").Default("info").Enum("info", "verbose")

	buildCMD    = app.Command("build", "Build an executable.")
	buildOutput = buildCMD.Flag("output", "Output binary name.").Short('o').Default("main").String()
	emitLLVM    = buildCMD.Flag("emit-llvm", "Emit LLVM to stdout").Short('S').Bool()
	buildInput  = buildCMD.Arg("input", "Ark source file or package").String()
)
