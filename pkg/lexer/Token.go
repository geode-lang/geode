package lexer

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/nickwanninger/geode/pkg/typesystem"
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

// TokenIsOperator will return if a given token is an operator or not
func TokenIsOperator(t TokenType) bool {
	return t > TokOperatorStart && t < TokOperatorEnd
}

// Token is a token in the program
type Token struct {
	source      *Sourcefile `json:"-"`
	Type        TokenType   `json:"type,omitempty"`
	Value       string      `json:"value,omitempty"`
	Pos         int         `json:"start_pos"`
	EndPos      int         `json:"end_pos"`
	StartLine   int         `json:"start_line"`
	StartColumn int         `json:"start_column"`
	EndLine     int         `json:"end_line"`
	EndColumn   int         `json:"end_column"`
}

// Is - returns if the token is a certain type as a string
func (t Token) Is(types ...TokenType) bool {
	for _, a := range types {
		if t.Type == a {
			return true
		}
	}
	return false
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

// SyntaxError returns a formatted syntax error
func (t *Token) SyntaxError() {
	src := strings.Replace(t.source.String(), "\t", "    ", -1)
	fmt.Printf("\n=========================================================\n")
	fmt.Printf("\nSyntax error!\n")
	fmt.Printf("Value: %q\n", t.Value)
	fmt.Printf("Line: %d\n", t.StartLine)
	fmt.Printf("Column: %d\n", t.StartColumn)

	lines := strings.Split(src, "\n")
	lineMargins := 3
	fmt.Printf("\n       ...\n")
	for i, line := range lines {
		ln := i + 1

		if ln >= t.StartLine-lineMargins && ln <= t.EndLine+lineMargins {
			es := "  "
			if ln >= t.StartLine && ln <= t.EndLine {
				es = "->"
			}
			fmt.Printf("%3d %s %s\n", ln, es, line)
		}

	}
	fmt.Printf("       ...\n")
	fmt.Printf("\n=========================================================\n\n")

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
