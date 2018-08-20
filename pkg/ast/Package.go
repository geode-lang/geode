package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
)

// RuntimePackage is the global runtime package
// var RuntimePackage *Package
var dependencyMap map[string]*Package

func init() {
	// RuntimePackage = GetRuntime()
	dependencyMap = make(map[string]*Package)
}

// Package packages up information around a set of nodes
type Package struct {
	// fmt.Stringer

	Name            string
	Files           map[string]*lexer.Sourcefile
	Nodes           []Node
	Program         *Program
	DependencyPaths []string
}

// NewPackage returns a pointer to a new package
func NewPackage(name string, prog *Program) *Package {
	p := &Package{}

	p.Name = name
	p.Program = prog
	p.Nodes = make([]Node, 0)
	p.Files = make(map[string]*lexer.Sourcefile)
	p.DependencyPaths = make([]string, 0)
	return p
}

// MergeWith merges two packages into a single one, modifying the original
func (p *Package) MergeWith(other *Package) {
	p.Nodes = append(p.Nodes, other.Nodes...)
	for name, src := range other.Files {
		p.Files[name] = src
	}
}

// HasAccessToPackage -
func (p *Package) HasAccessToPackage(name string) bool {

	// return true
	// Base case
	if name == p.Name {
		return true
	}

	for path, pkg := range p.Program.Packages {
		for _, dpath := range p.DependencyPaths {
			if ReduceToDir(path) == ReduceToDir(dpath) && pkg.Name == name {
				return true
			}

		}
	}

	return false

	// for _,
}
