package ast

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
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
	buildDir    string
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

// SetBuildDir sets the target output format of the linker
// be it binary, asm, etc...
func (l *Linker) SetBuildDir(t string) {
	l.buildDir = t
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

	// linkArgs = append(linkArgs, "-g")
	if l.target == BinaryTarget {
		linkArgs = append(linkArgs, "-lm", "-lc")
	}

	filename := l.output

	if l.optimize {
		linkArgs = append(linkArgs, "-O3")
	} else {
		linkArgs = append(linkArgs, "-O0")
	}

	if l.dump {

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
						fmt.Println(c)
					}

				}
			})
			return
		}

		for _, obj := range l.objectPaths {
			if strings.HasSuffix(obj, ".ll") {
				bs, _ := ioutil.ReadFile(obj)
				fmt.Println(string(bs))
			}
		}

	}

	for i, obj := range l.objectPaths {
		outbase := path.Join(l.buildDir, obj)

		extension := filepath.Ext(outbase)
		if extension == ".ll" {
			// outbase = obj
		}

		if extension == ".a" {

		}

		if extension == ".c" {
			outbase = outbase[0 : len(outbase)-len(extension)]

			hash := util.HashFile(obj)

			cachefile := outbase + ".cache"
			objFile := outbase + ".o"

			cachedat, err := ioutil.ReadFile(cachefile)
			if err != nil || strings.Compare(string(cachedat), hash) != 0 {
				os.MkdirAll(path.Dir(outbase), os.ModePerm)
				// the file doesnt exist, we need to compile it
				out, err := util.RunCommand("clang", "-c", "-o", objFile, obj)
				if err != nil {
					log.Fatal("(%s) %s\n", err, string(out))
				}
				// fmt.Println(string(out), err)
				ioutil.WriteFile(cachefile, []byte(hash), os.ModePerm)
			}
			l.objectPaths[i] = objFile
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
