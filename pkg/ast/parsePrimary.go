package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
)

func (p *Parser) parsePrimary() Node {
	// fmt.Println(lexer.GetTokenName(p.token.Type))
	switch p.token.Type {
	case lexer.TokSizeof:
		return p.parseSizeofExpr()
	case lexer.TokIdent:
		return p.parseIdentifierExpr(false)
	case lexer.TokNumber:
		return p.parseNumericExpr()
	case lexer.TokLeftParen:
		return p.parseParenExpr()
	case lexer.TokString:
		return p.parseStringExpr()

	case lexer.TokLeftBrace:
		return p.parseArrayDecl()
	}

	return nil
}
