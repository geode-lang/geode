package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
	"github.com/geode-lang/geode/pkg/util/log"
)

func validTypeInfoTokens(t lexer.Token) bool {
	allowed := map[string]bool{
		"*": true,
		"?": true,
		// "[": true,
		// "]": true,
	}
	_, ok := allowed[t.Value]
	return ok

}

func (p *Parser) atType() bool {

	if !p.token.Is(lexer.TokIdent) {
		return false
	}
	offset := 1
	for validTypeInfoTokens(p.peek(offset)) {
		offset++
	}

	if p.peek(offset).Type == lexer.TokIdent {
		return true
	}

	return false
}

func (p *Parser) parseType() (t GeodeTypeRef) {
	p.requires(lexer.TokIdent)

	t.Name = p.token.Value
	p.next()

	for {

		if p.token.Is(lexer.TokQuestionMark) {
			if t.Unknown {
				log.Fatal("Multiple Unknown Type operators for %q used.\n", t.Name)
			}

			t.Unknown = true
			p.next()
			continue
		}

		if p.token.Is(lexer.TokOper) {
			for _, c := range p.token.Value {
				if c == '*' {
					t.PointerLevel++
				}
			}

			p.next()
			continue
		}

		break

	}
	return t
}
