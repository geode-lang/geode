package ast

import (
	"github.com/nickwanninger/geode/pkg/lexer"
)

var blkidx = 0

func (p *Parser) parseBlockStmt() blockNode {
	p.requires(lexer.TokLeftCurly)
	blk := blockNode{}
	blk.NodeType = nodeBlock

	for {
		p.next()
		if p.token.Is(lexer.TokType) {
			blk.Nodes = append(blk.Nodes, p.parseVariableDefn(true))
			p.checkSemiColon()
			continue
		}

		if p.token.Is(lexer.TokReturn) {
			blk.Nodes = append(blk.Nodes, p.parseReturnStmt())
			p.checkSemiColon()
			continue
		}

		if p.token.Is(lexer.TokIdent) {
			blk.Nodes = append(blk.Nodes, p.parseIdentifierExpr())
			p.checkSemiColon()
			continue
		}

		if p.token.Is(lexer.TokIf) {
			blk.Nodes = append(blk.Nodes, p.parseIfStmt())
			// We have to step back a token when exiting any statement that parses a block
			// so the next time around in this block will step forward
			p.back()
			continue
			// p.checkSemiColon()
		}

		if p.token.Is(lexer.TokWhile) {
			blk.Nodes = append(blk.Nodes, p.parseWhileStmt())
			// We have to step back a token when exiting any statement that parses a block
			// so the next time around in this block will step forward
			p.back()
			continue
		}

		if p.token.Is(lexer.TokFor) {
			blk.Nodes = append(blk.Nodes, p.parseForStmt())
			// We have to step back a token when exiting any statement that parses a block
			// so the next time around in this block will step forward
			p.back()
			continue
		}

		// If the block is over.
		if p.token.Is(lexer.TokRightCurly) {
			// Set the cursor on the next character

			break
		}
	}
	p.next()

	blkidx++

	return blk
}
