package parser

// Alot of this based off this video:
// https://www.youtube.com/watch?v=HxaD_trXwRE

import (
	"fmt"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

// Tokens is a list of all tokens
var Tokens = [][]string{

	{"STRING", `"[^"]*"`},
	{"NUMBER", `[+-]?[0-9]*\.?[0-9]+`},
	{"MUL", `\*`},
	{"PLUS", `\+`},
	{"MINUS", `-`},
	{"DIV", `/`},
	{"EXP", `\^`},
	{"ASSIGNMENT", `<-`},
	{"EQUALITY", `=`},

	{"RIGHT_PAREN", `\)`},
	{"LEFT_PAREN", `\(`},

	{"RIGHT_CURLY", `}`},
	{"LEFT_CURLY", `{`},

	{"RIGHT_BRACE", `\[`},
	{"LEFT_BRACE", `\]`},

	{"ACT", `act\s`},
	{"IF", `if`},
	{"ELSE", `else`},
	{"RETURN", `return`},
	{"BOOLEAN", `true|false`},

	{"IDENTIFIER", `[a-zA-Z_][a-zA-Z0-9_]*`},

	{"COMMENT", `\/\/[^\n]*`},
	{"COMMENT", `{-.*-}`},
	{"WHITESPACE", `\s+`},

	{"ACTARROW", `->`},
	{"LT", `<`},
	{"LTE", `<=|≤`},
	{"GT", `>`},
	{"GTE", `>=|≥`},
}
var tokmap map[string]int

var tokRegexMap map[string]string

func init() {
	tokmap = make(map[string]int)
	tokRegexMap = make(map[string]string)
	for id, val := range Tokens {
		tokmap[val[0]] = id
		tokRegexMap[val[1]] = val[0]
	}
}

// LexState - an internal rep of the lexer
type LexState struct {
	lexer *lexmachine.Lexer
}

// Lex - takes a string and turns it into tokens
func (s *LexState) Lex(text []byte) ([]lexmachine.Token, error) {
	scanner, err := s.lexer.Scanner(text)
	toks := []lexmachine.Token{}
	if err != nil {
		return nil, err
	}
	for tk, err, eof := scanner.Next(); !eof; tk, err, eof = scanner.Next() {
		if ui, is := err.(*machines.UnconsumedInput); ui != nil && is {
			scanner.TC = ui.FailTC
			fmt.Printf("skipping %v", ui)
		} else if err != nil {
			return nil, err
		} else {
			toks = append(toks, *tk.(*lexmachine.Token))
		}
	}
	return toks, nil
}

// NewLexer produces a new lexer and poluates it with the configuration
func NewLexer() *LexState {

	getToken := func(tokenType int) lexmachine.Action {
		return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			if tokenType == tokmap["WHITESPACE"] {
				return nil, nil
			}
			return s.Token(tokenType, string(m.Bytes), m), nil
		}
	}

	lexer := lexmachine.NewLexer()

	for k, v := range tokRegexMap {
		lexer.Add([]byte(k), getToken(tokmap[v]))
	}
	s := &LexState{}

	s.lexer = lexer
	return s
}

func GetTokenName(id int) string {
	return Tokens[id][0]
}
