package ast

import (
	"os"

	"github.com/geode-lang/geode/pkg/lexer"
	"github.com/geode-lang/geode/pkg/util/log"
)

// Parser -
type Parser struct {
	tokens             []lexer.Token // channel of tokens from the lexer
	tokenIndex         int
	token              lexer.Token // current token, most recently recieved
	topLevelNodes      chan Node
	binaryOpPrecedence map[string]int // maps binary operators to the precidence determining the order of operations
}

// NewQuickParser is used to lex and build a parser from tokens quickly
// for small lexing tasks
func NewQuickParser(source string) *Parser {
	p := &Parser{}
	p.tokens = lexer.QuickLex(source)
	p.move(0)
	return p
}

func (p *Parser) fork() *Parser {
	n := &Parser{}

	n.binaryOpPrecedence = p.binaryOpPrecedence
	n.tokenIndex = p.tokenIndex
	n.token = p.token
	n.tokens = p.tokens
	n.token = p.token

	return n
}

func (p *Parser) reset() {
	p.tokenIndex = 0
	p.move(0)
}

// Parse creates and runs a new lexer, that returns the
// chan that the nodes will be passed through with
func Parse(tokens chan lexer.Token) <-chan Node {
	p := &Parser{
		tokens:        make([]lexer.Token, 0),
		topLevelNodes: make(chan Node),
		binaryOpPrecedence: map[string]int{
			"||": 1,
			"&&": 1,
			"^":  1,
			"=":  2,
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
		},
	}

	// os.Exit(1)

	// prime the next token for use by reading from the token channel (easier than handling in .next())

	for t := range tokens {
		if t.Type != lexer.TokWhitespace && t.Type != lexer.TokComment {
			p.tokens = append(p.tokens, t)
		}
	}

	p.move(0)
	go p.parse()

	return p.topLevelNodes
}

func (p *Parser) parse() {
	for p.token.Type > 0 {
		topLevelNode := p.parseTopLevelStmt()
		if topLevelNode != nil {
			p.topLevelNodes <- topLevelNode
		} else {
			break
		}
	}
	close(p.topLevelNodes)
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
	case lexer.TokIdent:
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
func (p *Parser) Errorf(format string, args ...interface{}) {
	p.token.SyntaxError()
	log.Fatal(format, args...)
	os.Exit(1)
}
