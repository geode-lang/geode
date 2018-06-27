package ast

import (
	"github.com/nickwanninger/geode/pkg/lexer"
	"github.com/nickwanninger/geode/pkg/util/log"
)

// Module is a grouping of top level nodes in a sourcefile and it's scope
type Module struct {
	Name         string
	Parent       *Module
	Dependencies []*Module
	Nodes        []Node
	Scope        *Scope
	Lexer        *lexer.LexState
	Tokens       []lexer.Token
	source       *lexer.Sourcefile
	IsRuntime    bool
	Compiler     *Compiler
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

		for _, node := range m.Nodes {
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
	rts.LoadString(RuntimeSource)
	mod := NewModule("runtime", rts)
	mod.IsRuntime = true
	for _ = range mod.Parse() {

	}

	// mod.Compile()

	m.AddDep(mod)
	// fmt.Println(mod.Compiler.Functions)
	// m.Inject(mod)
}

// Inject another module's defintions into this module
// This is how external dependencies work
func (m *Module) Inject(c *Module) {

	m.Dependencies = append(m.Dependencies, c)
	for k, v := range c.Compiler.Scope.Vals {
		m.Compiler.Scope.Set(k, v)
	}

	for _, fn := range c.Compiler.Functions {
		m.Compiler.AddExternalFunction(fn)
	}
}

// NewModule constructs a module from a source file and returns a pointer to it
func NewModule(name string, src *lexer.Sourcefile) *Module {
	m := &Module{}
	m.Name = name
	m.Nodes = make([]Node, 0)
	m.source = src
	m.Lexer = lexer.NewLexer()
	m.Tokens = make([]lexer.Token, 0) // construct the token array.
	return m
}

// RuntimeSource is the source the runtime will use when compiling
const RuntimeSource string = `
# Define the printf sig
func printf(string format, ...) ...

func exp(int x, int n) int {
	if n = 0 {
		return 1;
	}
	return x * exp(x, n - 1);
}
`
