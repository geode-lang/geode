package ast

import "github.com/nickwanninger/geode/pkg/typesystem"

func (p *Parser) parseNumericExpr() Node {
	t, val := p.token.InferType()

	if t == typesystem.GeodeI64 {
		n := IntNode{}
		n.TokenReference.Token = p.token
		n.NodeType = nodeInt
		n.Value = val.(int64)
		p.next()
		return n
	}

	if t == typesystem.GeodeF64 {
		n := FloatNode{}
		n.TokenReference.Token = p.token
		n.NodeType = nodeFloat
		n.Value = val.(float64)
		p.next()
		return n
	}

	if t == typesystem.GeodeI8 {
		n := CharNode{}
		n.TokenReference.Token = p.token
		n.NodeType = nodeChar
		n.Value = val.(int8)
		p.next()
		return n
	}

	p.Error("invalid number syntax")

	return nil
}
