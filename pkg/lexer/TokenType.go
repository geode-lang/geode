package lexer

//go:generate stringer -type=TokenType $GOFILE

// TokenType -
type TokenType int

// Assigning tokens integer values
const (
	TokError  TokenType = iota
	TokNoEmit           // NoEmit is to be ignored by the lexer
	TokWhitespace
	TokChar
	TokString
	TokNumber

	TokElipsis
	TokOper
	TokPtr

	TokNamespaceAccess

	TokOperatorStart
	TokStar
	TokPlus
	TokMinus
	TokDiv
	TokExp
	TokLT
	TokLTE
	TokGT
	TokGTE
	TokOperatorEnd

	TokSemiColon

	TokDefereference
	TokReference

	TokAssignment
	TokEquality

	TokRightParen
	TokLeftParen

	TokRightCurly
	TokLeftCurly

	TokRightBrace
	TokLeftBrace

	TokRightArrow
	TokLeftArrow

	TokFor
	TokWhile
	TokIf
	TokElse
	TokReturn
	TokFuncDefn
	TokClassDefn
	TokNamespace

	TokDependency

	TokType

	TokComma

	TokIdent

	TokComment
)
