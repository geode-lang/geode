package ast

import (
	"github.com/nickwanninger/act/pkg/parser"
	"github.com/nickwanninger/act/pkg/types"
)

func (p *Parser) parseFnDefn() functionNode {
	p.next()

	fn := functionNode{}
	fn.NodeType = nodeFunction

	fn.Name = p.token.Value

	p.next()

	if p.token.Type == parser.TokLeftParen {

		for {
			// If there is an arg
			if p.nextToken.Is(parser.TokType) {
				p.next()
				fn.Args = append(fn.Args, p.parseVariableDefn(false))
			}
			p.next()
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
		fn.ReturnType = types.GlobalTypeMap.GetType(p.token.Value)
		// move the token pointer along (no type, so we check the left curly brace)
		p.next()
	} else {
		fn.ReturnType = types.GlobalTypeMap.GetType("void")
	}

	// Get the token after the act arrow (->)
	if p.token.Is(parser.TokLeftCurly) {
		fn.Body = p.parseBlockStmt()
	}
	return fn
}
