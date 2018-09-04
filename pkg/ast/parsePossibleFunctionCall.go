package ast

import "github.com/geode-lang/geode/pkg/lexer"

func (p *Parser) parsePossibleFunctionCall(left Node) Node {
	_, ok := left.(Callable)
	if !ok || p.token.Is(lexer.TokLeftParen) {
		return left
	}

	return left
}
