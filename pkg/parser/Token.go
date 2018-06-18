package parser

import (
	"encoding/json"
	"github.com/nickwanninger/act/pkg/types"
	"strconv"
)

// Token is a token in the program
type Token struct {
	SourceCode *string `json:"-"`
	Type       int     `json:"type,omitempty"`
	Value      string  `json:"value,omitempty"`
	Pos        int     `json:"start_pos"`
	EndPos     int     `json:"end_pos"`
}

// Is - returns if the token is a certain type as a string
func (t *Token) Is(a int) bool {
	return t.Type == a
}

func (t Token) String() string {

	m := make(map[string]interface{})
	m["type"] = GetTokenName(t.Type)
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
func (t Token) InferType() (*types.VarType, interface{}) {
	if t.Type == TokNumber {
		intval, intErr := strconv.ParseInt(t.Value, 10, 8)
		if intErr == nil {
			return types.DefaultIntType, intval
		}

		floatval, floatErr := strconv.ParseFloat(t.Value, 8)
		if floatErr == nil {
			return types.DefaultFloatType, floatval
		}
	}

	// if t.Type == TokString {
	// 	return types.DefaultStringType, t.Value
	// }

	// if t.Type == TokChar {
	// 	return types.DefaultCharType, t.Value
	// }

	return nil, nil
}

// func (t Token) String() string {

// 	return fmt.Sprintf("%12s : %q", GetTokenName(t.Type), t.Value)
// }
