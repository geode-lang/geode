package ast

import (
	"fmt"

	"github.com/geode-lang/geode/pkg/lexer"
)

func (p *Parser) parseName() (string, error) {

	name := ""
	if !p.token.Is(lexer.TokIdent) {
		return "", fmt.Errorf("Invalid Name Reference")
	}

	for {
		if p.token.Is(lexer.TokIdent) {
			name += p.token.Value
		} else {
			return "", fmt.Errorf("Invalid Name Reference")
		}
		p.next()
		if p.token.Is(lexer.TokNamespaceAccess) {
			name += p.token.Value
			p.next()
			continue
		}
		break
	}

	// fmt.Println(name)

	return name, nil
}
