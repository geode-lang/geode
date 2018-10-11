package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
	"github.com/geode-lang/geode/pkg/util/log"
)

func (p *Parser) parseFunctionNode() FunctionNode {
	// func, pure, etc...
	declarationKeyword := p.token.Value

	fn := FunctionNode{}
	fn.TokenReference.Token = p.token
	fn.NodeType = nodeFunction
	fn.DeclKeyword = DeclKeywordFunc

	fn.line = p.token.Line
	fn.column = p.token.Column

	p.Next()

	if declarationKeyword == "pure" || declarationKeyword == "λ" {
		fn.DeclKeyword = DeclKeywordPure
	}

	if p.token.Type == lexer.TokIdent && p.token.Value == "nomangle" {
		fn.Nomangle = true
		p.Next()
	}

	rawNameString, _ := p.parseName()
	fn.Name = NewIdentNode(rawNameString)

	// The main function should never be mangled
	if rawNameString == "main" {
		fn.Nomangle = true
	}

	if p.token.Type == lexer.TokLeftParen {
		p.Next()

		for {

			// Parse a function argument
			if p.token.Is(lexer.TokIdent, lexer.TokType) {

				typ := p.parseType()

				if !p.token.Is(lexer.TokIdent) {
					p.token.SyntaxError()
					log.Fatal("invalid function argument\n")
				}

				for p.token.Is(lexer.TokIdent) {
					arg := FunctionArg{}
					arg.Type = typ
					arg.Name = p.token.Value
					p.Next()
					fn.Args = append(fn.Args, arg)
				}

			}

			if p.token.Is(lexer.TokElipsis) {
				fn.Variadic = true
				// Variadic functions are external, or should be. This means they shouldn't be mangled
				fn.Nomangle = true
				p.Next()
			}

			// Break out case (not a comma, or a right paren)
			if p.token.Is(lexer.TokRightParen) {
				p.Next()
				break
			}

			if p.token.Is(lexer.TokComma) {
				// Skip over the comma
				p.Next()
				continue
			}
		}

	}

	if p.token.Is(lexer.TokType) {
		fn.ReturnType = p.parseType()
	} else {
		fn.ReturnType = TypeNode{}
		fn.ReturnType.Name = "void"
		fn.ReturnType.PointerLevel = 0
		fn.ReturnType.Unknown = false
	}
	// fmt.Println(p.token.Value)

	if p.token.Is(lexer.TokLeftCurly) {
		fn.BodyParser = p.forkBlockParser()
	} else if p.token.Is(lexer.TokRightArrow, lexer.TokOper) {

		if p.token.Is(lexer.TokOper) && p.token.Value != "=" {
			p.token.SyntaxError()
			log.Fatal("unexpected token %q in function declaration\n", p.token.Value)
		}

		if p.token.Is(lexer.TokRightArrow) {
			log.Deprecated("Use of an arrow function will be removed. Replace '->' with '=' (%s)\n", p.token.FileInfo())
		}
		fn.Body = BlockNode{}
		fn.Body.NodeType = nodeBlock
		fn.Body.Nodes = make([]Node, 0)
		fn.ImplicitReturn = true
		p.Next()

		implReturnValue := p.parseExpression(false)
		implReturn := ReturnNode{}
		implReturn.Value = implReturnValue
		fn.Body.Nodes = []Node{implReturn}
		p.globTerminator()
	} else if p.token.Is(lexer.TokElipsis) {
		fn.External = true
		// External functions should not be mangled
		fn.Nomangle = true
		p.Next()
	} else {
		p.token.SyntaxError()
		log.Fatal("")
	}

	for _, arg := range fn.Args {
		if arg.Type.Unknown {
			fn.HasUnknownType = true
		}
	}

	return fn
}

// QuickParseFunction takes a stream of tokens and lexes them into a single node
func QuickParseFunction(src string) Node {
	return NewQuickParser(src).parseFunctionNode()
}
