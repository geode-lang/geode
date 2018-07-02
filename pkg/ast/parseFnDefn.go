package ast

import (
	"github.com/nickwanninger/geode/pkg/lexer"
)

func (p *Parser) parseFnDefn() functionNode {
	p.next()

	fn := functionNode{}
	fn.NodeType = nodeFunction

	if p.token.Type == lexer.TokIdent && p.token.Value == "nomangle" {
		fn.Nomangle = true
		p.next()
	}

	fn.Name = p.token.Value

	// The main function should never be mangled
	if fn.Name == "main" {
		fn.Nomangle = true
	}

	p.next()

	if p.token.Type == lexer.TokLeftParen {
		p.next()

		for {

			// If there is an arg
			if p.token.Is(lexer.TokIdent) {
				fn.Args = append(fn.Args, p.parseVariableDefn(false))
			}

			if p.token.Is(lexer.TokElipsis) {
				fn.Variadic = true
				// Variadic functions are external, or should be. This means they shouldn't be mangled
				fn.Nomangle = true
				p.next()
			}

			// Break out case (not a comma, or a right paren)
			if p.token.Is(lexer.TokRightParen) {
				p.next()
				break
			}

			if p.token.Is(lexer.TokComma) {
				// Skip over the comma
				p.next()
				continue
			}
		}

	}

	if p.token.Is(lexer.TokIdent) {
		fn.ReturnType = p.parseType()
	} else {
		fn.ReturnType = GeodeTypeRef{false, 0, "void"}
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
		fn.External = true
		// External functions should not be mangled
		fn.Nomangle = true
		p.next()
	}
	return fn
}