package parser

import (
	"fmt"
	"github.com/danielchatfield/go-chalk"
	"strings"
)

// TokenType gives information to what kind of token was parsed
type TokenType string

// Token is a single token, and can contain the source of the token, the type of the token, and where in the file the token exists
type Token struct {
	Type       TokenType
	Filename   string
	SourceCode *string
	Value      string
	Row        int
	Col        int
}

// SyntaxError prints a string that points an error in the source code
func SyntaxError(row int, col int, len int, sourcecode string, filename string) string {
	// Row and Col are zero indexed
	line := row + 1
	lineSrc := strings.Split(sourcecode, "\n")[row]
	valLen := col
	underline := "╹" // start out by assuming the token is of len 1 (point to just the token)
	// if the token is longer, make the underline longer as well
	if valLen > 1 {
		underline = fmt.Sprintf("┗%s┛", strings.Repeat("━", valLen-2))
	}

	errorString := ""

	errorString += fmt.Sprintf("ERROR: Syntax Error in '%s':\n", filename)
	errorString += fmt.Sprintf("\tLine: %d\n", line)
	errorString += fmt.Sprintf("\tCol: %d\n", col)
	errorString += fmt.Sprintf("\t%s\n", chalk.Red(lineSrc))
	errorString += fmt.Sprintf("\t%s%s\n", strings.Repeat(" ", col), underline)
	errorString += fmt.Sprintf("\n")

	return errorString
}
