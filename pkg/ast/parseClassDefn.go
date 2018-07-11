package ast

import (
	"github.com/nickwanninger/geode/pkg/lexer"
)

var classIndex = 0

func (p *Parser) parseClassDefn() Node {
	p.requires(lexer.TokClassDefn)
	n := ClassNode{}
	n.NodeType = nodeClass

	p.next()

	p.requires(lexer.TokIdent)
	n.Name = p.token.Value

	p.next()
	nodes := p.parseClassBody()
	n.Variables = make([]VariableNode, 0)
	n.Methods = make([]FunctionNode, 0)
	for _, node := range nodes {
		switch node.Kind() {
		case nodeVariableDecl:
			n.Variables = append(n.Variables, node.(VariableNode))
		case nodeFunction:
			n.Methods = append(n.Methods, node.(FunctionNode))
		}
	}

	// p.next()

	// return nil for now. This causes the program to think the
	// stream of nodes/tokens has ended
	return n
}

func (p *Parser) parseClassBody() []Node {
	p.requires(lexer.TokLeftCurly)
	nodes := make([]Node, 0)

	for {
		p.next()

		if p.token.Is(lexer.TokFuncDefn) {
			// fmt.Println("CLASS - FUNC")

			nodes = append(nodes, p.parseFnDefn())
			p.back()
			continue
		}

		if p.atType() {
			// fmt.Println("CLASS - VAR")
			// No initializer is allowed in class variable defns
			nodes = append(nodes, p.parseVariableDefn(false))
			p.checkSemiColon()
			continue
		}

		// If the block is over.
		if p.token.Is(lexer.TokRightCurly) {
			// fmt.Println("CLASS - FUNC/")
			break
		}
	}
	p.next()

	return nodes
}
