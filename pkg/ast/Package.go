package ast

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/geode-lang/geode/pkg/lexer"
	"github.com/geode-lang/geode/pkg/util"
	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/geode-lang/llvm/ir"
)

// RuntimePackage is the global runtime package
// var RuntimePackage *Package
var dependencyMap map[string]*Package

func init() {
	// RuntimePackage = GetRuntime()
	dependencyMap = make(map[string]*Package)
}

// Package is a wrapper around a module. It is able
// to compile and emit code, as well as lex and parse it.
type Package struct {
	fmt.Stringer

	Name          string
	Source        *lexer.Sourcefile
	Nodes         []Node
	Dependencies  []*Package
	TargetTripple string
	Scope         *Scope
	Compiler      *Compiler
	IsRuntime     bool
	Compiled      bool
	CLinkages     []string
	NamespaceName string
	Namespaces    map[string]*[]Node
}

// NewPackage returns a pointer to a new package
func NewPackage(name string, source *lexer.Sourcefile, scope *Scope, namespaces map[string]*[]Node) *Package {
	p := &Package{}

	p.Name = name
	p.Source = source
	p.Nodes = make([]Node, 0)
	p.Scope = scope
	p.Scope.InjectPrimitives()
	p.Namespaces = namespaces
	return p
}

// String will get the LLVM IR from the package's compiler
func (p *Package) String() string {
	// ir := ""
	ir := &bytes.Buffer{}
	// We need to build up the IR that will be emitted
	// so we can track this information later on.
	fmt.Fprintf(ir, "; ModuleID = %q\n", p.Name)
	fmt.Fprintf(ir, "target datalayout = %q\n", "e-m:o-i64:64-f80:128-n8:16:32:64-S128")
	fmt.Fprintf(ir, "target triple = %q\n", p.TargetTripple)

	// Append the module information
	fmt.Fprintf(ir, "\n%s", p.Compiler.Module.String())

	return ir.String()
}

// Emit will emit the package as IR to a file then build it into an object file for further usage.
// This function returns the path to the object file
func (p *Package) Emit(buildDir string) string {
	outPathBase, _ := filepath.Abs(p.Source.Path)

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

// Hash returns the truncated sha1 of the soruce file
func (p *Package) Hash() []byte {
	return p.Source.Hash()
}

// AddDepPackage appends a dependency from a pacakge
func (p *Package) AddDepPackage(pkg *Package) {
	// Here I check for circular dependencies, which are not allowed
	sourceHash := p.Source.HashName()
	for _, dep := range pkg.Dependencies {
		if dep.Source.HashName() == sourceHash {
			log.Fatal("Circular dependency detected: %s <-> %s\n", pkg.Name, p.Name)
		}
	}
	p.Dependencies = append(p.Dependencies, pkg)
}

// HasDependency returns if a package has access to a namespace
func (p *Package) HasDependency(ns string) bool {
	for _, dep := range p.Dependencies {
		if dep.Name == ns {
			return true
		}
	}
	return false
}

// AddClinkage - takes an absolute path to a c file, and adds it to the link list
func (p *Package) AddClinkage(libPath string) {
	p.CLinkages = append(p.CLinkages, libPath)
}

// LoadDep appends a dependency from a path
func (p *Package) LoadDep(depPath string) *Package {
	filename := path.Base(depPath)
	isStdlib := false

	if strings.HasPrefix(filename, "std:") {
		isStdlib = true
		filename = strings.Replace(filename, "std:", "", -1)
		// Join up the new filename to the standard library source location
		depPath = util.StdLibFile(filename)
	}

	depSource, err := lexer.NewSourcefile(filename)

	if err != nil {
		log.Fatal("Error creating dependency source structure\n")
	}

	depSource.ResolveFile(depPath)

	if isStdlib {
		depSource.Name = path.Join(".stdlib", filename)
	}

	depSource.Preprocess()

	pkgName := fmt.Sprintf("%s", filename)

	if pkg, ok := dependencyMap[depSource.HashName()]; ok {
		p.AddDepPackage(pkg)
		return pkg
	}

	depPkg := NewPackage(pkgName, depSource, p.Scope, p.Namespaces)
	for _ = range depPkg.Parse() {
	}
	dependencyMap[depPkg.Source.HashName()] = depPkg
	p.AddDepPackage(depPkg)
	return depPkg
}

// InjectExternalFunction injects the function without the body, just the sig
func (p *Package) InjectExternalFunction(from *Package, fn *ir.Function, node FunctionNode) {
	// ex := ir.NewFunction(fn.Name, fn.Sig.Ret, fn.Params()...)
	// ex.Sig.Variadic = fn.Sig.Variadic
	scopeItem := NewFunctionScopeItem(fn.Name, node, fn, PublicVisibility)
	p.Scope.Add(scopeItem)
}

// Inject another Package's defintions into this Package
// This is how external dependencies work
func (p *Package) Inject(c *Package) {
	p.Dependencies = append(p.Dependencies, c)
	// Copy over all Scope Variables
	for _, v := range c.Scope.Vals {
		if v.Visibility() == PublicVisibility {

			if v.Type() == ScopeItemFunctionType {
				p.InjectExternalFunction(c, v.Value().(*ir.Function), v.Node().(FunctionNode))
			} else {
				p.Scope.Add(v)
			}

		}
	}
}

// Parse returns a channel of new packages that will be compiled.
func (p *Package) Parse() chan *Package {

	chn := make(chan *Package)
	go func() {
		// There was no cache file
		tokens := lexer.Lex(p.Source)
		nodes := Parse(tokens)

		log.Debug("Parsing package %s\n", p.Name)
		// And append all those nodes to the package's nodes.
		for node := range nodes {
			p.Nodes = append(p.Nodes, node)
		}

		firstNode := p.Nodes[0]
		// The first node *should* always be a namespace node
		if firstNode.Kind() == nodeNamespace {
			p.NamespaceName = firstNode.(NamespaceNode).Name
		} else {
			firstNode.SyntaxError()
			log.Fatal("%q missing namespace. It must be the first statement.\n", p.Source.Path)
		}

		p.Namespaces[p.NamespaceName] = &p.Nodes

		chn <- p
		close(chn)
	}()
	return chn
}

// GetRuntime builds a runtime
func GetRuntime(scope *Scope, namespaces map[string]*[]Node) *Package {
	rts, err := lexer.NewSourcefile("runtime")
	if err != nil {
		log.Fatal("Error creating runtime source structure\n")
	}
	rts.ResolveFile(util.StdLibFile("/runtime"))
	rt := NewPackage("runtime", rts, scope, namespaces)
	rt.IsRuntime = true
	for _ = range rt.Parse() {
	}

	return rt
}

// Compile returns a codegen-ed compiler instance
func (p *Package) Compile(module *ir.Module, targetTripple string) chan *Package {
	p.TargetTripple = targetTripple
	packages := make(chan *Package)

	// fmt.Println(module)

	go func() {
		p.Compiler = NewCompiler(module, p.Name, p, p.Namespaces)
		log.Debug("Compiling Package %s\n", p.Name)

		// Go through all nodes and handle the ones that are dependencies
		for _, node := range p.Nodes {
			if node.Kind() == nodeDependency {
				node.(DependencyNode).Handle(p.Compiler)
			}
		}

		for _, dep := range p.Dependencies {
			if !dep.Compiled {
				dep.Compiled = true
				for pkg := range dep.Compile(module, targetTripple) {
					packages <- pkg
				}
			}
			p.Inject(dep)
		}
		p.Compiled = true

		classes := make([]ClassNode, 0)

		// go through and declare all classes/types
		for _, node := range p.Nodes {
			if node.Kind() == nodeClass {
				classes = append(classes, node.(ClassNode))
			}
		}

		// Declare all classes
		for _, c := range classes {
			c.Declare(p.Compiler.Scope, p.Compiler)
		}

		// Generate the code for all classes
		for _, c := range classes {
			c.Codegen(p.Compiler.Scope, p.Compiler)
		}

		allClassesCorrect := true

		// Verify all classes are correct
		for _, c := range classes {
			correct := c.VerifyCorrectness(p.Compiler.Scope, p.Compiler)
			if !correct {
				allClassesCorrect = false
			}
		}

		if !allClassesCorrect {
			log.Fatal("Exited at class check.\n")
		}

		// go through and declare all the functions
		for _, node := range p.Nodes {
			if node.Kind() == nodeFunction {
				fnNode := node.(FunctionNode)
				if len(fnNode.Generics) == 0 {
					fnNode.Declare(p.Scope.SpawnChild(), p.Compiler)
				}
			}
			// node.Codegen(p.Compiler.Scope.SpawnChild(), p.Compiler)
		}

		for _, node := range p.Nodes {
			if node.Kind() == nodeClass {
				continue
			}
			node.Codegen(p.Compiler.Scope.SpawnChild(), p.Compiler)

		}

		packages <- p
		close(packages)
	}()

	return packages
}
