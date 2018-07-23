package ast

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/geode-lang/geode/pkg/lexer"
	"github.com/geode-lang/geode/pkg/util"
	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/llir/llvm/ir"
)

// RuntimePackage is the global runtime package
var RuntimePackage *Package
var dependencyMap map[string]*Package

func init() {
	RuntimePackage = GetRuntime()
	dependencyMap = make(map[string]*Package)
}

// Package is a wrapper around a module. It is able
// to compile and emit code, as well as lex and parse it.
type Package struct {
	fmt.Stringer

	Name               string
	Source             *lexer.Sourcefile
	Nodes              []Node
	Dependencies       []*Package
	TargetTripple      string
	Scope              *Scope
	Compiler           *Compiler
	IsRuntime          bool
	objectFilesEmitted []string
	Compiled           bool
	CLinkages          []string
	NamespaceName      string
}

// NewPackage returns a pointer to a new package
func NewPackage(name string, source *lexer.Sourcefile) *Package {
	p := &Package{}

	p.Name = name
	p.Source = source
	p.Nodes = make([]Node, 0)
	p.Scope = NewScope()
	p.Scope.InjectPrimitives()
	return p
}

// String will get the LLVM IR from the package's compiler
func (p *Package) String() string {
	// ir := ""
	ir := &bytes.Buffer{}
	// We need to build up the IR that will be emitted
	// so we can track this information later on.
	fmt.Fprintf(ir, "; ModuleID = %q\n", p.Name)
	fmt.Fprintf(ir, "; UnixDate = %d\n", time.Now().Unix())
	fmt.Fprintf(ir, "target triple = %q\n", p.TargetTripple)

	// Append the module information
	fmt.Fprintf(ir, "\n%s", p.Compiler.Module.String())

	return ir.String()
}

// const buildDir = ".geode_build/"

// Emit will emit the package as IR to a file for further compiling
func (p *Package) Emit(buildDir string) string {
	name := strings.Replace(p.Source.Name, ".g", "", -1)
	filename := fmt.Sprintf("%s.ll", name)
	objFileName := path.Join(util.GetCacheDir(), fmt.Sprintf("%x.o", p.Hash()))

	pwd, _ := os.Getwd()
	filename = strings.Replace(filename, pwd, "", -1)
	ir := p.String()

	buildFolder := path.Join(buildDir, path.Dir(filename))

	filename = path.Join(buildFolder, path.Base(filename))

	os.MkdirAll(buildFolder, os.ModePerm)

	writeErr := ioutil.WriteFile(filename, []byte(ir), 0666)
	if writeErr != nil {
		panic(writeErr)
	}

	_, err := util.RunCommand("clang", filename, "-c", "-o", objFileName)
	if err != nil {
		log.Fatal("%s\n", err.Error())
	}

	// fmt.Println(filename, objFileName)

	p.objectFilesEmitted = append(p.objectFilesEmitted, filename)
	return filename
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

	pkgName := fmt.Sprintf("%s", filename)

	if pkg, ok := dependencyMap[depSource.HashName()]; ok {
		p.AddDepPackage(pkg)
		return pkg
	}

	depPkg := NewPackage(pkgName, depSource)
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
				// fmt.Println(p.Name, v.Name())
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

		cacheFolder := path.Join(util.GetCacheDir())
		cacheFile := fmt.Sprintf("%s/%x.cache", cacheFolder, p.Hash())

		var cacheBuffer bytes.Buffer

		if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
			// There was no cache file
			tokens := lexer.Lex(p.Source)
			nodes := Parse(tokens)

			log.Debug("Parsing package %s\n", p.Name)
			// And append all those nodes to the package's nodes.
			for node := range nodes {
				p.Nodes = append(p.Nodes, node)
			}

			enc := gob.NewEncoder(&cacheBuffer)
			err := enc.Encode(p.Nodes)
			if err != nil {
				log.Fatal("encode error:", err)
			}

			os.MkdirAll(cacheFolder, os.ModePerm)

			err = ioutil.WriteFile(cacheFile, cacheBuffer.Bytes(), 0644)
			if err != nil {
				log.Fatal("Error writing cache file to home folder")
			}
		} else {
			f, _ := os.Open(cacheFile) // Error handling elided for brevity.
			defer f.Close()
			io.Copy(&cacheBuffer, f) // Error handling elided for brevity.

			// There was a cache file.
			dec := gob.NewDecoder(&cacheBuffer)
			p.Nodes = make([]Node, 0)
			err = dec.Decode(&p.Nodes)
			if err != nil {
				log.Fatal("Error decoding node array in package %s\n", p.Name)
			}
		}

		chn <- p
		close(chn)
	}()
	return chn
}

// GetRuntime builds a runtime
func GetRuntime() *Package {
	rts, err := lexer.NewSourcefile("runtime")
	if err != nil {
		log.Fatal("Error creating runtime source structure\n")
	}
	rts.ResolveFile(util.StdLibFile("/runtime"))
	rt := NewPackage("runtime", rts)
	rt.IsRuntime = true
	for _ = range rt.Parse() {
	}

	return rt
}

// Compile returns a codegen-ed compiler instance
func (p *Package) Compile(module *ir.Module, targetTripple string) chan *Package {
	p.TargetTripple = targetTripple
	packages := make(chan *Package)

	go func() {
		p.Compiler = NewCompiler(module, p.Name, p)
		log.Debug("Compiling Package %s\n", p.Name)

		if !p.IsRuntime {
			log.Debug("Injecting runtime into '%s'\n", p.Name)
			// p.LoadDep("std:_runtime.g")
			p.AddDepPackage(RuntimePackage)
		}

		firstNode := p.Nodes[0]
		// The first node *should* always be a namespace node
		if firstNode.Kind() == nodeNamespace {
			p.NamespaceName = firstNode.(NamespaceNode).Name
		} else {
			firstNode.SyntaxError()
			log.Fatal("%q missing namespace. It must be the first statement.\n", p.Source.Path)
		}

		// Go through all nodes and handle the ones that are dependencies
		for _, node := range p.Nodes {
			if node.Kind() == nodeDependency {
				node.(DependencyNode).Handle(p.Compiler)
			}
		}

		for _, dep := range p.Dependencies {
			// fmt.Println(p.NamespaceName, len(p.Nodes))
			if !dep.Compiled {
				dep.Compiled = true
				for pkg := range dep.Compile(module, targetTripple) {
					packages <- pkg
				}
			}
			p.Inject(dep)
		}
		p.Compiled = true

		// go through and generate all classes/types
		for _, node := range p.Nodes {
			if node.Kind() == nodeClass {
				node.(ClassNode).Codegen(p.Compiler.Scope, p.Compiler)
			}
		}

		// go through and declare all the functions
		for _, node := range p.Nodes {
			if node.Kind() == nodeFunction {
				fnNode := node.(FunctionNode)
				if len(fnNode.Generics) == 0 {

				}
				fnNode.Declare(p.Scope.SpawnChild(), p.Compiler)

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
