package ast

import (
	"github.com/nickwanninger/geode/pkg/lexer"
)

func (p *Parser) parseIdentifierExpr(allowVariableDefn bool) Node {

	if allowVariableDefn && p.atType() {
		return p.parseVariableDefn(true)
	}

	name := p.parseName()
	// p.next()

	// Is the next value a paren? If it isnt it is a normal variable reference
	if !p.token.Is(lexer.TokLeftParen) {

		n := VariableNode{}
		n.Name = name
		n.RefType = ReferenceAccessValue
		n.HasValue = true
		n.NodeType = nodeVariable
		// fmt.Println(p.token)
		if p.token.Is(lexer.TokLeftArrow) {
			n.RefType = ReferenceAssign
			p.next()
			n.Body = p.parseExpression()
			return n
		}

		if p.token.Is(lexer.TokCompoundAssignment) {
			operator := string(p.token.Value[1])

			n.RefType = ReferenceAssign

			// The left side is just a reference to the variable.
			left := VariableNode{}
			left.Name = name
			left.RefType = ReferenceAccessValue

			// Parse the right side of the operator
			p.next()
			right := p.parseExpression()

			// Build out the bopNode
			bopNode := BinaryNode{}
			bopNode.Left = left
			bopNode.Right = right
			bopNode.OP = operator

			n.Body = bopNode
		}

		if p.token.Is(lexer.TokLeftBrace) {
			n.IndexExpr = p.parseIndexExpr()
			return n
		}

		return n
	}

	if p.token.Is(lexer.TokLeftParen) {

		// it was a paren, so we need to parse it as if it were a function call
		n := FunctionCallNode{}
		n.Name = name
		n.NodeType = nodeFunctionCall

		for p.next(); p.token.Type != lexer.TokRightParen; {
			switch p.token.Type {
			case lexer.TokComma:
				p.next()
			default:

				arg := p.parseExpression()
				if arg == nil {
					return nil
				}
				n.Args = append(n.Args, arg)
			}
		}
		p.next()

		return n
	}
	return nil
}

// QuickParseIdentifier takes a stream of tokens and lexes them into a single node
func QuickParseIdentifier(src string) Node {
	return NewQuickParser(src).parseIdentifierExpr(true)
}
