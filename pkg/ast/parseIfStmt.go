package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
)

var ifStmtIndex = 0

func (p *Parser) parseIfStmt() Node {
	p.requires(lexer.TokIf)
	n := IfNode{}
	n.TokenReference.Token = p.token
	n.NodeType = nodeIf
	n.Index = ifStmtIndex
	ifStmtIndex++

	p.Next()

	n.If = p.parseExpression(false)

	p.requires(lexer.TokLeftCurly)

	n.Then = p.parseBlockStmt()

	if p.token.Is(lexer.TokElse) {
		p.Next()
		n.Else = p.parseBlockStmt()
	}

	// increment the ifstmtindex for the next time around
	return n
}
