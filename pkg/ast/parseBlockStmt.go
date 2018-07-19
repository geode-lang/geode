package ast

import (
	"github.com/nickwanninger/geode/pkg/lexer"
)

var blkidx = 0

func (p *Parser) parseBlockStmt() BlockNode {
	p.requires(lexer.TokLeftCurly)
	blk := BlockNode{}
	blk.TokenReference.Token = p.token
	blk.NodeType = nodeBlock

	for {
		p.next()
		// if p.token.Is(lexer.TokLet) {
		// 	blk.Nodes = append(blk.Nodes, p.parseVariableDefn(true))
		// 	p.checkSemiColon()
		// 	continue
		// }

		if p.token.Is(lexer.TokReturn) {
			blk.Nodes = append(blk.Nodes, p.parseReturnStmt())

			p.checkSemiColon()
			continue
		}

		if p.token.Is(lexer.TokIdent) {
			blk.Nodes = append(blk.Nodes, p.parseIdentifierExpr(true))
			p.checkSemiColon()
			continue
		}

		if p.token.Is(lexer.TokIf) {
			blk.Nodes = append(blk.Nodes, p.parseIfStmt())
			p.back()
			continue
			// p.checkSemiColon()
		}

		if p.token.Is(lexer.TokWhile) {
			blk.Nodes = append(blk.Nodes, p.parseWhileStmt())
			p.back()
			continue
		}

		if p.token.Is(lexer.TokFor) {
			blk.Nodes = append(blk.Nodes, p.parseForStmt())
			p.back()
			continue
		}

		// If the block is over.
		if p.token.Is(lexer.TokRightCurly) {
			break
		}
	}
	p.next()

	blkidx++

	return blk
}
