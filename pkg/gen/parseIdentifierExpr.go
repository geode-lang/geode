package gen

import (
	"gitlab.com/nickwanninger/geode/pkg/lexer"
)

func (p *Parser) parseIdentifierExpr() Node {
	name := p.token.Value
	p.next()
	if p.token.Is(lexer.TokLeftArrow) {
		n := variableNode{}
		n.Name = name
		n.Reassignment = true
		n.HasValue = true
		p.next()
		n.NodeType = nodeVariable
		n.Body = p.parseExpression()
		return n
	}

	// Is the next value a paren? If it isnt it is a normal variable reference
	if !p.token.Is(lexer.TokLeftParen) {
		n := variableReferenceNode{}
		n.NodeType = nodeVariableReference
		n.Name = name
		return n
	}

	// it was a paren, so we need to parse it as if it were a function call
	n := functionCallNode{}
	n.Name = name
	n.NodeType = nodeFunctionCall

	for p.next(); p.token.Type != lexer.TokRightParen; {
		switch p.token.Type {
		case lexer.TokComma:
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
