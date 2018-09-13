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

	n.Init = p.parseExpression(true)
	n.Cond = p.parseExpression(false)
	n.Step = p.parseExpression(false)
	n.Body = p.parseBlockStmt()

	return n
}
