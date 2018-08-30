package ast

func (p *Parser) parseReturnStmt() ReturnNode {
	n := ReturnNode{}
	n.TokenReference.Token = p.token
	p.Next()

	n.Value = p.parseExpression()

	p.globTerminator()

	return n
}
