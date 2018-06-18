package ast

func (p *Parser) parseStringExpr() Node {
	n := stringNode{}
	n.NodeType = nodeString
	n.Value = p.token.Value
	p.next()
	return n
}
