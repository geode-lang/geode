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
	TokBool

	TokDot
	TokElipsis
	TokOper
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

	TokInfo

	TokCompoundAssignment

	TokQuestionMark

	TokFor
	TokWhile
	TokIf
	TokElse
	TokReturn
	TokFuncDefn
	TokClassDefn
	TokNamespace
	TokLet
	TokAs
	TokNil

	TokDependency

	TokType

	TokComma

	TokIdent
	TokSymbol

	TokComment
)
