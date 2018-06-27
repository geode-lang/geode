package ast

import (
	"github.com/nickwanninger/geode/pkg/lexer"
)

var forStmtIndex = 0

func (p *Parser) parseForStmt() Node {
	p.requires(lexer.TokFor)
	n := forNode{}
	n.NodeType = nodeFor
	n.Index = forStmtIndex
	forStmtIndex++
	p.next()

	n.Init = p.parseVariableDefn(true)
	p.requires(lexer.TokSemiColon)

	p.next()

	n.Cond = p.parseExpression()
	p.requires(lexer.TokSemiColon)
	p.next()

	n.Step = p.parseExpression()

	p.requires(lexer.TokLeftCurly)
	n.Body = p.parseBlockStmt()

	return n
}
