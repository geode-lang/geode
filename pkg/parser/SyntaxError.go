package parser

import (
	"fmt"
	"github.com/danielchatfield/go-chalk"
	"strings"
)

// SyntaxError prints a string that points an error in the source code
func SyntaxError(tok Token, message string) string {
	col := tok.StartCol - 1
	line := tok.StartRow - 1
	length := tok.EndCol - tok.StartCol
	src := string(*tok.SourceCode)

	lines := strings.Split(src, "\n")
	if len(lines) < line-1 {
		return ""
	}

	underline := "╹" // start out by assuming the token is of len 1 (point toqjust the token)

	// if the token is longer, make the underline longer as well
	if length > 1 {
		underline = fmt.Sprintf("┗%s┛", strings.Repeat("━", length-2))
	}

	border := fmt.Sprintf("%s\n", strings.Repeat("·", len(lines[line])))

	errorString := ""
	errorString += fmt.Sprintf("Syntax Error: %s\n", message)
	errorString += fmt.Sprintf("\tLine: %d\n", line)
	errorString += fmt.Sprintf("\tCol: %d\n", col)
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
