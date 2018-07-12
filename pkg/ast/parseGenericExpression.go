package ast

import (
	"unicode"

	"github.com/nickwanninger/geode/pkg/lexer"
	"github.com/nickwanninger/geode/pkg/util/log"
)

// parseGenericExpression is a bit differnet. It returns both the generic symbols
// as well as if the generation succeeded. This is because generic syntax shares
// a lot of syntax with math. for example, deciding if this is a generic angle
// bracket or math is hard `.. foo < bar ...`
func (p *Parser) parseGenericExpression(declaration bool) ([]*GenericSymbol, bool) {

	generics := make([]*GenericSymbol, 0)

	p.requires(lexer.TokOper)
	if p.token.Value != "<" {
		return nil, false
	}

	for {
		// Step forward (off of '<' or ',')
		p.next()
		if !p.token.Is(lexer.TokIdent) {
			return nil, false
		}

		g := &GenericSymbol{}
		g.Name = p.token.Value

		// Declaring Generics requires some more things.
		// declaration is true when in a function defn or a type defn
		if declaration {

			if len(g.Name) != 1 || !unicode.IsUpper([]rune(g.Name)[0]) {
				return nil, false
				// p.token.SyntaxError()
				// log.Fatal("Generic type '%s' is not uppercase. Generic types must be capital letters.\n", g.Name)
			}
		}

		generics = append(generics, g)

		p.next()

		if p.token.Is(lexer.TokComma) {
			continue
		}

		if p.token.Value == ">" {
			break
		}

		p.token.SyntaxError()
		log.Fatal("Unexpected token inside generic function\n")
	}

	p.next()
	return generics, true

}
