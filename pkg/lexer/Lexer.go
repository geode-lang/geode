package lexer

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/geode-lang/geode/pkg/info"
	"github.com/geode-lang/geode/pkg/util/log"
)

var tokenTypeOverrides = map[string]TokenType{
	"return":  TokReturn,
	"if":      TokIf,
	"else":    TokElse,
	"for":     TokFor,
	"while":   TokWhile,
	"func":    TokFuncDefn,
	"λ":       TokFuncDefn,
	"new":     TokNew,
	"class":   TokClassDefn,
	"include": TokDependency,
	"link":    TokDependency,
	"is":      TokNamespace,
	"sizeof":  TokSizeof,
	"as":      TokAs,
	"true":    TokBool,
	"false":   TokBool,
	"(":       TokLeftParen,
	")":       TokRightParen,
	"{":       TokLeftCurly,
	"}":       TokRightCurly,
	"[":       TokLeftBrace,
	"]":       TokRightBrace,
	"->":      TokRightArrow,
	"<-":      TokLeftArrow,
	"←":       TokLeftArrow,
	";":       TokSemiColon,
	":":       TokNamespaceAccess,
	":=":      TokAssignment,
	"...":     TokElipsis,
	".":       TokDot,

	"or": TokOper,
	"+=": TokCompoundAssignment,
	"-=": TokCompoundAssignment,
	"*=": TokCompoundAssignment,
	"/=": TokCompoundAssignment,
}

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*Lexer) stateFn

// Lexer - an internal rep of the lexer
type Lexer struct {
	source     *Sourcefile
	tokenCount int // the number of tokens lexed/emitted
	line       int
	col        int
	pos        int // current position in input
	start      int // beginning position of the current token
	width      int // width of last rune read from input
	input      string
	tokens     chan Token
}

// Lex - takes a string and turns it into tokens
func Lex(source *Sourcefile) chan Token {
	l := NewLexer()
	l.source = source
	l.input = source.String()
	go log.Timed(fmt.Sprintf("Lex %s", source.Path), l.run)
	return l.tokens
}

func (l *Lexer) run() {
	for state := lexTopLevel; state != nil; {
		state = state(l)
		if state == nil {
			break
		}
	}
	log.Verbose("Lexer emitted %d tokens from %s\n", l.tokenCount, l.source.Path)
	close(l.tokens) // No more tokens will be delivered.
}

// QuickLex takes a string and lexes it into a token array
func QuickLex(str string) []Token {
	source, _ := NewSourcefile("temp")
	source.LoadString(str)

	tokArr := make([]Token, 0)

	for t := range Lex(source) {
		// fmt.Println(t)
		tokArr = append(tokArr, t)
	}

	return tokArr
}

func (l *Lexer) value() string {
	return l.input[l.start:l.pos]
}
func (l *Lexer) emit(typ TokenType) {
	l.tokenCount++
	if typ != TokNoEmit {
		tok := Token{}
		tok.source = l.source
		tok.Value = l.value()
		tok.Pos = int(l.start)
		tok.EndPos = int(l.pos)
		tok.Line = l.line
		tok.Column = l.col

		newTyp, override := tokenTypeOverrides[tok.Value]
		if override {
			typ = newTyp
		}
		// inputRunes := []rune(l.input)

		// length := utf8.RuneCountInString(l.input)

		// fmt.Println("===")
		// fmt.Println("len", length)
		// if l.start > 1 {
		// 	fmt.Printf("%q %q\n", inputRunes[l.start-1], tok.Value)
		// 	tok.SpaceBefore = isSpace(inputRunes[l.start-1])
		// }

		tok.SpaceBefore = true

		// if l.pos < length {
		// 	fmt.Println("end", l.pos)
		// 	tok.SpaceAfter = isSpace(inputRunes[l.pos])
		// }

		tok.Type = typ
		info.AddToken(tok)
		l.tokens <- tok
	}
	l.start = l.pos
}

// l.next() returns eof to signal end of file to a stateFn.
const eof = -1

// next returns the next rune from the input and advances the scan.
// It returns the eof constant (-1) if the scanner is at the end of
// the input.
func (l *Lexer) next() rune {
	r, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width
	l.pos += l.width

	if l.width == 0 {
		return eof
	}
	if r == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col += width
	}
	return r
}

// peek returns the next rune without moving the scan forward.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup moves the scan back one rune.
func (l *Lexer) backup() {
	l.pos -= l.width
	r, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width
	if r == '\n' {
		l.line--
	}
}

// ignore skips the pending input before this point.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// acceptRun consumes a run of runes from valid set.
func (l *Lexer) acceptRun(valid string) {
	l.acceptRunPredicate(func(r rune) bool {
		return strings.IndexRune(valid, r) >= 0
	})
	l.backup()
}

func (l *Lexer) acceptRunPredicate(pred func(rune) bool) {
	for {
		if !pred(l.next()) {
			break
		}
	}
	l.backup()
}

func lexTopLevel(l *Lexer) stateFn {
	// Either whitespace, an empty line, a comment,
	// a number, a paren, identifier, or unary operator.
	r := l.next()

	switch {
	case r == eof:

		return nil
	case strings.IndexRune("-0123456789.", r) >= 0:
		l.backup()
		return lexNumber
	case isAlphaNumeric(r):
		l.backup()
		return lexIdentifer
	case r == '#':
		return lexComment
	case isSpace(r):
		l.backup()
		return lexSpace
	case isNewline(r):
		l.backup()
		return lexNewline
	case r == ';':
		l.emit(TokSemiColon)
		return lexTopLevel
	case r == ',':
		l.emit(TokComma)
		return lexTopLevel

	case r == '(':
		l.emit(TokLeftParen)
		return lexTopLevel
	case r == ')':
		l.emit(TokRightParen)
		return lexTopLevel

	case r == '{':
		l.emit(TokLeftCurly)
		return lexTopLevel
	case r == '}':
		l.emit(TokRightCurly)
		return lexTopLevel

	case r == '[':
		l.emit(TokLeftBrace)
		return lexTopLevel
	case r == ']':
		l.emit(TokLeftBrace)
		return lexTopLevel

	case isOperator(r):
		return lexSymbol
	case r == '"':
		// l.backup()
		return lexStringLiteral
	}
	return l.fatal("unrecognized character: %#U\n", r)
}

// fatal is just a statefn wrapper around log.Fatal
func (l *Lexer) fatal(format string, args ...interface{}) stateFn {
	log.Fatal(format, args...)
	return nil
}

func lexIdentifer(l *Lexer) stateFn {
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r) || r == '\'' || r == '-':
			// absorb
		default:
			l.backup()
			l.emit(TokIdent)
			return lexTopLevel
		}
	}
}

func lexNumber(l *Lexer) stateFn {
	l.acceptRun("-0123456789.xabcdefABCDEF")
	l.next()
	// There is a chance that the numeric expression lexer will
	// parse only a + or a - since it gets handled first in the list
	// so if it is only a minus,
	if isOnly(l.value(), '-') {
		l.backup()

		return lexSymbol
	}

	l.emit(TokNumber)
	return lexTopLevel
}

func lexComment(l *Lexer) stateFn {
	l.acceptRunPredicate(func(r rune) bool {
		// log.Verbose("Rune: %#U\n", r)
		return r != '\n' && r != -1
	})
	l.emit(TokComment)
	return lexTopLevel
}

// lexSpace globs contiguous whitespace and ignores them.
func lexSpace(l *Lexer) stateFn {
	l.acceptRunPredicate(isSpace)
	l.ignore()
	return lexTopLevel
}

// lexNewline globs contiguous newlines.
func lexNewline(l *Lexer) stateFn {
	l.acceptRunPredicate(isNewline)
	l.ignore()
	return lexTopLevel
}

func lexSymbol(l *Lexer) stateFn {
	var lastRune rune
	l.acceptRunPredicate(func(c rune) bool {
		if c == lastRune || l.peek() == lastRune {
			l.backup()
			l.emit(TokOper)
			l.next()
		}
		lastRune = c
		return isOperator(c)
	})
	l.emit(TokOper)
	return lexTopLevel
}

func lexStringLiteral(l *Lexer) stateFn {
	for {
		r := l.next()
		if r == eof {
			break
		}

		if r == '\\' {
			// Skip escape ('\' and next char)
			l.next()
		}
		if r == '"' {
			l.emit(TokString)
			return lexTopLevel
		}
	}
	return l.fatal("Unclosed string literal\n")
}

//
// Helper Functions
///

const operators = "&\\*+-/%:!=<>≤≥≠.←|&^"

func isOperator(r rune) bool {
	return strings.IndexRune(operators, r) >= 0
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}
func isNewline(r rune) bool {
	return r == '\n'
}

// isValidIdefRune reports if r may be part of an identifier name.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isOnly(s string, r rune) bool {
	for _, c := range s {
		if r != c {
			return false
		}
	}
	return true
}

// NewLexer produces a new lexer and poluates it with the configuration
func NewLexer() *Lexer {
	s := &Lexer{}
	s.line = 1
	s.col = 1
	s.tokens = make(chan Token)
	return s
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
					// m["type"] = t.Type.String()
					m["type_raw"] = t.Type
					m["value"] = t.Value
					m["start_pos"] = t.Pos
					m["end_pos"] = t.EndPos
					// _, m["type_inference"] = t.InferType()
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
