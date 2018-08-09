package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
)

func (p *Parser) parseExpression() Node {
	lhs := p.parseUnary()
	if lhs == nil {
		return nil
	}

	if p.token.Is(lexer.TokAs) {
		return p.parseCastExpr(lhs)
	}
	return p.parseBinaryOpRHS(1, lhs)
}
