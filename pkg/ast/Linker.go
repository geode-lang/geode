package ast

import (
	"os"
	"path"
	"path/filepath"
	"strings"

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

		// We don't want to remove .c files
		if filepath.Ext(objFile) == ".c" {
			continue
		}

		os.Remove(objFile)

	}
}

// type error interface {
// 	Error() string
// }

// linkError is a trivial implementation of error.
type linkError struct {
	s string
}

func (e linkError) Error() string {
	return e.s
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
		// We want to only write intel syntax. AT&T Sucks
		linkArgs = append(linkArgs, "-S", "-masm=intel")
		// Compile each of the objects to a .s file.
		for _, obj := range l.objectPaths {
			// We only want to leave user generated files in the filesystem
			if strings.HasSuffix(obj, ".ll") && obj != "runtime.ll" {
				// Pull the extension of the object file
				ext := path.Ext(obj)
				// Replace it with .s
				filename = obj[0:len(obj)-len(ext)] + ".s"
				// set the output to that of the .s file
				asmArgs := append(linkArgs, "-o", filename, obj)
				// run the compile to asm
				util.RunCommand(linker, asmArgs...)
			}

		}
		return
	}

	// Append input files to the end of the command
	linkArgs = append(linkArgs, l.objectPaths...)

	// set the output filename
	linkArgs = append(linkArgs, "-o", filename)

	util.RunCommand(linker, linkArgs...)
}
