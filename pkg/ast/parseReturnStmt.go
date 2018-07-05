package ast

func (p *Parser) parseReturnStmt() ReturnNode {
	n := ReturnNode{}
	p.next()

	n.Value = p.parseExpression()

	return n
}
