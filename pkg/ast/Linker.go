package ast

import (
	"os"
	"path"

	"github.com/nickwanninger/geode/pkg/util"
)

// CompileTarget is a target to build a binary for
type CompileTarget int

// Some default targets
const (
	ASMTarget CompileTarget = iota
	BinaryTarget
)

// Linker is an instance that can link several
// object files into a single binary or other format
type Linker struct {
	output      string
	target      CompileTarget
	objectPaths []string
	optimize    bool
}

// NewLinker constructs a linker with an outpu
func NewLinker(output string) *Linker {
	l := &Linker{}
	l.output = output
	return l
}

// AddObject appends an object path to a linker
func (l *Linker) AddObject(path string) {
	l.objectPaths = append(l.objectPaths, path)
}

// SetTarget sets the target output format of the linker
// be it binary, asm, etc...
func (l *Linker) SetTarget(t CompileTarget) {
	l.target = t
}

// SetOutput -
func (l *Linker) SetOutput(path string) {
	l.output = path
}

// SetOptimize -
func (l *Linker) SetOptimize(o bool) {
	l.optimize = o
}

// Cleanup removes all the
func (l *Linker) Cleanup() {
	for _, objFile := range l.objectPaths {
		os.Remove(objFile)
	}
}

// Run a list of objects through a linker and build
// into a single outfile with the given target
func (l *Linker) Run() {
	linker := "clang"
	linkArgs := make([]string, 0)

	linkArgs = append(linkArgs, "-lm", "-lc")

	filename := l.output

	if l.optimize {
		linkArgs = append(linkArgs, "-O3")
	}

	if l.target == ASMTarget {
		linkArgs = append(linkArgs, "-S", "-masm=intel")

		for _, obj := range l.objectPaths {
			ext := path.Ext(obj)
			filename = obj[0:len(obj)-len(ext)] + ".s"
			asmArgs := append(linkArgs, "-o", filename, obj)
			util.RunCommand(linker, asmArgs...)
		}
		return
	}

	// Append input files to the end of the command
	linkArgs = append(linkArgs, l.objectPaths...)

	// set the output filename
	linkArgs = append(linkArgs, "-o", filename)

	util.RunCommand(linker, linkArgs...)
}
