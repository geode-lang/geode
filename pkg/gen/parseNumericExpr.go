package gen

import "gitlab.com/nickwanninger/geode/pkg/typesystem"

func (p *Parser) parseNumericExpr() Node {
	t, val := p.token.InferType()

	if t == typesystem.GeodeI64 {
		n := intNode{}
		n.NodeType = nodeInt
		n.Value = val.(int64)
		p.next()
		return n
	}

	if t == typesystem.GeodeF64 {
		n := floatNode{}
		n.NodeType = nodeFloat
		n.Value = val.(float64)
		p.next()
		return n
	}

	if t == typesystem.GeodeI8 {
		n := charNode{}
		n.NodeType = nodeChar
		n.Value = val.(int8)
		p.next()
		return n
	}

	Error(p.token, "invalid number syntax")

	return nil
}
