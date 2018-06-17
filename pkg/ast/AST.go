package ast

import (
	"encoding/json"
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
		if topLevelNode != nil {
			p.topLevelNodes <- topLevelNode
		} else {
			break
		}
	}
	close(p.topLevelNodes)
}

func (p *Parser) next() parser.Token {
	for p.token = <-p.tokens; p.token.Type == parser.TokWhitespace || p.token.Type == parser.TokComment; p.token = <-p.tokens {
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

// parseFnDefn parses top level function definitions.
func (p *Parser) parseFnDefn() functionNode {
	p.next()

	fn := functionNode{}
	fn.nodeType = nodeFunction

	fn.Name = p.token.Value

	p.next()

	if p.token.Type == parser.TokLeftParen {

		for {
			p.next()

			// If there is an arg
			if p.token.Is(parser.TokType) {

				// Create the node
				v := variableNode{}
				v.nodeType = nodeVariable
				// set it's type
				v.Type = p.token.Value
				// Check the next value
				p.next()

				// if it is an identifier
				if p.token.Is(parser.TokIdent) {

					// Fill the variable's name
					v.Name = p.token.Value
					// and append it to the list
					fn.Args = append(fn.Args, v)
				} else {
					Error(p.token, "Syntax error, invalid parameters to function %s\n", fn.Name)
				}
			}
			// p.next()
			// Break out case (not a comma, or a right paren)
			if p.token.Is(parser.TokRightParen) {
				break
			}
			if p.token.Is(parser.TokComma) {
				continue
			}
		}
	}

	p.next()
	if p.token.Is(parser.TokType) {
		fn.ReturnType = p.token.Value
		// move the token pointer along (no type, so we check the left curly brace)
		p.next()
	} else {
		fn.ReturnType = "void"
	}

	// Get the token after the act arrow (->)
	if p.token.Is(parser.TokLeftCurly) {
		fn.Body = p.parseBlockStmt()
	}
	return fn
}

// func (p *Parser) parsePrimary() Node {

// }

func (p *Parser) parseIfStmt() ifNode {
	// initialize the if statement
	i := ifNode{}

	// // Parse the predicate for the if-statement
	// p.next()
	// ifE := p.parseExpression()

	// if ifE == nil {
	// 	Error(p.token, "Syntax Error, expected")
	// }

	return i
}

// parse while statement
func (p *Parser) parseWhileStmt() whileNode {
	whl := whileNode{}

	return whl
}

// parse any block statement
func (p *Parser) parseBlockStmt() blockNode {
	blk := blockNode{}
	blk.nodeType = nodeBlock

	for {
		p.next()

		if p.token.Is(parser.TokReturn) {
			// fmt.Println("RETURN")
		}

		if p.token.Is(parser.TokType) {
			// parse possible variable definition
		}

		if p.token.Is(parser.TokIf) {
			// Parse if statement
		}

		if p.token.Is(parser.TokFor) {
			// Parse for statement
		}

		if p.token.Is(parser.TokIdent) {
			// blk.nodes = append(blk.nodes, p.parseIdentifierExpr())
		}

		if p.token.Is(parser.TokWhile) {
			blk.Nodes = append(blk.Nodes, p.parseWhileStmt())
		}

		// If the block is over.
		if p.token.Is(parser.TokRightCurly) {
			// Set the cursor on the next character
			p.next()
			break
		}
	}

	return blk
}

// Error is a helper function to make logging easier
func Error(t parser.Token, format string, args ...interface{}) {

	fmt.Fprintf(os.Stderr, "\033[31;1m")
	fmt.Fprintf(os.Stderr, "Token Error\n")

	fmt.Fprintf(os.Stderr, "The token in question's data:\n")

	fmt.Fprintf(os.Stderr, format, args...)
	spew.Fdump(os.Stderr, t)
	fmt.Fprintf(os.Stderr, "\033[0m\n")

	os.Exit(1)
}

// DumpTree takes a channel of nodes and prints all Nodes it recieves,
// then pushes them back out a new channel it makes and returns
func DumpTree(in <-chan Node, useJSON bool) <-chan Node {
	out := make(chan Node)
	go func() {
		for {

			// Read from the input channel of nodes.
			n, stillOpen := <-in

			// If the channel is closed, exit out of the printing phase
			if !stillOpen {
				close(out)
				return
			}

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
