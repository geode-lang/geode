package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
)

func (p *Parser) parseArrayDecl() Node {

	n := ArrayNode{}
	n.NodeType = nodeArray
	p.requires(lexer.TokLeftBrace)
	p.next()
	elements := make([]Node, 0)

	for {
		if p.token.Is(lexer.TokRightBrace) {
			break
		}
		if p.token.Is(lexer.TokComma) {
			p.next()
			continue
		}
		elements = append(elements, p.parseExpression())
	}

	n.Elements = elements
	n.Length = len(elements)
	p.requires(lexer.TokRightBrace)
	p.next()

	return n
}
