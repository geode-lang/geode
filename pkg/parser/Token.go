package parser

import (
	"fmt"
	"github.com/danielchatfield/go-chalk"
	"strings"
)

// TokenType gives information to what kind of token was parsed
type TokenType int

const (
	tokenWhitespace TokenType = iota
	tokenEOF
	tokenError
	tokenNumber
	tokenString
)

// Token is a single token, and can contain the source of the token, the type of the token, and where in the file the token exists
type Token struct {
	typ TokenType
	val string
	idx int
}

// Type - Returns the type of the token
func (t Token) Type() TokenType {
	return t.typ
}

// Value - Returns the string value of the token
func (t Token) Value() string {
	return t.val
}

// Index - Returns the integer index of the token in the source code
func (t Token) Index() int {
	return t.idx
}

func (t Token) String() string {
	switch t.typ {
	case tokenEOF:
		return "EOF"
	case tokenError:
		return t.val
	}

	if len(t.val) > 30 {
		return fmt.Sprintf("\"%.30s...\"", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}

// SyntaxError prints a string that points an error in the source code
func SyntaxError(index int, length int, src string, filename string, message string) string {
	// Row and Col are zero indexed
	line := 0
	col := 0

	for i := 0; i < index; i++ {
		col++
		if src[i] == '\n' {
			line++
			col = 0
		}
	}

	fmt.Println(line, col)

	lines := strings.Split(src, "\n")
	if len(lines) < line-1 {
		return ""
	}

	underline := "╹" // start out by assuming the token is of len 1 (point to just the token)

	// if the token is longer, make the underline longer as well
	if length > 1 {
		underline = fmt.Sprintf("┗%s┛", strings.Repeat("━", length-2))
	}

	border := fmt.Sprintf("%s\n", strings.Repeat("·", len(lines[line])))

	errorString := ""
	errorString += fmt.Sprintf("ERROR: Syntax Error in %q:\n", filename)
	errorString += fmt.Sprintf("%s\n", message)
	errorString += fmt.Sprintf("\tLine: %d\n", line)
	errorString += fmt.Sprintf("\tCol: %d\n", col)
	errorString += fmt.Sprintf("\tIndex: %d\n", index)

	errorString += border

	// Print the line before if we can
	if line > 0 {
		errorString += fmt.Sprintf("%s\n", lines[line-1])
	}

	errorString += fmt.Sprintf("%s\n", chalk.Red(lines[line]))
	errorString += fmt.Sprintf("%s%s\n", strings.Repeat(" ", col), underline)
	// Print the line after if we can
	if len(lines) > line+1 {
		errorString += fmt.Sprintf("%s\n", lines[line+1])
	}

	errorString += border
	errorString += fmt.Sprintf("\n")

	return errorString
}
