package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
	"github.com/geode-lang/geode/pkg/util/log"
)

func (p *Parser) parseSubscriptExpr(source Accessable) Node {

	subN := SubscriptNode{}
	subN.Source = source
	p.requires(lexer.TokLeftBrace)
	p.Next()
	index := p.parseExpression()

	if indexAc, isAccessable := index.(Accessable); isAccessable {
		subN.Index = indexAc
	} else {
		index.SyntaxError()
		log.Fatal("Unable to index by an expression that isn't an accessable value\n")
	}
	p.requires(lexer.TokRightBrace)
	p.Next()
	return subN
}
