package ast

func (p *Parser) parseReturnStmt() returnNode {
	n := returnNode{}
	p.next()
	n.Value = p.parseExpression()
	return n
}
