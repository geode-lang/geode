package gen

import (
	"gitlab.com/nickwanninger/geode/pkg/lexer"
)

func (p *Parser) parseBlockStmt() blockNode {
	blk := blockNode{}
	blk.NodeType = nodeBlock

	for {
		p.next()
		if p.token.Is(lexer.TokType) {
			blk.Nodes = append(blk.Nodes, p.parseVariableDefn(true))
		}

		if p.token.Is(lexer.TokReturn) {
			blk.Nodes = append(blk.Nodes, p.parseReturnStmt())
		}

		// If the block is over.
		if p.token.Is(lexer.TokRightCurly) {
			// Set the cursor on the next character
			p.next()
			break
		}
	}

	return blk
}
