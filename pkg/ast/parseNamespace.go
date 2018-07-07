package ast

import (
	"github.com/nickwanninger/geode/pkg/lexer"
)

func (p *Parser) parseNamespace() Node {
	p.requires(lexer.TokNamespace)
	n := NamespaceNode{}
	n.NodeType = nodeNamespace
	p.next()

	p.requires(lexer.TokIdent)
	n.Name = p.token.Value
	p.next()
	return n
}
