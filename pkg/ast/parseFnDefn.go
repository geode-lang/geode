package ast

import (
	"github.com/nickwanninger/geode/pkg/lexer"
)

func (p *Parser) parseFnDefn() FunctionNode {
	// func, pure, etc...
	declarationKeyword := p.token.Value

	fn := FunctionNode{}
	fn.TokenReference.Token = p.token
	fn.NodeType = nodeFunction
	fn.DeclKeyword = DeclKeywordFunc

	fn.line = p.token.Line
	fn.column = p.token.Column

	p.next()

	if declarationKeyword == "pure" {
		fn.DeclKeyword = DeclKeywordPure
	}

	if p.token.Type == lexer.TokIdent && p.token.Value == "nomangle" {
		fn.Nomangle = true
		p.next()
	}

	rawNameString := p.parseName()
	fn.Name = NewNamedReference(rawNameString)

	// fmt.Println(declarationKeyword, fn.Name)

	// The main function should never be mangled
	if rawNameString == "main" {
		fn.Nomangle = true
	}

	if p.token.Type == lexer.TokOper && p.token.Value == "<" {
		fn.Generics, _ = p.parseGenericExpression(true)
	}

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
				// fn.Nomangle = true
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
		fn.Body = BlockNode{}
		fn.Body.NodeType = nodeBlock
		fn.Body.Nodes = make([]Node, 0)
		fn.ImplicitReturn = true
		p.next()

		implReturnValue := p.parseExpression()
		implReturn := ReturnNode{}
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

// QuickParseFunction takes a stream of tokens and lexes them into a single node
func QuickParseFunction(src string) Node {
	return NewQuickParser(src).parseFnDefn()
}
