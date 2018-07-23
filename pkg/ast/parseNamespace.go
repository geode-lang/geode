package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
)

func (p *Parser) parseNamespace() Node {
	p.requires(lexer.TokNamespace)
	n := NamespaceNode{}
	n.TokenReference.Token = p.token
	n.NodeType = nodeNamespace
	p.next()

	p.requires(lexer.TokIdent)
	n.Name = p.token.Value
	p.next()
	return n
}
