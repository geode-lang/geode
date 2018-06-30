package ast

import (
	"github.com/nickwanninger/geode/pkg/lexer"
)

var classIndex = 0

func (p *Parser) parseClassDefn() Node {
	p.requires(lexer.TokClassDefn)
	n := classNode{}
	n.NodeType = nodeClass

	p.next()

	p.requires(lexer.TokIdent)
	n.Name = p.token.Value

	p.next()
	nodes := p.parseClassBody()
	n.Variables = make([]variableNode, 0)
	n.Methods = make([]functionNode, 0)
	for _, node := range nodes {
		switch node.Kind() {
		case nodeVariableDecl:
			n.Variables = append(n.Variables, node.(variableNode))
		case nodeFunction:
			n.Methods = append(n.Methods, node.(functionNode))
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
		if p.token.Is(lexer.TokType) {
			// No initializer is allowed in class variable defns
			nodes = append(nodes, p.parseVariableDefn(false))
			p.checkSemiColon()
			continue
		}

		if p.token.Is(lexer.TokFuncDefn) {
			nodes = append(nodes, p.parseFnDefn())
			p.back()
			continue
		}

		// If the block is over.
		if p.token.Is(lexer.TokRightCurly) {
			break
		}
	}
	p.next()

	return nodes
}
