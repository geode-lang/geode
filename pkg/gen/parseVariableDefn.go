package gen

import (
	"gitlab.com/nickwanninger/geode/pkg/lexer"
	"gitlab.com/nickwanninger/geode/pkg/types"
)

func (p *Parser) parseVariableDefn(allowDefn bool) variableNode {
	n := variableNode{}
	n.NodeType = nodeVariableDecl

	if p.token.Is(lexer.TokType) {
		n.Type = types.GlobalTypeMap.GetType(p.token.Value)
		p.next()

		if p.token.Is(lexer.TokRightBrace) {
			n.IsArray = true

			p.next()
			if p.token.Is(lexer.TokLeftBrace) {
				p.next()
			} else {
				Error(p.token, "Malformed variable array definition")
			}
		}
		if p.token.Is(lexer.TokIdent) {
			n.Name = p.token.Value
		} else {
			Error(p.token, "Missing Variable name")
		}

	} else {
		Error(p.token, "Invalid variable declaration")
	}

	if allowDefn && p.nextToken.Is(lexer.TokAssignment) {
		n.HasValue = true
		p.next()
		p.next()
		n.Body = p.parseExpression()
	}

	return n
}
