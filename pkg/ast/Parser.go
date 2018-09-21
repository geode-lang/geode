package ast

import (
	"fmt"

	"github.com/geode-lang/geode/pkg/info"

	"github.com/geode-lang/geode/pkg/lexer"
)

var parserid = 0

// ParseContext is a wrapper around information that allows the parser to understand the world
// around it. This will contain the program that is currently running, etc.
type ParseContext struct {
	Program    *Program
	ClassNames map[string]lexer.Token
}

// Parser -
type Parser struct {
	tokens             []lexer.Token // channel of tokens from the lexer
	tokenIndex         int
	token              lexer.Token // current token, most recently recieved
	topLevelNodes      []Node
	binaryOpPrecedence map[string]int // maps binary operators to the precidence determining the order of operations
	context            *ParseContext
	isFork             bool
	forkParent         *Parser
	ID                 int
}

// NewQuickParser is used to lex and build a parser from tokens quickly
// for small lexing tasks
func NewQuickParser(source string) *Parser {
	p := NewParser()
	p.tokens = lexer.QuickLex(source)
	p.move(0)
	return p
}

// NewParser constructs a new parser and returns a pointer to it
func NewParser() *Parser {
	p := &Parser{
		tokens:             make([]lexer.Token, 0),
		topLevelNodes:      make([]Node, 0),
		binaryOpPrecedence: parserOpPrec,
		ID:                 parserid,
	}
	parserid++

	return p
}

// Fork forks the parser into a child fork that has
// a reference back to the parent
func (p *Parser) Fork() *Parser {
	n := NewParser()
	n.forkParent = p
	n.isFork = true
	n.binaryOpPrecedence = p.binaryOpPrecedence
	n.tokenIndex = p.tokenIndex
	n.token = p.token
	n.tokens = p.tokens
	n.token = p.token
	return n
}

// Join up to a forked parser
func (p *Parser) Join(fork *Parser) error {
	if fork.isFork {
		p.token = fork.token
		p.tokens = fork.tokens
		p.tokenIndex = fork.tokenIndex
		return nil
	}
	return fmt.Errorf("parser join failed because joinee (parser %d) is not a fork", fork.ID)
}

func (p *Parser) reset() {
	p.tokenIndex = 0
	p.move(0)
}

var parserOpPrec = map[string]int{
	"=":  0,
	"+=": 0,
	"-=": 0,
	"*=": 0,
	"/=": 0,
	"||": 1,
	"&&": 1,
	"^":  1,
	"==": 2,
	"!=": 2,
	"<":  10,
	"<=": 10,
	">":  10,
	">=": 10,
	">>": 15,
	"<<": 15,
	"+":  20,
	"-":  20,
	"*":  40,
	"/":  40,
	"%":  40,
}

// Parse creates and runs a new lexer, that returns the
// chan that the nodes will be passed through with
func Parse(tokens []lexer.Token) []Node {
	p := NewParser()

	// prime the next token for use by reading from the token channel (easier than handling in .next())
	for _, t := range tokens {
		if t.Type != lexer.TokWhitespace && t.Type != lexer.TokComment {
			p.tokens = append(p.tokens, t)
		}
	}

	p.move(0)
	p.parse()
	return p.topLevelNodes
}

// Context returns the context of a parser
func (p *Parser) Context() *ParseContext {
	// If the parser doesn't have a context, make a new one
	if p.context == nil {
		p.context = &ParseContext{
			ClassNames: make(map[string]lexer.Token),
		}
	}
	return p.context
}

func (p *Parser) parse() {
	for p.token.Type > 0 {
		topLevelNode := p.parseTopLevelStmt()
		if topLevelNode != nil {
			p.topLevelNodes = append(p.topLevelNodes, topLevelNode)

			info.AddNode(topLevelNode)
		} else {
			break
		}
	}
}

func (p *Parser) requires(t lexer.TokenType) {
	if p.token.Is(t) {
		return
	}

	p.token.SyntaxError()
	p.Errorf("Required token '%s' is missing. Has '%s' instead.\n", t.String(), p.token.Type.String())
}

// Back walks the parser back one token
func (p *Parser) Back() lexer.Token {
	return p.move(-1)
}

// Next advances the parser's state
func (p *Parser) Next() lexer.Token {
	return p.move(1)
}

func (p *Parser) move(o int) lexer.Token {
	p.tokenIndex += o
	p.token = p.Peek(0)
	return p.token
}

// ParserSaveState is a wrapper around some parser
// data so you can save and restore a parser
type ParserSaveState struct {
	index int
}

// Save returns a state that can be restored from
func (p *Parser) Save() ParserSaveState {
	return ParserSaveState{
		index: p.tokenIndex,
	}
}

// Restore takes a state that was previously saved and
// restores the parser's state
func (p *Parser) Restore(state ParserSaveState) {
	p.tokenIndex = state.index
	p.move(0) // make sure to update the token and whatnot, this is the easiest way.
}

// Peek returns the token at an integer offset from the current index
func (p *Parser) Peek(o int) lexer.Token {
	target := p.tokenIndex + o
	if target < 0 || target > len(p.tokens)-1 {
		return lexer.Token{}
	}
	return p.tokens[target]
}

func (p *Parser) globTerminator() {
	if p.token.Is(lexer.TokSemiColon) {
		p.Next()
	}
}

func (p *Parser) parseTopLevelStmt() Node {

	switch p.token.Type {
	case lexer.TokNamespace:
		return p.parseNamespace()
	case lexer.TokDependency:
		return p.parseDependencyStmt()
	case lexer.TokClassDefn:
		return p.parseClassDefn()
	case lexer.TokFuncDefn:
		return p.parseFunctionNode()
	case lexer.TokType:
		node := p.parseGlobalVariableDecl()
		return node
	}
	p.token.SyntaxError()
	p.Errorf("Invalid syntax in root\n")
	return nil
}

func (p *Parser) getTokenPrecedence(token string) int {
	return p.binaryOpPrecedence[token]
}

// Errorf is a helper function to make logging easier
func (p *Parser) Errorf(format string, a ...interface{}) error {

	return fmt.Errorf("%s\n%s", p.token.SyntaxErrorS(), fmt.Sprintf(format, a...))
}
