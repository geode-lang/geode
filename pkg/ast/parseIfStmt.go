package ast

import (
	"github.com/nickwanninger/geode/pkg/lexer"
)

var ifStmtIndex = 0

func (p *Parser) parseIfStmt() Node {
	p.requires(lexer.TokIf)
	n := ifNode{}
	n.NodeType = nodeIf
	n.Index = ifStmtIndex
	ifStmtIndex++

	p.next()

	n.If = p.parseExpression()
	p.requires(lexer.TokLeftCurly)

	n.Then = p.parseBlockStmt()

	if p.token.Is(lexer.TokElse) {
		// fmt.Println(p.token)
		p.next()
		n.Else = p.parseBlockStmt()
	}
	// increment the ifstmtindex for the next time around
	return n
}
