package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
)

var whileStmtIndex = 0

func (p *Parser) parseWhileStmt() Node {
	p.requires(lexer.TokWhile)
	n := WhileNode{}
	n.TokenReference.Token = p.token
	n.NodeType = nodeWhile
	n.Index = whileStmtIndex
	whileStmtIndex++
	p.Next()

	n.If = p.parseExpression(false)
	p.requires(lexer.TokLeftCurly)

	n.Body = p.parseBlockStmt()
	return n
}
