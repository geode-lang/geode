package ast

import (
	"github.com/nickwanninger/geode/pkg/lexer"
)

func (p *Parser) parsePrimary() Node {
	// fmt.Println(lexer.GetTokenName(p.token.Type))
	switch p.token.Type {

	case lexer.TokIdent:
		return p.parseIdentifierExpr()
	case lexer.TokNumber:
		return p.parseNumericExpr()
	case lexer.TokLeftParen:
		return p.parseParenExpr()
	case lexer.TokString:
		return p.parseStringExpr()
	// case tokEndOfTokens:
	// 	return nil // this token should not be skipped
	default:
		// p.next()
		return nil
	}
}
