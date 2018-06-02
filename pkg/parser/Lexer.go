package parser

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const eof = -1

// stateFN - a function that takes a pointer to the
// lexer, and returns a new stateFN
type stateFN func(*Lexer) stateFN

// Lexer - stores the current parsing state
type Lexer struct {
	name   string     // The name of the file/module/whatnot (used for debugging and errors)
	input  string     // The input source code
	start  int        // the starting position of this token
	pos    int        // current position in the input
	width  int        // width of the last rune read
	tokens chan Token // channel of scanned tokens
	done   bool
}

func (l *Lexer) emit(t TokenType) {
	l.tokens <- Token{t, l.input[l.start:l.pos], l.start}
	l.start = l.pos
}

// Run the Lexer. Tokenize then Build AST
func (l *Lexer) run() {
	for state := lexProgram; state != nil; {
		state = state(l)
	}
	l.done = true
	close(l.tokens)
}

// NewLexer will generate a new Lexer to be run on some input
func NewLexer(name, input string) (*Lexer, chan Token) {
	l := &Lexer{
		name:   name,
		input:  input,
		tokens: make(chan Token),
		done:   false,
	}
	go l.run() // Concurrently run state machine.
	l.catchTokens()
	return l, l.tokens
}

func (l *Lexer) catchTokens() {
	for {
		if l.done {
			fmt.Println("DONE")
			break
		}
		t := <-l.tokens
		if t.Type() == tokenError {
			fmt.Printf(SyntaxError(l.start, l.pos-l.start, l.input, l.name, t.Value()))
			break
		}
		if t.Type() == tokenEOF {
			fmt.Println("EOF")
			break
		}

		fmt.Println(t.Type(), t)
	}
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFN {
	l.tokens <- Token{
		typ: tokenError,
		val: fmt.Sprintf(format, args...),
	}
	return nil
}

// next returns the next rune in the input.
func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	rn, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width

	return rn
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) peek() rune {
	rune := l.next()
	l.backup()
	return rune
}

func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}

func lexNumber(l *Lexer) stateFN {
	// Optional leading sign.
	l.accept("+-")
	// is it hex?
	digits := "0123456789"
	if l.accept("0") && l.accept("x") {
		digits = "0123456789abcdefABCDEF"
	}

	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}

	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}

	if isAlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}

	l.emit(tokenNumber)
	return lexProgram
}

func lexString(l *Lexer) stateFN {
	isEscaped := false

	for !(isEscaped || l.accept("\"")) {
		if isEscaped {
			isEscaped = false
		}
		r := l.next()
		if r == '\\' {
			isEscaped = true
		}

	}

	l.emit(tokenString)

	return lexProgram(l)
}

func lexProgram(l *Lexer) stateFN {
	for {
		switch r := l.next(); {
		case r == eof:
			return nil
		case unicode.IsSpace(r):
			l.ignore()
		case r == '+' || r == '-' || '0' <= r && r <= '9':
			return lexNumber(l)
		case r == '"':
			return lexString(l)
		default:
			l.ignore()
			return lexProgram(l)
		}
	}
}
