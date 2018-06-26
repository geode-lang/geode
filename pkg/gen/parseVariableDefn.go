package gen

import (
	"github.com/nickwanninger/geode/pkg/lexer"
	"github.com/nickwanninger/geode/pkg/util/log"
)

func (p *Parser) parseVariableDefn(allowDefn bool) variableNode {
	p.requires(lexer.TokType)
	n := variableNode{}
	n.NodeType = nodeVariableDecl

	if p.token.Is(lexer.TokType) {
		n.Type, n.IsPointer = p.parseType()

		if p.token.Is(lexer.TokIdent) {
			n.Name = p.token.Value
			p.next()
		} else {
			log.Debug("%s\n", p.token)
			log.Fatal("Missing Variable name")
		}

	} else {
		log.Fatal("Invalid variable declaration")
	}

	if allowDefn && p.token.Is(lexer.TokAssignment) {
		n.HasValue = true
		p.next()
		n.Body = p.parseExpression()
	}

	return n
}
