package ast

func (p *Parser) parseExpression() Node {
	lhs := p.parseUnary()
	if lhs == nil {
		return nil
	}
	return p.parseBinaryOpRHS(1, lhs) // TODO: check on this value wrt our : = and 0 val for not found instead of tut's -1
}
