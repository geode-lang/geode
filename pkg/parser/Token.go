package parser

import (
	"encoding/json"
)

// Token is a token in the program
type Token struct {
	SourceCode *string `json:"source_code,omitempty"`
	Type       int     `json:"type,omitempty"`
	Value      string  `json:"value,omitempty"`
	Lexeme     []byte  `json:"lexeme,omitempty"`
	Pos        int     `json:"pos"`
	StartCol   int     `json:"start_col,omitempty"`
	StartRow   int     `json:"start_row,omitempty"`
	EndCol     int     `json:"end_col,omitempty"`
	EndRow     int     `json:"end_row,omitempty"`
}

// Is - returns if the token is a certain type as a string
func (t *Token) Is(a int) bool {
	return t.Type == a
}

func (t Token) String() string {

	m := make(map[string]interface{})
	m["type"] = GetTokenName(t.Type)
	m["value"] = t.Value
	m["start"], m["end"] = t.getPosRange()

	encoded, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(encoded)
}

func (t Token) getPosRange() (int, int) {

	src := *t.SourceCode

	start := t.Pos
	end := 0

	row := 1
	col := 1

	for i, c := range src {
		col++
		if c == '\n' {
			row++
			col = 1
		}
		if row == t.EndRow && col == t.EndCol {
			end = i
			break
		}
	}

	return start, end
}

// func (t Token) String() string {

// 	return fmt.Sprintf("%12s : %q", GetTokenName(t.Type), t.Value)
// }
