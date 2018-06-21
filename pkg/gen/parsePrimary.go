package gen

import (
	"gitlab.com/nickwanninger/geode/pkg/lexer"
)

func (p *Parser) parsePrimary() Node {
	// fmt.Println(lexer.GetTokenName(p.token.Type))
	switch p.token.Type {
	// case tokIdentifier:
	// 	return p.parseIdentifierExpr()
	// case tokIf:
	// 	return p.parseIfExpr()
	// case tokFor:
	// 	return p.parseForExpr()
	// case tokVariable:
	// 	return p.parseVarExpr()
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
