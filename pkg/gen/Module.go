package gen

import (
	"gitlab.com/nickwanninger/geode/pkg/lexer"
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
func (m *Module) Compile() *Compiler {
	c := NewCompiler(m.Name)
	for _, node := range m.Nodes {
		node.Codegen(c.Scope.SpawnChild(), c)
	}
	return c
}

func (m *Module) String() string {
	return m.Name
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
# The exponent operator function
func exp(int x, int n) int {
	if n = 0 {
		return 1;
	}
	return x * exp(x, n - 1);
}
`
