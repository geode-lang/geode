package ast

func (p *Parser) parseReturnStmt() ReturnNode {
	n := ReturnNode{}
	n.TokenReference.Token = p.token
	p.next()

	n.Value = p.parseExpression()

	return n
}
