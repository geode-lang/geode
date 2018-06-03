package parser

// Token is a token in the program
type Token struct {
	Type     int
	Value    string
	Lexeme   []byte
	StartCol int
	StartRow int
	EndCol   int
	EndRow   int
}
