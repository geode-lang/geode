package parser

// Alot of this based off this video:
// https://www.youtube.com/watch?v=HxaD_trXwRE

import (
	"fmt"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
	"os"
)

// Tokens is a list of all tokens
var Tokens = [][]string{
	{"CHAR", `'.'`},
	{"STRING", `"(\\"|.)*?"`},
	{"NUMBER", `[+-]?[0-9]*\.?[0-9]+`},
	{"MUL", `\*`},
	{"PLUS", `\+`},
	{"MINUS", `-`},
	{"DIV", `/`},
	{"EXP", `\^`},

	{"DEREFERENCE", `@`},
	{"REFERENCE", `\*`},

	{"ASSIGNMENT", `:=`},
	{"EQUALITY", `=`},

	{"RIGHT_PAREN", `\)`},
	{"LEFT_PAREN", `\(`},

	{"RIGHT_CURLY", `}`},
	{"LEFT_CURLY", `{`},

	{"RIGHT_BRACE", `\[`},
	{"LEFT_BRACE", `\]`},

	{"ACT_DEFN", `act\s`},
	{"ACTARROW", `->`},
	// The main parser won't work correctly and will just look these up later
	{"IF", ""},
	{"ELSE", ""},
	{"RETURN", ""},

	{"IDENTIFIER", `[a-zA-Z_][a-zA-Z0-9_]*`},

	{"COMMENT", `\/\/[^\n]*`},
	{"COMMENT", `{-.*-}`},
	{"WHITESPACE", `\s+`},

	{"LT", `<`},
	{"LTE", `<=|≤`},
	{"GT", `>`},
	{"GTE", `>=|≥`},
}

var keywordrmap = map[string]string{
	"return": "RETURN",
	"if":     "IF",
	"else":   "ELSE",
}

var tokmap map[string]int

var tokRegexMap map[string]string

func init() {
	tokmap = make(map[string]int)
	tokRegexMap = make(map[string]string)
	for id, val := range Tokens {
		tokmap[val[0]] = id
		if val[1] != "" {
			tokRegexMap[val[1]] = val[0]
		}

	}
}

// LexState - an internal rep of the lexer
type LexState struct {
	lexer *lexmachine.Lexer
}

// Lex - takes a string and turns it into tokens
func (s *LexState) Lex(text []byte) ([]Token, error) {
	scanner, err := s.lexer.Scanner(text)
	toks := []Token{}
	if err != nil {
		return nil, err
	}
	for tk, err, eof := scanner.Next(); !eof; tk, err, eof = scanner.Next() {
		if ui, is := err.(*machines.UnconsumedInput); ui != nil && is {
			e := err.(*machines.UnconsumedInput)
			scanner.TC = ui.FailTC
			fmt.Println(SyntaxError(e.FailLine, e.StartColumn, e.FailColumn-e.StartColumn-1, string(text), "Tokenize Failed"))
			os.Exit(1)
		} else if err != nil {
			return nil, err
		} else {
			// I don't like lexmachine's token, so I will convert it to my own
			to := *tk.(*lexmachine.Token)
			t := Token{}
			t.StartCol = to.StartColumn
			t.StartRow = to.StartLine
			t.EndCol = to.EndColumn
			t.EndRow = to.EndLine
			t.Type = to.Type
			t.Lexeme = to.Lexeme
			t.Value = string(to.Value.(string))
			toks = append(toks, t)
		}
	}

	return toks, nil
}

// NewLexer produces a new lexer and poluates it with the configuration
func NewLexer() *LexState {

	getToken := func(tokenType int) lexmachine.Action {
		return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			kw, isKwInMap := keywordrmap[string(m.Bytes)]
			if isKwInMap {
				return s.Token(tokmap[kw], string(m.Bytes), m), nil
			}
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
