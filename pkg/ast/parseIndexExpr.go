package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
)

func (p *Parser) parseIndexExpr() Node {
	p.requires(lexer.TokLeftBrace)
	p.next()
	n := p.parseExpression()
	p.requires(lexer.TokRightBrace)
	p.next()
	return n
}
