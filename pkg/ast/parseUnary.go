package ast

func (p *Parser) parseUnary() Node {
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
