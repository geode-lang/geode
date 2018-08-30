package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
)

func (p *Parser) parseSizeofExpr() SizeofNode {
	p.requires(lexer.TokSizeof)

	n := SizeofNode{}
	n.NodeType = nodeSizeof
	n.Token = p.token
	p.Next()

	p.requires(lexer.TokLeftParen)
	p.Next()

	n.Type = p.parseType()
	p.requires(lexer.TokRightParen)
	p.Next()

	return n
}
