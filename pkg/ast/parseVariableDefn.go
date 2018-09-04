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

	if p.token.Is(lexer.TokLet) {
		n.NeedsInference = true
		p.Next()

	} else if p.atType() {
		n.Type = p.parseType()
	} else {
		p.token.SyntaxError()
		log.Fatal("let: Invalid variable declaration\n")
	}

	if p.token.Is(lexer.TokIdent) {
		n.Name = NewNamedReference(p.token.Value)
		p.Next()
	} else if p.token.Is(lexer.TokAssignment) {

	} else {
		n.SyntaxError()
		log.Fatal("type: Invalid variable declaration\n")
	}

	if p.token.Is(lexer.TokAssignment) {
		if allowDefn {
			n.HasValue = true
			p.Next()
			n.Body = p.parseExpression()
		} else {
			log.Fatal("Variable Initialization of '%s' is not allowed in it's context\n", n.Name)
		}
	} else if n.NeedsInference {
		n.SyntaxError()
		log.Fatal("When declaring a variable with let, it must have an assignment\n")
	}

	return n
}
