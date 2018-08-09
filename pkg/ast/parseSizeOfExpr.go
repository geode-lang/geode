package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
)

func (p *Parser) parseSizeofExpr() SizeofNode {
	p.requires(lexer.TokSizeof)

	n := SizeofNode{}
	n.NodeType = nodeSizeof
	n.Token = p.token
	p.next()

	p.requires(lexer.TokLeftParen)
	p.next()

	n.Type = p.parseType()
	p.requires(lexer.TokRightParen)
	p.next()

	return n
}
