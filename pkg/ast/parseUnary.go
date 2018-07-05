package ast

// Parse unary will parse a single side of a binary statement

func (p *Parser) parseUnary() Node {

	ptrOps := map[string]bool{
		"&": true,
		"*": true,
	}

	// _, isBinaryOp := p.binaryOpPrecedence[p.token.Value]
	_, isPtrOp := ptrOps[p.token.Value]
	if !isPtrOp {
		return p.parsePrimary()
	}

	if p.token.Value == "&" {

	}

	op := p.token.Value
	p.next()
	operand := p.parseUnary()
	if operand != nil {
		n := UnaryNode{}
		n.NodeType = nodeUnary
		n.Operator = op
		n.Operand = operand
		return n
	}
	return nil
}
