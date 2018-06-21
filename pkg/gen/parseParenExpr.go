package gen

import (
	"gitlab.com/nickwanninger/geode/pkg/lexer"
)

func (p *Parser) parseParenExpr() Node {
	// skip over the parens
	p.next()
	v := p.parseExpression()
	if v == nil {
		return nil
	}
	if p.token.Type != lexer.TokRightParen {
		Error(p.token, "expected ')'")
		return nil
	}
	p.next()
	return v
}
