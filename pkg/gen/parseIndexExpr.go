package gen

import (
	"github.com/nickwanninger/geode/pkg/lexer"
)

func (p *Parser) parseIndexExpr() Node {
	p.requires(lexer.TokLeftBrace)
	p.next()
	n := p.parseExpression()
	p.requires(lexer.TokRightBrace)
	p.next()
	return n
}
