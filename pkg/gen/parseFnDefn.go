package gen

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

	if p.token.Is(parser.TokRightArrow) {
		fn.Body = blockNode{}
		fn.Body.NodeType = nodeBlock
		fn.Body.Nodes = make([]Node, 0)
		p.next()

		implReturnValue := p.parseExpression()
		implReturn := returnNode{}
		implReturn.Value = implReturnValue
		fn.Body.Nodes = []Node{implReturn}
		if p.token.Is(parser.TokSemiColon) {
			p.next()
		} else {
			Error(p.token, "Missing semicolon after implicit return in function %q", fn.Name)
		}
	} else if p.token.Is(parser.TokLeftCurly) {
		fn.Body = p.parseBlockStmt()
	} else if p.token.Is(parser.TokElipsis) {
		fn.IsExternal = true
		p.next()
	}
	return fn
}
