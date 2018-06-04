package parser

import (
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
	"os"
)

// TokenType - The type of token as a string
type TokenType string

// TokenRegexRelation - allows a relationship between a token type and a certain regex
type TokenRegexRelation struct {
	token TokenType
	regex string
}

// Tokens is a list of all tokens
var Tokens = []TokenRegexRelation{
	{"ERROR", ""},
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

	{"ASSIGNMENT", `<-`},
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

	{"TYPE", ""},

	{"COMMA", `,`},

	{"IDENTIFIER", `[a-zA-Z_][a-zA-Z0-9_]*`},

	{"COMMENT", `\/\/[^\n]*`},
	{"COMMENT", `{-.*-}`},
	{"WHITESPACE", `\s+`},

	{"LT", `<`},
	{"LTE", `<=|≤`},
	{"GT", `>`},
	{"GTE", `>=|≥`},
}

var keywordrmap = map[string]TokenType{
	"return": "RETURN",
	"if":     "IF",
	"else":   "ELSE",

	// We also determine type mapping in here as well
	"num": "TYPE",
}

var tokmap map[TokenType]int

var tokRegexMap map[string]TokenType

func init() {
	tokmap = make(map[TokenType]int)
	tokRegexMap = make(map[string]TokenType)
	for id, val := range Tokens {
		tokmap[val.token] = id
		if val.regex != "" {
			tokRegexMap[val.regex] = val.token
		}

	}
}

// LexState - an internal rep of the lexer
type LexState struct {
	lexer  *lexmachine.Lexer
	Tokens chan Token
	Done   bool
}

// Lex - takes a string and turns it into tokens
func (s *LexState) Lex(text []byte) error {
	scanner, err := s.lexer.Scanner(text)
	if err != nil {
		return err
	}
	for tk, err, eof := scanner.Next(); !eof; tk, err, eof = scanner.Next() {
		if ui, is := err.(*machines.UnconsumedInput); ui != nil && is {
			// e := err.(*machines.UnconsumedInput)
			// scanner.TC = ui.FailTC
			// // fmt.Println(SyntaxError(e.FailLine, e.StartColumn, e.FailColumn-e.StartColumn-1, string(text), "Tokenize Failed"))
			os.Exit(1)
		} else if err != nil {
			return err
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
			t.SourceCode = &text
			s.Tokens <- t
		}
	}

	close(s.Tokens)

	s.Done = true

	return nil
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
	s.Tokens = make(chan Token)
	s.lexer = lexer
	return s
}

// GetTokenName takes an id (an index) and returns the token's representation
func GetTokenName(id int) TokenType {
	return Tokens[id].token
}

// GetTokenId -
func GetTokenId(t TokenType) int {
	return tokmap[t]
}

// TokenIs -
func TokenIs(t Token, a TokenType) bool {
	return t.Type == GetTokenId(a)
}
