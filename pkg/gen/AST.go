package gen

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"gitlab.com/nickwanninger/geode/pkg/lexer"
)

// Parser -
type Parser struct {
	name               string             // the filename of the program
	tokens             []lexer.Token      // channel of tokens from the lexer
	tokenChan          <-chan lexer.Token // channel of tokens from the lexer
	tokenIndex         int
	token              lexer.Token // current token, most recently recieved
	topLevelNodes      chan Node
	binaryOpPrecedence map[string]int // maps binary operators to the precidence determining the order of operations
}

// Parse creates and runs a new lexer, that returns the
// chan that the nodes will be passed through with
func Parse(tokens <-chan lexer.Token) <-chan Node {
	p := &Parser{
		tokens:        make([]lexer.Token, 0),
		topLevelNodes: make(chan Node, 100),
		binaryOpPrecedence: map[string]int{
			"=":  2,
			"!=": 2,
			"<":  10,
			"+":  20,
			"-":  20,
			"*":  40,
			"/":  40,
		},
	}

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
	p.Error("Required token '%s' is missing. Has '%s' instead.", t.String(), p.token.Type.String())
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

func (p *Parser) peek(o int) lexer.Token {
	target := p.tokenIndex + o
	if target < 0 || target > len(p.tokens)-1 {
		return lexer.Token{}
	}
	return p.tokens[target]
}

func (p *Parser) checkSemiColon() {
	if !p.token.Is(lexer.TokSemiColon) {
		p.Error("Missing Semicolon")
	}
}

func (p *Parser) parseTopLevelStmt() Node {
	switch p.token.Type {
	case lexer.TokFuncDefn:
		return p.parseFnDefn()
	}

	p.Error("Invalid syntax in root")

	return nil
}

func (p *Parser) getTokenPrecedence(token string) int {
	return p.binaryOpPrecedence[token]
}

// Parse expression is the important ast lexer function.Node
// It can parse any expression like `1`, `1 + 2`, `func()`, `func(1)`
// or any mixture. It's basically the workhorse function

// parse any block statement

// Error is a helper function to make logging easier
func (p *Parser) Error(format string, args ...interface{}) {

	t := p.token
	fmt.Fprintf(os.Stderr, "\033[31;1m")
	fmt.Fprintf(os.Stderr, "Token Error\n")
	t.SyntaxError()
	// spew.Dump(t)
	fmt.Fprintf(os.Stderr, format, args...)
	// spew.Fdump(os.Stderr, t)
	fmt.Fprintf(os.Stderr, "\033[0m\n")

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
