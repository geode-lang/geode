package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
)

func (p *Parser) parseParenExpr() Node {
	// skip over the parens
	p.next()
	v := p.parseExpression()
	if v == nil {
		return nil
	}

	if p.token.Type != lexer.TokRightParen {
		p.Error("expected ')'")
		return nil
	}
	p.next()

	return v
}
