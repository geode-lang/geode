package ast

import (
	"github.com/nickwanninger/act/pkg/parser"
)

func (p *Parser) parseIdentifierExpr() Node {
	name := p.token.Value
	p.next()

	// Is the next value a paren? If it isnt it is a normal variable reference
	if !p.token.Is(parser.TokLeftParen) {
		n := variableReferenceNode{}
		n.NodeType = nodeVariableReference
		n.Name = name
		return n
	}

	// it was a paren, so we need to parse it as if it were a function call
	n := functionCallNode{}
	n.Name = name
	n.NodeType = nodeFunctionCall

	for p.next(); p.token.Type != parser.TokRightParen; {
		switch p.token.Type {
		case parser.TokComma:
			p.next()
		default:
			arg := p.parseExpression()
			if arg == nil {
				return nil
			}
			n.Args = append(n.Args, arg)
		}
	}
	p.next()

	return n
}
