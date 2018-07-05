package ast

import "github.com/nickwanninger/geode/pkg/lexer"

func (p *Parser) parseTypeCast() Node {
	n := CastNode{}
	n.NodeType = nodeCast
	n.To = p.token.Value

	p.next()

	p.requires(lexer.TokLeftParen)
	p.next()
	// Parse the value that will be casted
	n.From = p.parseExpression()
	p.requires(lexer.TokRightParen)
	p.next()
	return n
}
