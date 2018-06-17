package parser

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
	"os"
)

// int - The type of token as a string
// type TokenType int

// Assigning tokens integer values
const (
	TokError int = iota
	TokWhitespace
	TokChar
	TokString
	TokNumber

	TokOperatorStart
	TokStar
	TokPlus
	TokMinus
	TokDiv
	TokExp
	TokLT
	TokLTE
	TokGT
	TokGTE
	TokOperatorEnd

	TokSemiColon

	TokDefereference
	TokReference

	TokAssignment
	TokEquality

	TokRightParen
	TokLeftParen

	TokRightCurly
	TokLeftCurly

	TokRightBrace
	TokLeftBrace

	TokRightArrow
	TokLeftArrow

	TokFor
	TokWhile
	TokIf
	TokElse
	TokReturn
	TokFuncDefn

	TokType

	TokComma

	TokIdent

	TokComment
)

// TokenIsOperator will return if a given token is an operator or not
func TokenIsOperator(t int) bool {
	return t > TokOperatorStart && t < TokOperatorEnd
}

// TokenInfoRelation - allows a relationship between a token type and a certain regex
type TokenInfoRelation struct {
	token int
	regex string
	name  string
}

// Tokens is a list of all tokens
var Tokens = []TokenInfoRelation{
	{TokError, "", "error"},

	{TokChar, `'.'`, "char"},
	{TokString, `"([^\"]|(\\.))*"`, "string"},
	{TokNumber, `[+-]?[0-9]*\.?[0-9]+`, "number"},

	{TokStar, `\*`, "star"},
	{TokPlus, `\+`, "plus"},
	{TokMinus, `-`, "minus"},
	{TokDiv, `/`, "divide"},
	{TokExp, `\^`, "exp"},
	{TokSemiColon, `;`, "semicolon"},

	// {TokDefereference, `@`, "defereference"},
	// {TokReference, `\*`, "reference"},

	{TokAssignment, `:=`, "assign"},
	{TokEquality, `=`, "equality"},

	{TokRightParen, `\)`, "rightParen"},
	{TokLeftParen, `\(`, "leftParen"},

	{TokRightCurly, `}`, "rightCurly"},
	{TokLeftCurly, `{`, "leftCurly"},

	{TokRightBrace, `\[`, "rightBrace"},
	{TokLeftBrace, `\]`, "leftBrace"},

	{TokRightArrow, `->`, "rightArrow"},
	{TokLeftArrow, `<-`, "leftArrow"},
	// The main parser won't work correctly and will just look these up later
	{TokIf, "", "keywordIf"},
	{TokElse, "", "keywordElse"},
	{TokReturn, "", "keywordReturn"},
	{TokFuncDefn, "", "keywordFunc"},

	{TokType, "", "type"},

	{TokComma, `,`, "comma"},

	{TokIdent, `[a-zA-Z_][a-zA-Z0-9_]*`, "ident"},

	{TokComment, `\/\/[^\n]*`, "lineComment"},
	{TokComment, `{-.*-}`, "blockComment"},
	{TokWhitespace, `\s+`, "whitespace"},

	{TokLT, `<`, "lessThan"},
	{TokLTE, `<=|≤`, "lessThanOrEqual"},
	{TokGT, `>`, "greaterThan"},
	{TokGTE, `>=|≥`, "greaterThanOrEqual"},
}

var keyWordMap = map[string]int{
	"return": TokReturn,
	"if":     TokIf,
	"else":   TokElse,
	"for":    TokFor,
	"while":  TokWhile,
	"func":   TokFuncDefn,

	// We also determine type mapping in here as well
	"void":   TokType,
	"int":    TokType,
	"long":   TokType,
	"float":  TokType,
	"string": TokType,
	"char":   TokType,
}

var tokRegexMap map[string]int
var tokNameMap map[int]string

func init() {
	tokRegexMap = make(map[string]int)
	tokNameMap = make(map[int]string)
	for _, val := range Tokens {
		if val.regex != "" {
			tokRegexMap[val.regex] = val.token
		}
		tokNameMap[val.token] = val.name
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
			t.StartCol = to.StartColumn
			t.StartRow = to.StartLine
			t.EndCol = to.EndColumn
			t.EndRow = to.EndLine
			t.Type = to.Type
			t.Lexeme = to.Lexeme
			t.Value = string(to.Value.(string))
			fmt.Println(t)
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

	getToken := func(tokenType int) lexmachine.Action {
		return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			kw, isKwInMap := keyWordMap[string(m.Bytes)]
			if isKwInMap {
				return s.Token(kw, string(m.Bytes), m), nil
			}
			if tokenType == TokWhitespace {
				return nil, nil
			}
			return s.Token(tokenType, string(m.Bytes), m), nil
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
	go func() {
		for {
			// Read from the input channel of nodes.
			n, stillOpen := <-in
			// If the channel is closed, exit out of the printing phase
			if !stillOpen {
				close(out)
				return
			}
			spew.Dump(n)
			out <- n
		}
	}()
	return out
}

// GetTokenName returns the real name of a token
func GetTokenName(tok int) string {
	return tokNameMap[tok]
}
