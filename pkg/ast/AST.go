package ast

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
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

// Parse creates and runs a new lexer, that returns the
// chan that the nodes will be passed through with
func Parse(tokens chan lexer.Token) <-chan Node {
	// spew.Dump(QuickParseExpression("1 + 1"))
	p := &Parser{
		tokens:        make([]lexer.Token, 0),
		topLevelNodes: make(chan Node),
		binaryOpPrecedence: map[string]int{
			"|":  1,
			"&":  1,
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
	p.Error("Required token '%s' is missing. Has '%s' instead.\n", t.String(), p.token.Type.String())
}

func (p *Parser) back() lexer.Token {
	return p.move(-1)
}
func (p *Parser) next() lexer.Token {
	return p.move(1)
}

func (p *Parser) move(o int) lexer.Token {
	p.tokenIndex += o
	p.token = p.peek(0)
	return p.token
}

// ParserSaveState is a wrapper around some parser
// data so you can save and restore a parser
type ParserSaveState struct {
	index int
}

func (p *Parser) save() ParserSaveState {
	return ParserSaveState{
		index: p.tokenIndex,
	}
}

func (p *Parser) restore(state ParserSaveState) {
	p.tokenIndex = state.index
	p.move(0) // make sure to update the token and whatnot, this is the easiest way.
}

func (p *Parser) peek(o int) lexer.Token {
	target := p.tokenIndex + o
	if target < 0 || target > len(p.tokens)-1 {
		return lexer.Token{}
	}
	return p.tokens[target]
}

func (p *Parser) allowSemiColon() {
	p.requires(lexer.TokSemiColon)
	// fmt.Println(p.token)
	// if p.token.Is(lexer.TokSemiColon) {
	// 	p.next()
	// 	fmt.Println("->", p.token)
	// }
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
		return p.parseFnDefn()
		// case lexer.TokType:
		// 	log.Debug("parseTopLevelStmt - TokFuncDefn\n")
		// 	return p.parseVariableDefn(true)
	}

	p.token.SyntaxError()
	p.Error("Invalid syntax in root\n")
	return nil
}

func (p *Parser) getTokenPrecedence(token string) int {
	return p.binaryOpPrecedence[token]
}

// Error is a helper function to make logging easier
func (p *Parser) Error(format string, args ...interface{}) {
	log.Fatal(format, args...)
	os.Exit(1)
}

// DumpTree takes a channel of nodes and prints all Nodes it recieves,
// then pushes them back out a new channel it makes and returns
func DumpTree(in <-chan Node, useJSON bool) <-chan Node {
	out := make(chan Node)
	go func() {
		for n := range in {

			if useJSON {
				// Attempt to parse the
				j, jsonParseError := json.MarshalIndent(n, "", "    ")
				// We need to warn the user of some error printing the node, but don't fail
				// Instead print it as a spew dump, this way presentation is still given,
				// but maybe not in a json format.
				if jsonParseError != nil {
					fmt.Println("Error printing node: ", jsonParseError)
					fmt.Println("Raw representation of the node:")
					spew.Dump(n)
				}
				fmt.Println(string(j))
			} else {
				spew.Dump(n)
			}

			out <- n
		}
	}()
	return out
}
