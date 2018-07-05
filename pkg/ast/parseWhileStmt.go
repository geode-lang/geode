package ast

import (
	"github.com/nickwanninger/geode/pkg/lexer"
)

var whileStmtIndex = 0

func (p *Parser) parseWhileStmt() Node {
	p.requires(lexer.TokWhile)
	n := WhileNode{}
	n.NodeType = nodeWhile
	n.Index = whileStmtIndex
	whileStmtIndex++
	p.next()

	n.If = p.parseExpression()
	p.requires(lexer.TokLeftCurly)

	n.Body = p.parseBlockStmt()
	return n
}
