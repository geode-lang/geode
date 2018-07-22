package ast

import (
	"fmt"

	"github.com/nickwanninger/geode/pkg/lexer"
)

func validTypeReferenceToken(t lexer.Token) bool {
	allowed := map[string]bool{
		"*": true,
		// "[": true,
		// "]": true,
	}
	_, ok := allowed[t.Value]
	return ok

}

func (p *Parser) atType() bool {
	offset := 1
	for validTypeReferenceToken(p.peek(offset)) {
		offset++
	}

	// fmt.Println(p.peek(offset))

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

		if p.token.Is(lexer.TokOper) {
			for _, c := range p.token.Value {
				fmt.Printf("  -  %c\n", c)
				if c == '*' {
					t.PointerLevel++
					fmt.Println(t.PointerLevel)
				}
			}

			p.next()
			continue
		}

		break

	}
	return t
}
