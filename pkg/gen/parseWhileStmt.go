package gen

import (
	"github.com/nickwanninger/geode/pkg/lexer"
)

var whileStmtIndex = 0

func (p *Parser) parseWhileStmt() Node {
	p.requires(lexer.TokWhile)
	n := whileNode{}
	n.NodeType = nodeWhile
	n.Index = whileStmtIndex
	whileStmtIndex++
	p.next()

	n.If = p.parseExpression()
	p.requires(lexer.TokLeftCurly)

	n.Body = p.parseBlockStmt()
	return n
}
