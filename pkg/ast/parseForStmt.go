package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
)

var forStmtIndex = 0

func (p *Parser) parseForStmt() Node {
	p.requires(lexer.TokFor)
	n := ForNode{}
	n.TokenReference.Token = p.token
	n.NodeType = nodeFor
	n.Index = forStmtIndex
	forStmtIndex++
	p.Next()

	n.Init = p.parseIdentifierExpr(true)

	n.Cond = p.parseExpression()
	n.Step = p.parseExpression()

	n.Body = p.parseBlockStmt()

	return n
}
