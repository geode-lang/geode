package ast

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/nickwanninger/act/pkg/parser"
	"os"
)

// Parser -
type Parser struct {
	name               string              // the filename of the program
	tokens             <-chan parser.Token // channel of tokens from the lexer
	token              parser.Token        // current token, most recently recieved
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

	go p.parse()

	return p.topLevelNodes
}

func (p *Parser) parse() {
	for p.next(); p.token.Type > 0; {
		topLevelNode := p.parseTopLevelStmt()
		spew.Dump(p.token)
		if topLevelNode != nil {
			p.topLevelNodes <- topLevelNode
		} else {
			break
		}
	}
	close(p.topLevelNodes)
}

func (p *Parser) next() parser.Token {
	for p.token = <-p.tokens; p.token.Type == parser.GetTokenId("WHITESPACE") || p.token.Type == parser.GetTokenId("COMMENT"); p.token = <-p.tokens {
	}
	return p.token
}

func (p *Parser) parseTopLevelStmt() Node {
	switch p.token.Type {
	case parser.GetTokenId("ACT_DEFN"):
		return p.parseFnDefn()
	}

	return nil
}

// parseFnDefn parses top level function definitions.
func (p *Parser) parseFnDefn() functionNode {
	p.next()

	fn := functionNode{}

	fn.name = p.token.Value

	p.next()

	if p.token.Is("LEFT_PAREN") {
		p.next()
		for {

			if p.token.Is("TYPE") {
				v := variableNode{}
				v.typ = p.token.Value
				p.next()
				if p.token.Is("IDENTIFIER") {
					v.name = p.token.Value
					fn.args = append(fn.args, v)
				} else {
					Error(p.token, "Syntax error, invalid parameters to function %s\n", fn.name)
				}
			}
			p.next()
			// Break out case (not a comma, or a right paren)
			if p.token.Is("RIGHT_PAREN") {
				break
			}
			if p.token.Is("COMMA") {
				p.next()
				continue
			}
			Error(p.token, "")
		}
	}

	if p.token.Is("ACTARROW") {
		Error(p.token, "Syntax Error, function %s missing '->'\n", fn.name)
	}

	// Get the token after the act arrow (->)
	p.next()
	if p.token.Is("LEFT_CURLY") {
		fn.body = p.parseBlockStmt()
	}
	return fn
}

func (p *Parser) parseBlockStmt() blockNode {
	blk := blockNode{}
	for {
		p.next()
		if p.token.Is("RIGHT_CURLY") {
			break
		}
	}

	p.next()

	return blk
}

// Error is a helper function to make logging easier
func Error(t parser.Token, format string, args ...interface{}) {

	fmt.Fprintf(os.Stderr, "\033[31;1m")
	fmt.Fprintf(os.Stderr, "Token Error\n")
	fmt.Fprintf(os.Stderr, parser.SyntaxError(t, fmt.Sprintf(format, args...)))
	fmt.Fprintf(os.Stderr, "The token in question's data:\n")
	spew.Fdump(os.Stderr, t)
	fmt.Fprintf(os.Stderr, "\033[0m\n")

	os.Exit(1)
}

// DumpTree takes a channel of nodes and prints all Nodes it recieves,
// then pushes them back out a new channel it makes and returns
func DumpTree(in <-chan Node) <-chan Node {
	out := make(chan Node)
	go func() {
		for {
			n, ok := <-in
			if !ok {
				close(out)
				return
			}
			fmt.Println("")
			fmt.Println("================ Root Node Parsed ================")
			fmt.Println("")
			spew.Dump(n)
			out <- n
		}
	}()
	return out
}
