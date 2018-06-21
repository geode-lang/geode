package gen

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"gitlab.com/nickwanninger/geode/pkg/parser"
)

// Parser -
type Parser struct {
	name               string              // the filename of the program
	tokens             <-chan parser.Token // channel of tokens from the lexer
	token              parser.Token        // current token, most recently recieved
	nextToken          parser.Token        // next token in the list (allows lookahead)
	topLevelNodes      chan Node
	binaryOpPrecedence map[string]int // maps binary operators to the precidence determining the order of operations
}

// Parse creates and runs a new parser, that returns the
// chan that the nodes will be passed through with
func Parse(tokens <-chan parser.Token) <-chan Node {
	p := &Parser{
		tokens:        tokens,
		topLevelNodes: make(chan Node, 100),
		binaryOpPrecedence: map[string]int{
			"=": 2,
			"<": 10,
			"+": 20,
			"-": 20,
			"*": 40,
			"/": 40,
		},
	}

	// prime the next token for use by reading from the token channel (easier than handling in .next())
	p.nextToken = <-p.tokens
	go p.parse()

	return p.topLevelNodes
}

func (p *Parser) parse() {
	for p.next(); p.token.Type > 0; {
		topLevelNode := p.parseTopLevelStmt()
		if topLevelNode != nil {
			p.topLevelNodes <- topLevelNode
		} else {
			break
		}
	}
	close(p.topLevelNodes)
}

func (p *Parser) next() parser.Token {
	for {
		// Set token to the value of next token and get a new value for nextToken
		p.token = p.nextToken
		p.nextToken = <-p.tokens
		// if that token is valid, break from the loop
		if !(p.token.Type == parser.TokWhitespace || p.token.Type == parser.TokComment) {
			break
		}
	}

	return p.token
}

func (p *Parser) parseTopLevelStmt() Node {
	switch p.token.Type {
	case parser.TokFuncDefn:
		return p.parseFnDefn()
	}

	Error(p.token, "Invalid syntax in root")

	return nil
}

func (p *Parser) getTokenPrecedence(token string) int {
	return p.binaryOpPrecedence[token]
}

// Parse expression is the important ast parser function.Node
// It can parse any expression like `1`, `1 + 2`, `func()`, `func(1)`
// or any mixture. It's basically the workhorse function

// parse any block statement

// Error is a helper function to make logging easier
func Error(t parser.Token, format string, args ...interface{}) {

	fmt.Fprintf(os.Stderr, "\033[31;1m")
	fmt.Fprintf(os.Stderr, "Token Error\n")

	fmt.Fprintf(os.Stderr, "The token in question's data:\n")
	spew.Dump(t)
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
