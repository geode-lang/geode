package lexer

import (
	"encoding/json"
	"fmt"
	"os"

	"gitlab.com/nickwanninger/geode/pkg/typesystem"

	"github.com/davecgh/go-spew/spew"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

// TokenInfoRelation - allows a relationship between a token type and a certain regex
type TokenInfoRelation struct {
	token TokenType
	regex string
}

// Tokens is a list of all tokens
var Tokens = []TokenInfoRelation{
	{TokError, ""},

	{TokChar, `'.'`},
	{TokString, `"([^\"]|(\\.))*"`},
	{TokNumber, `[0-9]*\.?[0-9]+`},

	{TokElipsis, `\.\.\.`},

	{TokOper, `\*`},
	{TokOper, `\+`},
	{TokOper, `-`},
	{TokOper, `/`},
	{TokOper, `\^`},
	{TokOper, `!=`},
	{TokOper, `=`},
	{TokOper, `<`},
	{TokOper, `<=|≤`},
	{TokOper, `>`},
	{TokOper, `>=|≥`},
	{TokSemiColon, `;`},

	// {TokDefereference, `@`},
	// {TokReference, `\*`},

	{TokAssignment, `:=`},

	{TokRightParen, `\)`},
	{TokLeftParen, `\(`},

	{TokRightCurly, `}`},
	{TokLeftCurly, `{`},

	{TokRightBrace, `\[`},
	{TokLeftBrace, `\]`},

	{TokRightArrow, `->`},
	{TokLeftArrow, `<-`},
	// The main lexer won't work correctly and will just look these up later
	{TokIf, ""},
	{TokElse, ""},
	{TokReturn, ""},
	{TokFuncDefn, ""},

	{TokType, ""},

	{TokComma, `,`},

	{TokIdent, `[a-zA-Z_][a-zA-Z0-9_]*`},

	{TokComment, `\#[^\n]*`},
	{TokComment, `{-.*-}`},
	{TokWhitespace, `\s+`},
}

var keyWordMap = map[string]TokenType{
	"return": TokReturn,
	"if":     TokIf,
	"else":   TokElse,
	"for":    TokFor,
	"while":  TokWhile,
	"func":   TokFuncDefn,
}

var tokRegexMap map[string]TokenType

func init() {
	tokRegexMap = make(map[string]TokenType)
	for _, val := range Tokens {
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
	srcString := string(text)
	scanner, err := s.lexer.Scanner(text)
	if err != nil {
		return err
	}
	for tk, err, eof := scanner.Next(); !eof; tk, err, eof = scanner.Next() {
		if ui, is := err.(*machines.UnconsumedInput); ui != nil && is {
			e := err.(*machines.UnconsumedInput)
			spew.Dump(e)
			// scanner.TC = ui.FailTC
			// fmt.Println(SyntaxError(e.FailLine, e.StartColumn, e.FailColumn-e.StartColumn-1, string(text), "Tokenize Failed"))
			os.Exit(1)
		} else if err != nil {
			return err
		} else {

			// I don't like lexmachine's token, so I will convert it to my own
			to := *tk.(*lexmachine.Token)

			t := Token{}
			t.SourceCode = &srcString
			t.Pos = to.TC
			t.buildEndPos(to.EndColumn, to.EndLine)
			t.Type = TokenType(to.Type)
			t.Value = string(to.Value.(string))
			// t.SourceCode = &text
			s.Tokens <- t
		}
	}

	close(s.Tokens)

	s.Done = true

	return nil
}

// NewLexer produces a new lexer and poluates it with the configuration
func NewLexer() *LexState {

	getToken := func(tokenType TokenType) lexmachine.Action {
		return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			if typesystem.GlobalTypeMap.GetType(string(m.Bytes)) != nil {
				return s.Token(int(TokType), string(m.Bytes), m), nil
			}

			kw, isKwInMap := keyWordMap[string(m.Bytes)]
			if isKwInMap {
				return s.Token(int(kw), string(m.Bytes), m), nil
			}
			if tokenType == TokWhitespace {
				return nil, nil
			}
			return s.Token(int(tokenType), string(m.Bytes), m), nil
		}
	}

	lexer := lexmachine.NewLexer()

	for k, v := range tokRegexMap {
		lexer.Add([]byte(k), getToken(v))
	}
	s := &LexState{}
	s.Tokens = make(chan Token)
	s.lexer = lexer
	return s
}

// GetTokenID -
func GetTokenID(t int) int {
	return t
}

// DumpTokens takes a channel of tokens and prints all tokens it recieves,
// then pushes them back out a new channel it makes and returns
func DumpTokens(in chan Token) chan Token {
	out := make(chan Token)
	tokens := make([]Token, 0)
	go func() {
		for {
			// Read from the input channel of nodes.
			n, stillOpen := <-in
			// If the channel is closed, exit out of the printing phase
			if !stillOpen {
				tokenMaps := make([]map[string]interface{}, 0)
				for _, t := range tokens {
					m := make(map[string]interface{})
					m["type"] = t.Type.String()
					m["type_raw"] = t.Type
					m["value"] = t.Value
					m["start_pos"] = t.Pos
					m["end_pos"] = t.EndPos
					_, m["type_inference"] = t.InferType()
					tokenMaps = append(tokenMaps, m)
				}
				j, _ := json.MarshalIndent(tokenMaps, "", "   ")
				fmt.Println(string(j))
				close(out)
				return
			}
			tokens = append(tokens, n)
			out <- n
		}
	}()
	return out
}
