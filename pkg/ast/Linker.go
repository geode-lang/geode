package ast

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/geode-lang/geode/pkg/util"
	"github.com/geode-lang/geode/pkg/util/log"
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
	dump        bool // dump result of compilation to stdout
}

// NewLinker constructs a linker with an outpu
func NewLinker(output string) *Linker {
	l := &Linker{}
	l.output = output
	return l
}

// HasObject returns if the linker already has a path
func (l *Linker) HasObject(path string) bool {
	for _, p := range l.objectPaths {
		if p == path {
			return true
		}
	}
	return false
}

// AddObject appends an object path to a linker
func (l *Linker) AddObject(path string) {
	if !l.HasObject(path) {
		l.objectPaths = append(l.objectPaths, path)
	}

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

// SetDump -
func (l *Linker) SetDump(o bool) {
	l.dump = o
}

// Cleanup removes all the
func (l *Linker) Cleanup() {
	for _, objFile := range l.objectPaths {
		ext := filepath.Ext(objFile)
		// We only remove ll and s files.
		if ext == ".ll" || ext == ".s" {
			os.Remove(objFile)
		}

	}
}

// Run a list of objects through a linker and build
// into a single outfile with the given target
func (l *Linker) Run() {
	linker := "clang"
	linkArgs := make([]string, 0)

	if l.target == BinaryTarget {
		linkArgs = append(linkArgs, "-lm", "-lc")
	}

	filename := l.output

	if l.optimize {
		linkArgs = append(linkArgs, "-O3")
	} else {
		linkArgs = append(linkArgs, "-Oz")
	}

	if l.target == ASMTarget {
		log.Timed("Assembly Generation", func() {
			// We want to only write intel syntax. AT&T Sucks
			linkArgs = append(linkArgs, "-S", "-masm=intel", "-Wno-everything")
			// Compile each of the objects to a .s file.
			for _, obj := range l.objectPaths {
				// We only want to leave user generated files in the filesystem
				if strings.HasSuffix(obj, ".ll") {
					// ext := path.Ext(obj)
					out := strings.Replace(obj, ".ll", ".s", -1)
					// Pull the extension of the object file
					if l.dump {
						out = "/dev/stdout"
					}
					// Replace it with .s
					// set the output to that of the .s file
					asmArgs := append(linkArgs, "-o", out, obj)
					// run the compile to asm
					c, _ := util.RunCommandStr(linker, asmArgs...)
					if l.dump {
						fmt.Println(c)
					}
				}

			}
		})

		return
	}

	if l.dump {
		for _, obj := range l.objectPaths {
			if strings.HasSuffix(obj, ".ll") {
				bs, _ := ioutil.ReadFile(obj)
				fmt.Println(string(bs))
			}
		}
	}

	// Append input files to the end of the command
	linkArgs = append(linkArgs, l.objectPaths...)

	// set the output filename
	linkArgs = append(linkArgs, "-o", filename)

	out, err := util.RunCommand(linker, linkArgs...)

	if err != nil {
		log.Fatal("failed to run command `%s %s`: `%s`\n\n%s",
			linker, strings.Join(linkArgs, " "),
			err.Error(), string(out))
	}
}
