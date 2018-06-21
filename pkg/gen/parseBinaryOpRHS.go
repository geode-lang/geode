package gen

import (
	"gitlab.com/nickwanninger/geode/pkg/parser"
)

func (p *Parser) parseBinaryOpRHS(exprPrec int, lhs Node) Node {
	for {
		_, isBinaryOp := p.binaryOpPrecedence[p.token.Value]
		if !isBinaryOp || p.token.Is(parser.TokSemiColon) {
			return lhs // an expression like '5' will get sent back up to parseTopLevelExpr or parseDefinition from here.
		}
		tokenPrec := p.getTokenPrecedence(p.token.Value)
		if tokenPrec < exprPrec {
			return lhs
		}
		binOp := p.token.Value
		p.next()

		rhs := p.parseUnary()
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
		n := binaryNode{}
		n.NodeType = nodeBinary
		n.OP = binOp
		n.Left = lhs
		n.Right = rhs
		lhs = n
	}
}
