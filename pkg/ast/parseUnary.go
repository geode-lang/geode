package ast

import "github.com/nickwanninger/geode/pkg/lexer"

// Parse unary will parse a single side of a binary statement

func (p *Parser) parseUnary() Node {

	if p.token.Is(lexer.TokType) {
		return p.parseTypeCast()
	}

	_, isBinaryOp := p.binaryOpPrecedence[p.token.Value]
	if !isBinaryOp {
		return p.parsePrimary()
	}
	name := p.token.Value
	p.next()
	operand := p.parseUnary()
	if operand != nil {
		n := unaryNode{}
		n.NodeType = nodeUnary
		n.Name = name
		n.Operand = operand
		return n
	}
	return nil
}
