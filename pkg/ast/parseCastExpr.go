package ast

import "github.com/geode-lang/geode/pkg/lexer"

func (p *Parser) parseCastExpr(source Node) Node {
	p.requires(lexer.TokAs)
	n := CastNode{}
	n.Token = p.token
	n.NodeType = nodeCast
	n.Source = source
	p.Next()
	n.Type = p.parseType()

	return n
}
