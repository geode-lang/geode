package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
	"github.com/geode-lang/geode/pkg/util/log"
)

func (p *Parser) parseVariableDefn(allowDefn bool) VariableDefnNode {
	n := VariableDefnNode{}

	n.Token = p.token
	n.NodeType = nodeVariableDecl
	n.TokenReference.Token = p.token
	if p.atType() {
		n.Typ = p.parseType()
	} else {
		p.token.SyntaxError()
		log.Fatal("let: Invalid variable declaration\n")
	}

	if p.token.Is(lexer.TokIdent) {
		n.Name = NewIdentNode(p.token.Value)
		p.Next()
	} else {
		n.SyntaxError()
		log.Fatal("type: Invalid variable declaration\n")
	}

	if p.token.Is(lexer.TokAssignment) {
		if allowDefn {
			n.HasValue = true
			p.Next()
			n.Body = p.parseExpression(false)
		} else {
			log.Fatal("Variable Initialization of '%s' is not allowed in it's context\n", n.Name)
		}
	} else if n.NeedsInference {
		n.SyntaxError()
		log.Fatal("When declaring a variable with let, it must have an assignment\n")
	}

	return n
}
