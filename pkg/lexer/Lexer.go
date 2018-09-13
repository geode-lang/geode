package lexer

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
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
	"let":     TokLet,
	"class":   TokClassDefn,
	"include": TokDependency,
	"link":    TokDependency,
	"is":      TokNamespace,
	"info":    TokInfo,
	"as":      TokAs,
	"true":    TokBool,
	"false":   TokBool,
	"nil":     TokNil,
	"(":       TokLeftParen,
	")":       TokRightParen,
	"{":       TokLeftCurly,
	"}":       TokRightCurly,
	"[":       TokLeftBrace,
	"]":       TokRightBrace,
	"->":      TokRightArrow,
	";":       TokSemiColon,
	":":       TokNamespaceAccess,
	"...":     TokElipsis,
	".":       TokDot,
	"?":       TokQuestionMark,

	"<-": TokOper,
	":=": TokOper,
	"+=": TokOper,
	"-=": TokOper,
	"*=": TokOper,
	"/=": TokOper,
}

var tokenAliasOverrides = map[string]string{
	"≠": "!=",
	"≤": "<=",
	"≥": ">=",
	"¬": "!",
	"⊕": "^",
	"∨": "||",
	"∧": "&&",
	"λ": "func",
	"←": "<-",
}

func getTokenValueAlias(value string) string {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	if alias, exists := tokenAliasOverrides[value]; exists {
		return alias
	}

	return value
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
	tokens     []Token
}

// Lex - takes a string and turns it into tokens
func Lex(source *Sourcefile) []Token {
	l := NewLexer()
	l.source = source
	l.input = source.String()
	log.Timed(fmt.Sprintf("Lex %s", source.Path), l.run)
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
}

// QuickLex takes a string and lexes it into a token array
func QuickLex(str string) []Token {
	source, _ := NewSourcefile("temp")
	source.LoadString(str)

	tokArr := Lex(source)

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

		tok.Value = getTokenValueAlias(tok.Value)

		// aliasenAliasOverrides[tok.Value]
		// if hasAlias {
		// 	tok.Value = alias
		// }

		tok.Pos = int(l.start)
		tok.EndPos = int(l.pos)
		tok.Line = l.line
		tok.Column = l.col

		newTyp, override := tokenTypeOverrides[tok.Value]
		if override {
			typ = newTyp
		}

		tok.SpaceBefore = true

		tok.Type = typ

		info.AddToken(tok)

		l.tokens = append(l.tokens, tok)
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
	case strings.IndexRune("-0123456789", r) >= 0:
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
		l.backup()
		return lexSymbol
	case r == '"':
		// l.backup()
		return lexStringLiteral
	case r == '\'':
		// l.backup()
		return lexCharLiteral
	}
	return l.fatal("unrecognized character: %#U\n", r)
}

// fatal is just a statefn wrapper around log.Fatal
func (l *Lexer) fatal(format string, args ...interface{}) stateFn {
	log.Fatal(format, args...)
	return nil
}

func lexIdentifer(l *Lexer) stateFn {
	sColonCount := 0
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r) || r == '\'' || r == ':':
			if r == ':' {
				sColonCount++
			}
			if sColonCount > 1 {
				l.emit(TokError)
			}
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

	// If the lexer's state ever results in an of these runs,
	// the lexer will emit that value. the finalRuns are a list
	// of the maxiumum repeats of tokens. They will be a list
	// of valid tokens in the language as repeats that aren't in
	// this list must be invalid
	finalRuns := map[string]bool{
		"...": true,
		"*":   true,
		"*=":  true,
		"&&":  true,
	}

	l.acceptRunPredicate(func(c rune) bool {
		if finalRuns[l.value()] && finalRuns[l.value()] {
			l.emit(TokOper)
		}
		return isOperator(c)
	})

	if len(l.value()) > 0 {
		l.emit(TokOper)
	}

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

func lexCharLiteral(l *Lexer) stateFn {
	for {
		r := l.next()
		if r == eof {
			break
		}

		if r == '\\' {
			// Skip escape ('\' and next char)
			l.next()
		}
		if r == '\'' {
			l.emit(TokChar)
			return lexTopLevel
		}
	}
	return l.fatal("Unclosed char literal\n")
}

//
// Helper Functions
//

const operators = "&\\*+-/%:!=<>≤≥≠.←|&^?"

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
	s.tokens = make([]Token, 0)
	return s
}
