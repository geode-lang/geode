package ast

import "github.com/geode-lang/geode/pkg/lexer"

func (p *Parser) parseBooleanExpr() BooleanNode {

	p.requires(lexer.TokBool)
	n := BooleanNode{}
	n.TokenReference.Token = p.token
	n.NodeType = nodeBool
	n.Value = p.token.Value
	p.Next()
	return n
}
