package ast

import (
	"strings"

	"github.com/geode-lang/geode/pkg/lexer"
)

func (p *Parser) parseDependencyStmt() Node {
	d := DependencyNode{}
	d.TokenReference.Token = p.token
	d.NodeType = nodeDependency
	p.requires(lexer.TokDependency)
	if p.token.Value == "link" {
		d.CLinkage = true
	}
	p.Next()

	d.Paths = make([]string, 0)

	for {
		p.requires(lexer.TokString)
		thisPath := strings.Trim(p.token.Value, "\"")
		d.Paths = append(d.Paths, thisPath)
		p.Next()
		if p.token.Type != lexer.TokComma {
			break
		}
	}

	return d
}
