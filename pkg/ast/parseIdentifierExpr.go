package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
)

func (p *Parser) parseIdentifierExpr(allowVariableDefn bool) Node {

	if allowVariableDefn && p.atType() {
		return p.parseVariableDefn(true)
	}

	nameToken := p.token

	name := p.parseName()
	namedRef := NewNamedReference(name)

	var generics []*GenericSymbol

	state := p.save()
	genValid := false
	if p.token.Is(lexer.TokOper) && p.token.Value == "<" {
		generics, genValid = p.parseGenericExpression(false)
	}

	if !genValid {
		p.restore(state)
	}
	// p.next()

	// Is the next value a paren? If it isnt it is a normal variable reference
	if !p.token.Is(lexer.TokLeftParen) {

		n := VariableNode{}
		n.TokenReference.Token = nameToken
		n.Name = namedRef
		n.RefType = ReferenceAccessValue
		n.HasValue = true
		n.NodeType = nodeVariable
		// fmt.Println(p.token)
		if p.token.Is(lexer.TokLeftArrow) {
			assignment := VariableAssignNode{}
			assignment.Target = namedRef
			p.next()
			assignment.Body = p.parseExpression()
			return assignment
		}

		if p.token.Is(lexer.TokCompoundAssignment) {
			operator := string(p.token.Value[0])

			n.RefType = ReferenceAssign

			// The left side is just a reference to the variable.
			left := VariableNode{}
			left.TokenReference.Token = p.token
			left.Name = namedRef
			left.RefType = ReferenceAccessValue

			// Parse the right side of the operator
			p.next()
			right := p.parseExpression()

			// Build out the bopNode
			bopNode := BinaryNode{}
			bopNode.TokenReference.Token = p.token
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
		n.TokenReference.Token = nameToken
		n.Name = namedRef
		n.NodeType = nodeFunctionCall
		n.Generics = generics

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
