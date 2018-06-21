package gen

import (
	"gitlab.com/nickwanninger/geode/pkg/parser"
)

func (p *Parser) parseBlockStmt() blockNode {
	blk := blockNode{}
	blk.NodeType = nodeBlock

	for {
		p.next()
		if p.token.Is(parser.TokType) {
			blk.Nodes = append(blk.Nodes, p.parseVariableDefn(true))
		}

		if p.token.Is(parser.TokReturn) {
			blk.Nodes = append(blk.Nodes, p.parseReturnStmt())
		}

		// If the block is over.
		if p.token.Is(parser.TokRightCurly) {
			// Set the cursor on the next character
			p.next()
			break
		}
	}

	return blk
}
