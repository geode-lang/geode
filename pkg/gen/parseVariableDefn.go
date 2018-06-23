package gen

import (
	"gitlab.com/nickwanninger/geode/pkg/lexer"
	"gitlab.com/nickwanninger/geode/pkg/typesystem"
)

func (p *Parser) parseVariableDefn(allowDefn bool) variableNode {
	n := variableNode{}
	n.NodeType = nodeVariableDecl

	if p.token.Is(lexer.TokType) {
		n.Type = typesystem.GlobalTypeMap.GetType(p.token.Value)
		p.next()

		if p.token.Is(lexer.TokRightBrace) {
			n.IsArray = true

			p.next()
			if p.token.Is(lexer.TokLeftBrace) {
				p.next()
			} else {
				p.Error("Malformed variable array definition")
			}
		}
		if p.token.Is(lexer.TokIdent) {
			n.Name = p.token.Value
		} else {
			p.Error("Missing Variable name")
		}

	} else {
		p.Error("Invalid variable declaration")
	}

	if allowDefn && p.peek(1).Is(lexer.TokAssignment) {
		n.HasValue = true
		p.next()
		p.next()
		n.Body = p.parseExpression()
	}

	return n
}
