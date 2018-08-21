package ast

import (
	"strconv"
	"strings"

	"github.com/geode-lang/geode/pkg/typesystem"
	"github.com/geode-lang/geode/pkg/util/log"
)

func (p *Parser) parseNumericExpr() Node {
	t, val := p.token.InferType()

	// Parse Hex Literals
	if strings.Contains(p.token.Value, "x") {
		if !strings.Contains(p.token.Value, "0x") {
			p.token.SyntaxError()
			log.Fatal("Hex Literal must be of the following format: 0x...\n")
		} else {
			n := IntNode{}
			n.TokenReference.Token = p.token
			n.NodeType = nodeInt
			parsed, e := strconv.ParseInt(strings.TrimPrefix(p.token.Value, "0x"), 16, 64)
			if e != nil {
				p.token.SyntaxError()
				log.Fatal("Error decoding hex token\n")
			}
			n.Value = parsed
			p.next()
			return n
		}
	}

	// Parse Binary Literals
	if strings.Contains(p.token.Value, "b") {
		if !strings.Contains(p.token.Value, "0b") {
			p.token.SyntaxError()
			log.Fatal("Binary Literal must be of the following format: 0b...\n")
		} else {
			n := IntNode{}
			n.TokenReference.Token = p.token
			n.NodeType = nodeInt
			parsed, e := strconv.ParseInt(strings.TrimPrefix(p.token.Value, "0b"), 2, 64)
			if e != nil {
				p.token.SyntaxError()
				log.Fatal("Error decoding binary token\n")
			}
			n.Value = parsed
			p.next()
			return n
		}
	}

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
		n.Value = val.(rune)
		p.next()
		return n
	}

	p.Error("invalid number syntax")

	return nil
}
