package ast

import (
	"strings"

	"github.com/geode-lang/geode/pkg/lexer"
	"github.com/geode-lang/geode/pkg/util/log"
)

func (p *Parser) parseClassDefn() Node {
	p.requires(lexer.TokClassDefn)
	n := ClassNode{}
	n.TokenReference.Token = p.token
	n.NodeType = nodeClass

	p.Next()

	if !p.token.Is(lexer.TokType) {
		p.token.SyntaxError()
		log.Fatal("Class names must be capitalized. Use %q instead\n", strings.Title(p.token.Value))
	}
	n.Name = p.token.Value

	p.Context().ClassNames[n.Name] = p.token

	p.Next()
	nodes := p.parseClassBody()
	n.Variables = make([]VariableDefnNode, 0)
	n.Methods = make([]FunctionNode, 0)

	for _, node := range nodes {
		switch node.Kind() {
		case nodeVariableDecl:
			n.Variables = append(n.Variables, node.(VariableDefnNode))
		case nodeFunction:
			fn := node.(FunctionNode)
			n.Methods = append(n.Methods, fn)
		}
	}

	// p.Next()

	// return nil for now. This causes the program to think the
	// stream of nodes/tokens has ended
	return n
}

func (p *Parser) parseClassBody() []Node {

	p.requires(lexer.TokLeftCurly)
	nodes := make([]Node, 0)
	p.Next()

	for {
		if p.token.Is(lexer.TokFuncDefn) {
			fn := p.parseFunctionNode()
			fn.IsMethod = true
			nodes = append(nodes, fn)
			continue
		}

		if p.atType() {
			// No initializer is allowed in class variable defns
			nodes = append(nodes, p.parseVariableDefn(false))
			p.globTerminator()
			continue
		}

		// If the block is over.
		if p.token.Is(lexer.TokRightCurly) {
			break
		}
	}
	p.Next()

	return nodes
}
