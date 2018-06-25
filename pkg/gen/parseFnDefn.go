package gen

import (
	"gitlab.com/nickwanninger/geode/pkg/lexer"
	"gitlab.com/nickwanninger/geode/pkg/typesystem"
	"gitlab.com/nickwanninger/geode/pkg/util/log"
)

func (p *Parser) parseFnDefn() functionNode {
	p.next()

	fn := functionNode{}
	fn.NodeType = nodeFunction

	fn.Name = p.token.Value

	p.next()

	if p.token.Type == lexer.TokLeftParen {
		p.next()
		for {
			log.Debug("parseTopLevelStmt - TokFuncDefn\n")
			// If there is an arg
			if p.token.Is(lexer.TokType) {
				fn.Args = append(fn.Args, p.parseVariableDefn(false))

			}
			// Break out case (not a comma, or a right paren)
			if p.token.Is(lexer.TokRightParen) {
				p.next()
				break
			}
			if p.token.Is(lexer.TokComma) {
				continue
			}
		}

	}

	if p.token.Is(lexer.TokType) {
		fn.ReturnType, _ = p.parseType()
	} else {
		fn.ReturnType = typesystem.GlobalTypeMap.GetType("void")
	}

	if p.token.Is(lexer.TokLeftCurly) {
		fn.Body = p.parseBlockStmt()

	} else if p.token.Is(lexer.TokRightArrow) {
		fn.Body = blockNode{}
		fn.Body.NodeType = nodeBlock
		fn.Body.Nodes = make([]Node, 0)
		p.next()

		implReturnValue := p.parseExpression()
		implReturn := returnNode{}
		implReturn.Value = implReturnValue
		fn.Body.Nodes = []Node{implReturn}
		if p.token.Is(lexer.TokSemiColon) {
			p.next()
		} else {
			p.Error("Missing semicolon after implicit return in function %q", fn.Name)
		}
	} else if p.token.Is(lexer.TokElipsis) {
		fn.IsExternal = true
		p.next()
	}
	return fn
}
