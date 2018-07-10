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

	unaryOp := p.token.Value

	p.next()
	operand := p.parseUnary()

	if unaryOp == "&" {
		if operand.Kind() == nodeVariable {
			// Update operand's RefType if it is a nodeVariable
			n := (operand).(VariableNode)
			n.RefType = ReferenceAccessStackAddress
			operand = n
		}
	}
	if operand != nil {
		n := UnaryNode{}
		n.NodeType = nodeUnary
		n.Operator = unaryOp
		n.Operand = operand
		return n
	}
	return nil
}
