package ast

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"path/filepath"

	"github.com/geode-lang/geode/pkg/lexer"
	"github.com/geode-lang/geode/pkg/util"
	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/geode-lang/llvm/ir"
)

// Program is a wrapper for information used
// in codegen and dependency resolution
type Program struct {
	Scope         *Scope
	Compiler      *Compiler
	Module        *ir.Module
	ParsedFiles   []string
	Packages      map[string]*Package
	Package       *Package // the currently active package
	CLinkages     []string
	Entry         string
	TargetTripple string
}

// NewProgram creates a program and returns a pointer to it
func NewProgram() *Program {
	p := &Program{}
	p.Scope = NewScope()
	p.Scope.InjectPrimitives()
	p.Compiler = &Compiler{}
	p.Module = ir.NewModule()
	p.Packages = make(map[string]*Package)
	return p
}

// ParsePath parses from some some path and handles
// everything required to get a final compiled program from some
// basic source location
func (p *Program) ParsePath(dir string) {

	// Determine if the path is a directory or not.
	// if isDir, _ := PathIsDir(dir); !isDir {
	// 	// The path isn't a directory, so we just pull the base of the file
	// 	dir = filepath.Dir(dir)
	// }

	dir = ReduceToDir(dir)

	absEntry, err := filepath.Abs(dir)

	if err != nil {
		log.Fatal("Error with parsing entry location\n")
	}

	files, err := p.ParseDir(absEntry)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error parsing folder for geode source files\n")
	}

	for _, file := range files {
		p.ParseFile(file)
	}

	// for
}

// CanParse helps decide whether or not to parse a file based on previously parsed files
func (p *Program) CanParse(file string) bool {
	for _, parsed := range p.ParsedFiles {
		if parsed == file {
			return false
		}
	}
	return true
}

// ParseDir parses a directory for all package information
func (p *Program) ParseDir(path string) ([]string, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	list, err := fd.Readdir(-1)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0, len(list))

	// pkgs = make(map[string]*ast.Package)
	for _, file := range list {
		if strings.HasSuffix(file.Name(), ".g") {
			filename := filepath.Join(path, file.Name())
			if p.CanParse(filename) {
				files = append(files, filename)
			}
		}
	}

	return files, nil

}

// ParseFile will parse the contents of the file at some path into a Package
func (p *Program) ParseFile(path string) {
	p.ParsedFiles = append(p.ParsedFiles, path)
	src, err := lexer.NewSourcefile(path)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error creating Sourcefile context for file at %q\n", path)
	}
	src.LoadFile(path)

	tokens := lexer.Lex(src)

	nodes := make([]Node, 0)

	for node := range Parse(tokens) {
		nodes = append(nodes, node)
	}

	name, err := NamespaceFromNodes(nodes)
	if err != nil {
		log.Fatal("Unable to decide on namespace for file %q", filepath.Clean(path))
	}

	newPkg := NewPackage(name, p)
	newPkg.Program = p
	newPkg.Files[path] = src
	newPkg.Nodes = nodes

	_, found := p.Packages[path]
	if !found {
		p.Packages[path] = newPkg
	}

	for node := range FilterNodes(newPkg.Nodes, nodeDependency) {
		base := filepath.Dir(path)
		dep := node.(DependencyNode)
		for _, depPath := range dep.Paths {
			if dep.CLinkage {
				p.CLinkages = append(p.CLinkages, ResolveDepPath(base, depPath))
			} else {
				newPkg.DependencyPaths = append(newPkg.DependencyPaths, ReduceToDir(ResolveDepPath(base, depPath)))
				p.ParseDep(base, depPath)

			}
		}

	}

}

// ParseDep will parse any dependency relative to the current base
func (p *Program) ParseDep(base, path string) {
	depPath := ResolveDepPath(base, path)
	if p.CanParse(depPath) {
		p.ParsePath(depPath)
	}
}

// ReduceToDir takes a path and reduces it down into its directory
func ReduceToDir(path string) string {
	if isDir, err := PathIsDir(path); !isDir || err != nil {
		path = filepath.Dir(path)
	}
	return path
}

// Codegen sets the programs module to one with nodes filled out
func (p *Program) Codegen() *ir.Module {
	p.Module = ir.NewModule()

	nodes := make([]*PackagedNode, 0)

	for _, pkg := range p.Packages {
		for _, node := range pkg.Nodes {
			nodes = append(nodes, PackageNode(node, pkg, p))
		}
	}

	p.Compiler = NewCompiler(p)

	for node := range FilterPackagedNodes(nodes, nodeClass) {
		node.SetupContext()
		node.Node.(ClassNode).Declare(p)
	}

	for node := range FilterPackagedNodes(nodes, nodeFunction) {
		node.SetupContext()
		node.Node.(FunctionNode).Declare(p)
	}

	for _, node := range nodes {
		node.SetupContext()
		node.Codegen(p)
	}
	// Sort the assorted items in a module because we want to have reproducable
	// hashes in the produced code. As a sideeffect of using a hashmap for path->pkg
	// mapping, it will be out of order most of the time.
	sort.SliceStable(p.Module.Funcs, func(i, j int) bool {
		return p.Module.Funcs[i].Name < p.Module.Funcs[j].Name
	})

	sort.SliceStable(p.Module.Types, func(i, j int) bool {
		return p.Module.Funcs[i].Name < p.Module.Funcs[j].Name
	})
	return p.Module
}

// Emit will emit the package as IR to a file then build it into an object file for further usage.
// This function returns the path to the object file
func (p *Program) Emit(buildDir string) string {
	outPathBase, _ := filepath.Abs(p.Entry)

	outPathBase = path.Join(buildDir, outPathBase)
	extension := filepath.Ext(outPathBase)
	outPathBase = outPathBase[0 : len(outPathBase)-len(extension)]

	baseDir := filepath.Dir(outPathBase)

	os.MkdirAll(baseDir, os.ModePerm)

	llvmFileName := fmt.Sprintf("%s.ll", outPathBase)

	// objFileName := fmt.Sprintf("%s.o", outPathBase)
	ir := p.String()

	writeErr := ioutil.WriteFile(llvmFileName, []byte(ir), 0666)
	if writeErr != nil {
		panic(writeErr)
	}

	// out, err := util.RunCommand("clang", "-flto=thin", "-Wl,-mllvm,-threads=4,-mllvm,-O0", llvmFileName, "-c", "-o", objFileName)
	// if err != nil {
	// 	log.Fatal("%s\n%s\n", string(out), err.Error())
	// }
	return llvmFileName
}

// String will get the LLVM IR from the package's compiler
func (p *Program) String() string {
	// ir := ""
	ir := &bytes.Buffer{}
	// We need to build up the IR that will be emitted
	// so we can track this information later on.
	fmt.Fprintf(ir, "target datalayout = %q\n", "e-m:o-i64:64-f80:128-n8:16:32:64-S128")
	fmt.Fprintf(ir, "target triple = %q\n", p.TargetTripple)
	fmt.Fprintf(ir, "\n")

	// Append the module information
	fmt.Fprintf(ir, "\n%s", p.Compiler.Module.String())

	return ir.String()
}

// ResolveDepPath returns the absolute location to a dependency
func ResolveDepPath(base, filename string) string {
	if strings.HasPrefix(filename, "std:") {
		filename = strings.Replace(filename, "std:", "", -1)
		// Join up the new filename to the standard library source location
		base = util.StdLibFile(filename)
	}
	return filepath.Join(base, filename)
}

// PathIsDir returns if a given path is a directory or not
func PathIsDir(pth string) (bool, error) {
	fd, err := os.Open(pth)
	if err != nil {
		return false, err
	}
	defer fd.Close()
	stat, err := fd.Stat()
	if err != nil {
		return false, err
	}
	return stat.IsDir(), nil
}

// NamespaceFromNodes takes an array of nodes and returns the namespace name of them
func NamespaceFromNodes(nodes []Node) (string, error) {
	for _, n := range nodes {
		if n.Kind() == nodeNamespace {
			return n.(NamespaceNode).Name, nil
		}
	}

	return "error", fmt.Errorf("nodes have no package name")
}

// FilterNodes returns only the nodes that have the type passed in
func FilterNodes(nodes []Node, t NodeType) chan Node {
	filtered := make(chan Node)

	go func() {
		for _, n := range nodes {
			if n.Kind() == t {
				filtered <- n
			}
		}
		close(filtered)
	}()

	return filtered
}

// FilterPackagedNodes returns only the nodes that have the type passed in
func FilterPackagedNodes(nodes []*PackagedNode, t NodeType) chan *PackagedNode {
	filtered := make(chan *PackagedNode)
	go func() {
		for _, n := range nodes {
			if n.Node.Kind() == t {
				filtered <- n
			}
		}
		close(filtered)
	}()

	return filtered
}

// PackagedNode wraps around a certain node and allows better codegen
// in the context of a certain package
type PackagedNode struct {
	Pkg     *Package
	Program *Program
	Node    Node
}

// Codegen will generate the node this PackagedNode wraps
func (p *PackagedNode) Codegen(prog *Program) {
	p.SetupContext()
	p.Node.Codegen(prog)
}

// SetupContext modifies the program to help with context information
func (p *PackagedNode) SetupContext() {
	p.Program.Package = p.Pkg
	p.Program.Scope.PackageName = p.Pkg.Name
}

// PackageNode takes a node, it's package and the program context
// and creates an encapsulated context for it
func PackageNode(node Node, pkg *Package, prog *Program) *PackagedNode {
	n := &PackagedNode{}
	n.Node = node
	n.Pkg = pkg
	n.Program = prog
	return n
}
