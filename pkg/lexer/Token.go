package lexer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nickwanninger/geode/pkg/typesystem"
	"github.com/nickwanninger/geode/pkg/util/color"
)

// TokenIsOperator will return if a given token is an operator or not
func TokenIsOperator(t TokenType) bool {
	return t > TokOperatorStart && t < TokOperatorEnd
}

// Token is a token in the program
type Token struct {
	source *Sourcefile
	Type   TokenType `json:"type,omitempty"`
	Value  string    `json:"value,omitempty"`
	Pos    int       `json:"start_pos"`
	EndPos int       `json:"end_pos"`
	Line   int       `json:"line"`
	Column int       `json:"column"`
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

	return fmt.Sprintf("%q - %s", t.Value, t.Type.String())
}

// SyntaxError prints a formatted syntax error
func (t *Token) SyntaxError() {
	// fmt.Println(t.Value)
	// width := runewidth.StringWidth(t.Value)
	// underline := util.Red(strings.Repeat("^", width))
	// fmt.Println(underline)

	lineMargins := 2
	src := t.source.String()

	// Highlight the source string at the error
	src = src[:t.Pos] + color.Red(src[t.Pos:t.EndPos]) + src[t.EndPos:]
	// Replace tabs with a fixed number of spaces
	src = strings.Replace(src, "\t", "    ", -1)
	lines := strings.Split(src, "\n")

	// Start printing
	fmt.Printf("\nSyntax error:\n")
	fmt.Printf(color.Blue("   | %s:%d\n"), t.source.Path, t.Line)
	fmt.Printf(color.Blue("   |\n"))
	for i, line := range lines {
		ln := i + 1

		if ln >= t.Line-lineMargins && ln <= t.Line+lineMargins {
			lineString := ""
			lineNumber := "  "
			if ln == t.Line {
				lineNumber = fmt.Sprintf("%2d", ln)
			}
			lineString = color.Blue(fmt.Sprintf("%s |", lineNumber))
			fmt.Printf("%s %s\n", lineString, line)
		}

	}
	fmt.Printf(color.Blue("   |\n\n"))

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
