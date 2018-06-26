package lexer

import (
	"encoding/json"
	"strconv"
	"strings"

	"gitlab.com/nickwanninger/geode/pkg/typesystem"
)

//go:generate stringer -type=TokenType $GOFILE

// TokenType -
type TokenType int

// Assigning tokens integer values
const (
	TokError TokenType = iota
	TokWhitespace
	TokChar
	TokString
	TokNumber

	TokElipsis
	TokOper
	TokPtr

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

	TokType

	TokComma

	TokIdent

	TokComment
)

// TokenIsOperator will return if a given token is an operator or not
func TokenIsOperator(t TokenType) bool {
	return t > TokOperatorStart && t < TokOperatorEnd
}

// Token is a token in the program
type Token struct {
	SourceCode *string   `json:"-"`
	Type       TokenType `json:"type,omitempty"`
	Value      string    `json:"value,omitempty"`
	Pos        int       `json:"start_pos"`
	EndPos     int       `json:"end_pos"`
}

// Is - returns if the token is a certain type as a string
func (t Token) Is(a TokenType) bool {
	return t.Type == a
}

func (t Token) String() string {

	m := make(map[string]interface{})
	m["type"] = string(t.Type)
	m["value"] = t.Value

	encoded, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(encoded)
}

func (t *Token) buildEndPos(endCol int, endRow int) {

	src := *t.SourceCode

	end := 0

	row := 1
	col := 1

	for i, c := range src {
		col++
		if c == '\n' {
			row++
			col = 1
		}
		if row == endRow && col == endCol {
			end = i
			break
		}
	}

	t.EndPos = end
}

// InferType takes some token and guesses the type
func (t Token) InferType() (*typesystem.VarType, interface{}) {
	if t.Type == TokNumber {
		intval, intErr := strconv.ParseInt(t.Value, 10, 64)
		if intErr == nil {
			return typesystem.GeodeI64, intval
		}

		floatval, floatErr := strconv.ParseFloat(t.Value, 64)
		if floatErr == nil {
			return typesystem.GeodeF64, floatval
		}
	}

	if t.Type == TokChar {
		c := strings.Trim(t.Value, "'")[0]
		return typesystem.GeodeI8, c
	}

	return nil, nil
}
