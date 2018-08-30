package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
	"github.com/geode-lang/geode/pkg/util/log"
)

var blkidx = 0

func (p *Parser) parseBlockStmt() BlockNode {

	p.requires(lexer.TokLeftCurly)
	blk := BlockNode{}
	blk.TokenReference.Token = p.token
	blk.NodeType = nodeBlock
	p.Next()
	for {

		if p.token.Is(lexer.TokReturn) {
			blk.Nodes = append(blk.Nodes, p.parseReturnStmt())
			continue
		}

		if p.token.Is(lexer.TokIdent, lexer.TokLet) {
			node := p.parseIdentifierExpr(true)
			blk.Nodes = append(blk.Nodes, node)
			continue
		}

		if p.token.Is(lexer.TokIf) {
			blk.Nodes = append(blk.Nodes, p.parseIfStmt())
			continue
		}

		if p.token.Is(lexer.TokWhile) {
			blk.Nodes = append(blk.Nodes, p.parseWhileStmt())
			continue
		}

		if p.token.Is(lexer.TokFor) {
			blk.Nodes = append(blk.Nodes, p.parseForStmt())
			continue
		}

		// If the block is over.
		if p.token.Is(lexer.TokRightCurly) {
			break
		}

		p.token.SyntaxError()
		log.Fatal("Unknown token in block statement\n")
	}
	p.Next()

	blkidx++

	return blk
}

// forkBlockParser returns a new, forked parser that only has a subset of tokens that
// contain an entire block. ex: starting at {, ending at }.
// This funciton correctly nests.
func (p *Parser) forkBlockParser() *Parser {
	p.requires(lexer.TokLeftCurly)
	parser := p.fork()
	parser.tokenIndex = 0
	index := p.tokenIndex
	offset := 0
	// Start out at 1 level of nesting because the token is
	nesting := 1

	for nesting != 0 {
		offset++
		tok := p.Next()
		if tok.Is(lexer.TokLeftCurly) {
			nesting++
		} else if tok.Is(lexer.TokRightCurly) {
			nesting--
		}
	}
	offset++
	p.Next()
	tokens := p.tokens[index : index+offset]
	parser.tokens = tokens
	parser.reset()
	return parser
}
