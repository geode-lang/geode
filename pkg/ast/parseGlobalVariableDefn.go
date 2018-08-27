package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
	"github.com/geode-lang/geode/pkg/util/log"
)

func (p *Parser) parseGlobalVariableDecl() GlobalVariableDeclNode {
	n := GlobalVariableDeclNode{}
	n.Token = p.token
	n.NodeType = nodeGlobalDecl
	n.TokenReference.Token = p.token

	if p.atType() {
		n.Type = p.parseType()

		if p.token.Is(lexer.TokIdent) {
			n.Name = NewNamedReference(p.token.Value)
			p.next()
		} else if p.token.Is(lexer.TokAssignment) {

		} else {
			n.SyntaxError()
			log.Fatal("Invalid Global variable declaration\n")
		}

	} else {
		p.token.SyntaxError()
		log.Fatal("Invalid Global variable declaration")
	}

	if p.token.Is(lexer.TokAssignment) {
		p.next()
		n.Body = p.parseExpression()
	} else if p.token.Is(lexer.TokElipsis) {
		n.External = true
		p.next()
	}

	if !n.External {
		p.allowSemiColon()
		p.next()
	}

	return n
}
