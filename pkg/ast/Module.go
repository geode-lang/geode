package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/nickwanninger/geode/pkg/lexer"
	"github.com/nickwanninger/geode/pkg/util/log"
)

// Module is a grouping of top level nodes in a sourcefile and it's scope
type Module struct {
	Name         string
	Parent       *Module
	Dependencies []*Module
	Nodes        []Node
	// Scope        *Scope
	Lexer     *lexer.LexState
	Tokens    []lexer.Token
	source    *lexer.Sourcefile
	IsRuntime bool
	Compiler  *Compiler
}

// Modules have a parse method on them that takes a channel of Modules.
// When the module is done parsing, it will pass itself back through
// the channel as a signal it is done. This is because the parser runs
// goroutines that could take time

// Parse runs the lexer and parser on the source file already passed in
func (m *Module) Parse() chan *Module {

	chn := make(chan *Module)
	go func() {
		srcBytes := m.source.Bytes()
		go m.Lexer.Lex(srcBytes) // run the lexer

		nodes := Parse(m.Lexer.Tokens)

		for node := range nodes {
			m.Nodes = append(m.Nodes, node)
		}

		chn <- m
		close(chn)
	}()
	return chn
}

// Compile returns a codegen-ed compiler instance
func (m *Module) Compile() chan *Compiler {
	compilers := make(chan *Compiler)
	go func() {
		m.Compiler = NewCompiler(m.Name)
		if !m.IsRuntime {
			m.AddRuntime()
		}

		for _, dep := range m.Dependencies {
			for c := range dep.Compile() {
				compilers <- c
			}
			m.Inject(dep)
		}

		// First we *Need* to go through and declare all the functions. This is because
		for _, node := range m.Nodes {
			if node.Kind() == nodeFunction {
				node.(functionNode).Declare(m.Compiler.Scope.SpawnChild(), m.Compiler)
			}
			// node.Codegen(m.Compiler.Scope.SpawnChild(), m.Compiler)
		}

		for _, node := range m.Nodes {
			// node.Codegen(m.Compiler.Scope, m.Compiler)
			node.Codegen(m.Compiler.Scope.SpawnChild(), m.Compiler)
		}

		compilers <- m.Compiler
		close(compilers)
	}()

	return compilers
}

func (m *Module) String() string {
	return m.Name
}

// AddDep appends a dependency
func (m *Module) AddDep(mod *Module) {
	m.Dependencies = append(m.Dependencies, mod)
}

// AddRuntime builds the runtime and injects it into the module
func (m *Module) AddRuntime() {
	rts, err := lexer.NewSourcefile("runtime")
	if err != nil {
		log.Fatal("Error creating runtime source structure\n")
	}
	rts.LoadString(RuntimeGeode)
	mod := NewModule("runtime", rts)
	mod.IsRuntime = true
	for _ = range mod.Parse() {

	}

	m.AddDep(mod)
}

// InjectExternalFunction injects the function without the body, just the sig
func (m *Module) InjectExternalFunction(fn *ir.Function) {
	ex := ir.NewFunction(fn.Name, fn.Sig.Ret, fn.Params()...)
	ex.Sig.Variadic = fn.Sig.Variadic
	// m.Compiler.Module.AppendFunction(ex)
	scopeItem := NewFunctionScopeItem(fn.Name, ex, PublicVisibility)
	m.Compiler.Scope.Add(scopeItem)

}

// Inject another module's defintions into this module
// This is how external dependencies work
func (m *Module) Inject(c *Module) {
	m.Dependencies = append(m.Dependencies, c)
	// Copy over all Scope Variables
	for _, v := range c.Compiler.Scope.Vals {
		if v.Visibility() == PublicVisibility {

			if v.Type() == ScopeItemFunctionType {

				m.InjectExternalFunction(v.Value().(*ir.Function))
			} else {
				m.Compiler.Scope.Add(v)
			}

		}
	}
}

// NewModule constructs a module from a source file and returns a pointer to it
func NewModule(name string, src *lexer.Sourcefile) *Module {
	m := &Module{}
	m.Name = name
	m.Nodes = make([]Node, 0)
	m.source = src
	m.Lexer = lexer.NewLexer()
	// m.Scope = NewScope()
	m.Tokens = make([]lexer.Token, 0) // construct the token array.
	return m
}

// RuntimeGeode is the source the runtime will use when compiling
const RuntimeGeode string = `

# This print function is a wrapper around printf
func print(string format, ...) ...

# Some overloading of the print method for 
# different types. These are mangled
func print(int a) -> print("%d\n", a);
func print(float a) -> print("%f\n", a);


func readfile(string path) byte* ...

func exp(int x, int n) int {
	if n = 0 {
		return 1;
	}
	return x * exp(x, n - 1);
}
`

// RuntimeC is the implementation of some of the functions from
// the runtime. Things that should just be implemented in c.
const RuntimeC string = `
#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>


// Readfile just takes some path
// and returns bytes containing
// the content of the file.
char* readfile(char* path) {
	FILE *f = fopen(path, "rb");
	fseek(f, 0, SEEK_END);
	long fsize = ftell(f);
	fseek(f, 0, SEEK_SET); // same as rewind(f);
	char *string = malloc(fsize + 1);
	fread(string, fsize, 1, f);
	fclose(f);
	string[fsize] = 0;
	return string;
}

// the print function wrapper.
void print(char *fmt, ...) {
	va_list args;
	va_start(args, fmt);
	vprintf(fmt, args);
	va_end(args);
}
`
