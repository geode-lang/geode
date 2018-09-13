package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
)

func (p *Parser) parseBinaryOpRHS(exprPrec int, lhs Node) Node {

	// parse plain binary operator
	for {
		tokenPrec, isBinaryOp := p.binaryOpPrecedence[p.token.Value]
		if !isBinaryOp || p.token.Is(lexer.TokSemiColon) {
			return lhs
		}

		if tokenPrec < exprPrec {
			return lhs
		}
		binOp := p.token.Value
		p.Next()

		// right hand sides will never have a declaration, so pass false
		rhs := p.parseUnary(false)
		if rhs == nil {
			return nil
		}

		nextPrec := p.getTokenPrecedence(p.token.Value)
		if tokenPrec < nextPrec {
			rhs = p.parseBinaryOpRHS(tokenPrec+1, rhs)
			if rhs == nil {
				return nil
			}
		}
		n := BinaryNode{}
		n.TokenReference.Token = p.token
		n.NodeType = nodeBinary
		n.OP = binOp
		n.Left = lhs
		n.Right = rhs
		lhs = n
	}
}
