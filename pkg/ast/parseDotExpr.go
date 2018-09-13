package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
)

func (p *Parser) parseDotExpr(base Reference) Reference {

	n := DotReference{}
	n.Token = p.token
	n.NodeType = nodeDot
	n.Base = base
	p.requires(lexer.TokDot)
	p.Next()
	p.requires(lexer.TokIdent)
	n.Field = NewIdentNode(p.token.Value)
	p.Next()

	if p.token.Is(lexer.TokDot) {
		return p.parseDotExpr(n)
	}

	return n
}

// QuickParseExpression takes a stream of tokens and lexes them into a single node
func QuickParseExpression(src string) Node {
	return NewQuickParser(src).parseExpression(true)
}
