package gen

import (
	"github.com/nickwanninger/act/pkg/parser"
)

func (p *Parser) parsePrimary() Node {
	// fmt.Println(parser.GetTokenName(p.token.Type))
	switch p.token.Type {
	// case tokIdentifier:
	// 	return p.parseIdentifierExpr()
	// case tokIf:
	// 	return p.parseIfExpr()
	// case tokFor:
	// 	return p.parseForExpr()
	// case tokVariable:
	// 	return p.parseVarExpr()
	case parser.TokIdent:
		return p.parseIdentifierExpr()
	case parser.TokNumber:
		return p.parseNumericExpr()
	case parser.TokLeftParen:
		return p.parseParenExpr()
	case parser.TokString:
		return p.parseStringExpr()
	// case tokEndOfTokens:
	// 	return nil // this token should not be skipped
	default:
		// p.next()
		return nil
	}
}
