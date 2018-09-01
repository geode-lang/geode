package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
)

func (p *Parser) parsePrimary() Node {

	switch p.token.Type {

	case lexer.TokSizeof:
		return p.parseSizeofExpr()
	case lexer.TokIdent:
		return p.parseIdentifierExpr(false)
	case lexer.TokBool:
		return p.parseBooleanExpr()
	case lexer.TokNumber:
		return p.parseNumericExpr()
	case lexer.TokLeftParen:
		return p.parseParenExpr()
	case lexer.TokString:
		return p.parseStringExpr()
	case lexer.TokNil:
		n := NilNode{NodeType: nodeNil}
		p.Next()
		return n
	case lexer.TokChar:
		return p.parseCharExpr()

	case lexer.TokLeftBrace:
		return p.parseArrayDecl()
	}

	return nil
}
