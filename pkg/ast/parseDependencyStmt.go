package ast

import (
	"strings"

	"github.com/nickwanninger/geode/pkg/lexer"
)

func (p *Parser) parseDependencyStmt() Node {
	d := dependencyNode{}
	d.NodeType = nodeDependency
	p.requires(lexer.TokDependency)
	if p.token.Value == "link" {
		d.CLinkage = true
	}
	p.next()

	d.Paths = make([]string, 0)

	for {
		p.requires(lexer.TokString)
		thisPath := strings.Trim(p.token.Value, "\"")
		d.Paths = append(d.Paths, thisPath)
		p.next()
		if p.token.Type != lexer.TokComma {
			break
		}
	}

	return d
}
