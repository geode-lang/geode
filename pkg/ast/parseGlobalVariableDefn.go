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
			n.Name = NewIdentNode(p.token.Value)
			p.Next()
		} else if p.token.Is(lexer.TokOper) && p.token.Value == "=" {

		} else {
			n.SyntaxError()
			log.Fatal("Invalid Global variable declaration\n")
		}

	} else {
		p.token.SyntaxError()
		log.Fatal("Invalid Global variable declaration")
	}

	if p.token.Is(lexer.TokOper) && p.token.Value == "=" {
		p.Next()
		n.Body = p.parseExpression(false)
	} else if p.token.Is(lexer.TokElipsis) {
		n.External = true
		p.Next()
	}

	if !n.External {
		p.globTerminator()
		// p.Next()
	}

	return n
}
