package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
	"github.com/geode-lang/geode/pkg/util/log"
)

func (p *Parser) parseName() string {

	name := ""
	if !p.token.Is(lexer.TokIdent) {
		p.token.SyntaxError()
		log.Fatal("Invalid Name Reference")
	}

	for {
		if p.token.Is(lexer.TokIdent) {
			name += p.token.Value
		} else {
			p.token.SyntaxError()
			log.Fatal("Invalid Name Reference")
		}
		p.next()
		if p.token.Is(lexer.TokNamespaceAccess, lexer.TokDot) {
			name += p.token.Value
			p.next()
			continue
		}
		break
	}

	// fmt.Println(name)

	return name
}
