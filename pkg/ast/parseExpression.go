package ast

func (p *Parser) parseExpression() Node {
	lhs := p.parseUnary()
	if lhs == nil {
		return nil
	}
	return p.parseBinaryOpRHS(1, lhs)
}

// QuickParseExpression takes a stream of tokens and lexes them into a single node
func QuickParseExpression(src string) Node {
	return NewQuickParser(src).parseExpression()
}
