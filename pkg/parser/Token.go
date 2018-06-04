package parser

// Token is a token in the program
type Token struct {
	SourceCode *[]byte `json:"source_code,omitempty"`
	Type       int     `json:"type,omitempty"`
	Value      string  `json:"value,omitempty"`
	Lexeme     []byte  `json:"lexeme,omitempty"`
	StartCol   int     `json:"start_col,omitempty"`
	StartRow   int     `json:"start_row,omitempty"`
	EndCol     int     `json:"end_col,omitempty"`
	EndRow     int     `json:"end_row,omitempty"`
}

// Is - returns if the token is a certain type as a string
func (t *Token) Is(a TokenType) bool {
	return t.Type == GetTokenId(a)
}
