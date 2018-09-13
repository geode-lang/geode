package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
	"github.com/geode-lang/geode/pkg/util/log"
)

func (p *Parser) parseIdentifierExpr(allowVariableDefn bool) Node {

	if allowVariableDefn && (p.atType() || p.token.Is(lexer.TokLet)) {
		return p.parseVariableDefn(true)
	}

	nameToken := p.token

	name, _ := p.parseName()

	var target Reference
	target = NewIdentNode(name)

	for p.token.Is(lexer.TokDot) {
		target = p.parseDotExpr(target)

	}

	defer p.globTerminator()

	// Is the next value a paren? If it isnt it is a normal variable reference
	if !p.token.Is(lexer.TokLeftParen) {
		n := VariableNode{}
		n.TokenReference.Token = nameToken
		n.Name = target
		n.RefType = ReferenceAccessValue
		n.HasValue = true
		n.NodeType = nodeVariable

		if p.token.Is(lexer.TokLeftArrow) {
			assignment := AssignmentNode{}
			assignment.Token = p.token
			assignment.NodeType = nodeAssignment
			var toStore Assignable
			toStore = target

			assignment.Assignee = toStore

			p.Next()

			store := p.parseExpression(false)
			if access, isAccess := store.(Accessable); isAccess {
				assignment.Value = access
			} else {
				store.SyntaxError()
				log.Fatal("Not a valid value to assign to a variable node. (Not accessable)\n")
			}
			return assignment
		}

		if p.token.Is(lexer.TokCompoundAssignment) {
			operator := string(p.token.Value[0])

			assignment := AssignmentNode{}
			assignment.Token = p.token
			assignment.NodeType = nodeAssignment
			assignment.Assignee = target
			p.Next()

			// The left side is just a reference to the variable.
			left := VariableNode{}
			left.TokenReference.Token = p.token
			left.Name = target
			left.RefType = ReferenceAccessValue

			// Parse the right side of the operator
			right := p.parseExpression(false)
			// Build out the bopNode
			bopNode := BinaryNode{}
			bopNode.TokenReference.Token = p.token
			bopNode.Left = left
			bopNode.Right = right
			bopNode.OP = operator

			assignment.Value = bopNode
			return assignment
		}

		// Handle all things arrays accessing and assignment
		if p.token.Is(lexer.TokLeftBrace) {
			var s Node
			s = n
			for p.token.Is(lexer.TokLeftBrace) {
				if a, isAccessable := s.(Accessable); isAccessable {
					s = p.parseSubscriptExpr(a)
				} else {
					s.SyntaxError()
					log.Fatal("Unable to index.")
				}
			}

			if p.token.Is(lexer.TokLeftArrow) {
				assignment := AssignmentNode{}
				if a, is := s.(Assignable); is {
					assignment.Assignee = a
				} else {
					s.SyntaxError()
					log.Fatal("%T is not Assignable\n", s)
				}

				p.Next()

				val := p.parseExpression(false)

				if a, is := val.(Accessable); is {
					assignment.Value = a
				} else {
					val.SyntaxError()
					log.Fatal("Value is not an accessable value\n")
				}

				return assignment

			}
			return s
		}
		return p.parsePossibleFunctionCall(n)
	}

	if p.token.Is(lexer.TokLeftParen) {
		// it was a paren, so we need to parse it as if it were a function call
		n := FunctionCallNode{}
		n.NodeType = nodeFnCall
		n.TokenReference.Token = nameToken

		callable, ok := target.(Callable)
		if !ok {
			log.Fatal("%T is not Callable", target)
		}
		n.Name = callable
		n.NodeType = nodeFunctionCall

		for p.Next(); p.token.Type != lexer.TokRightParen; {
			switch p.token.Type {
			case lexer.TokComma:
				p.Next()
			default:

				arg := p.parseExpression(false)
				if arg == nil {
					return nil
				}
				n.Args = append(n.Args, arg)
			}
		}
		p.Next()

		return n
	}
	return nil
}

// QuickParseIdentifier takes a stream of tokens and lexes them into a single node
func QuickParseIdentifier(src string) Node {
	return NewQuickParser(src).parseIdentifierExpr(true)
}
