package ast

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/llir/llvm/ir"
	"github.com/nickwanninger/geode/pkg/lexer"
	"github.com/nickwanninger/geode/pkg/util/log"
)

// Package is a wrapper around a module. It is able
// to compile and emit code, as well as lex and parse it.
type Package struct {
	fmt.Stringer

	Name               string
	Lexer              *lexer.LexState
	Source             *lexer.Sourcefile
	Nodes              []Node
	Dependencies       []*Package
	Scope              *Scope
	Compiler           *Compiler
	IsRuntime          bool
	objectFilesEmitted []string
}

// NewPackage returns a pointer to a new package
func NewPackage(name string, source *lexer.Sourcefile) *Package {
	p := &Package{}

	p.Name = name
	p.Source = source
	p.Nodes = make([]Node, 0)
	p.Scope = NewScope()
	p.Lexer = lexer.NewLexer()

	return p
}

// String will get the LLVM IR from the package's compiler
func (p *Package) String() string {
	ir := ""
	// We need to build up the IR that will be emitted
	// so we can track this information later on.
	ir += fmt.Sprintf("; ModuleID = %q\n", p.Name)
	ir += fmt.Sprintf("; SourceHash = %x\n", p.Hash())
	ir += fmt.Sprintf("; UnixDate = %d\n", time.Now().Unix())
	ir += fmt.Sprintf("source_filename = %q\n", p.Source.Path)

	ir += "\n"
	// Append the module information
	ir += fmt.Sprintf("%s\n", p.Compiler.Module.String())

	return ir
}

// Emit will emit the package as IR to a file for further compiling
func (p *Package) Emit() string {
	name := strings.Replace(p.Name, ".g", "", -1)
	filename := fmt.Sprintf("%s.%x.ll", name, p.Hash())
	ir := p.String()

	writeErr := ioutil.WriteFile(filename, []byte(ir), 0666)
	if writeErr != nil {
		panic(writeErr)
	}

	p.objectFilesEmitted = append(p.objectFilesEmitted, filename)
	return filename
}

// Hash returns the truncated sha1 of the soruce file
func (p *Package) Hash() []byte {
	return p.Source.Hash()
}

// AddDep appends a dependency
func (p *Package) AddDep(pkg *Package) {
	p.Dependencies = append(p.Dependencies, pkg)
}

// InjectExternalFunction injects the function without the body, just the sig
func (p *Package) InjectExternalFunction(fn *ir.Function) {
	ex := ir.NewFunction(fn.Name, fn.Sig.Ret, fn.Params()...)
	ex.Sig.Variadic = fn.Sig.Variadic
	scopeItem := NewFunctionScopeItem(fn.Name, ex, PublicVisibility)
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

				p.InjectExternalFunction(v.Value().(*ir.Function))
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
		// Pull the source bytes out
		srcBytes := p.Source.Bytes()
		// go and lex the bytes
		go p.Lexer.Lex(srcBytes) // run the lexer
		// Parse the bytes into a channel of nodes
		nodes := Parse(p.Lexer.Tokens)
		// And append all those nodes to the package's nodes.
		for node := range nodes {
			p.Nodes = append(p.Nodes, node)
		}

		chn <- p
		close(chn)
	}()
	return chn
}

// AddRuntime builds the runtime and injects it into the module
func (p *Package) AddRuntime() {
	rts, err := lexer.NewSourcefile("runtime")
	if err != nil {
		log.Fatal("Error creating runtime source structure\n")
	}
	gopath := os.Getenv("GOPATH")
	rts.LoadFile(gopath + "/src/github.com/nickwanninger/geode/lib/lib.g")
	rt := NewPackage("runtime", rts)
	rt.IsRuntime = true
	for _ = range rt.Parse() {

	}

	p.AddDep(rt)
}

// Compile returns a codegen-ed compiler instance
func (p *Package) Compile() chan *Package {
	packages := make(chan *Package)
	go func() {
		p.Compiler = NewCompiler(p.Name, p)
		if !p.IsRuntime {
			p.AddRuntime()
		}

		for _, dep := range p.Dependencies {
			for pkg := range dep.Compile() {
				packages <- pkg
			}
			p.Inject(dep)
		}

		// First we *Need* to go through and declare all the functions. This is because
		for _, node := range p.Nodes {
			if node.Kind() == nodeFunction {
				node.(functionNode).Declare(p.Compiler.Scope.SpawnChild(), p.Compiler)
			}
			// node.Codegen(p.Compiler.Scope.SpawnChild(), p.Compiler)
		}

		for _, node := range p.Nodes {
			// node.Codegen(p.Compiler.Scope, p.Compiler)
			node.Codegen(p.Compiler.Scope.SpawnChild(), p.Compiler)
		}

		packages <- p
		close(packages)
	}()

	return packages
}
