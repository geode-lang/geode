package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
)

// Parse unary will parse a single side of a binary statement

func (p *Parser) parseUnary() Node {
	startTok := p.token
	validUnaryOps := map[string]bool{
		"&": true,
		"*": true,
		"-": true,
	}

	// parse the "as"
	if p.token.Is(lexer.TokAs) {
	}

	// _, isBinaryOp := p.binaryOpPrecedence[p.token.Value]
	_, isPtrOp := validUnaryOps[p.token.Value]
	if !isPtrOp {
		return p.parsePrimary()
	}

	unaryOp := p.token.Value

	p.next()
	operand := p.parseUnary()
	if unaryOp == "&" {
		if operand.Kind() == nodeVariable {
			// Update operand's RefType if it is a nodeVariable
			n := (operand).(VariableNode)
			n.TokenReference.Token = startTok
			n.RefType = ReferenceAccessStackAddress
			operand = n
		}
	}
	if operand != nil {
		n := UnaryNode{}
		n.TokenReference.Token = startTok
		n.NodeType = nodeUnary
		n.Operator = unaryOp
		n.Operand = operand

		return n
	}

	return nil
}
