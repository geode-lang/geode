package parser

// Alot of this based off this video:
// https://www.youtube.com/watch?v=HxaD_trXwRE

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
	name      string     // The name of the file/module/whatnot (used for debugging and errors)
	input     string     // The input source code
	start     int        // the starting position of this token
	pos       int        // current position in the input
	width     int        // width of the last rune read
	tokenchan chan Token // channel of scanned tokens
	tokens    []Token
	done      bool
}

func (l *Lexer) emit(t TokenType) {
	l.tokenchan <- Token{t, l.input[l.start:l.pos], l.start}
	l.start = l.pos
}

// Run the Lexer. Tokenize then Build AST
func (l *Lexer) run() {
	for state := lexProgram; state != nil; {
		state = state(l)
	}
	l.done = true
	close(l.tokenchan)
}

// NewLexer will generate a new Lexer to be run on some input
func NewLexer(name, input string) (*Lexer, chan Token) {
	l := &Lexer{
		name:      name,
		input:     input,
		tokenchan: make(chan Token),
		done:      false,
	}
	go l.run() // Concurrently run state machine.
	l.catchTokens()
	return l, l.tokenchan
}

func (l *Lexer) catchTokens() {
	for {
		if l.done {
			break
		}
		t := <-l.tokenchan

		if t.Type() == tokenNONE {
			continue
		}
		if t.Type() == tokenError {
			fmt.Printf(SyntaxError(l.start, l.pos-l.start, l.input, l.name, t.Value()))
			break
		}
		if t.Type() == tokenEOF {
			fmt.Println("EOF")
			break
		}

		l.tokens = append(l.tokens, t)
	}

	for _, t := range l.tokens {
		fmt.Println(t)
	}

}

func (l *Lexer) errorf(format string, args ...interface{}) stateFN {
	l.tokenchan <- Token{
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

func (l *Lexer) atTerminator() bool {
	r := l.peek()
	if unicode.IsSpace(r) || isEndOfLine(r) {
		return true
	}
	switch r {
	case eof, '.', ',', '|', ':', ')', '(':
		return true
	}
	return false
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
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

func lexIdentifier(l *Lexer) stateFN {
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb.
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			if !l.atTerminator() {
				return l.errorf("bad character %#U", r)
			}

			_, isKW := keywords[word]
			switch {
			case isKW:
				l.emit(keywords[word])
			case word[0] == '.':
				l.emit(tokenField)
			case word == "true", word == "false":
				l.emit(tokenBool)
			default:
				l.emit(tokenIdentifier)
			}
			break Loop
		}
	}
	return lexProgram
}

func lexOperator(l *Lexer) stateFN {
	var currentType TokenType
	var currentop string

	s := l.input[l.start:]

	if (s[0] == '-' || s[0] == '+') && l.accept("0123456789") {
		return lexNumber(l)
	}

	for k, t := range operators {
		if strings.HasPrefix(s, k) {
			currentop = k
			currentType = t
		}
	}
	width := utf8.RuneCountInString(currentop)
	l.pos += width

	l.emit(currentType)

	return lexProgram(l)
}

func isOperator(r rune) bool {
	for k := range operators {
		if []rune(k)[0] == r {
			return true
		}
	}
	return false
}

func lexProgram(l *Lexer) stateFN {
	for {
		switch r := l.next(); {
		case r == eof:
			return nil
		// We want to ignore all spaces.
		case unicode.IsSpace(r):
			l.ignore()

		case r == ',':
			l.emit(tokenComma)
			return lexProgram

		case isOperator(r):
			return lexOperator(l)

		case r == '+' || r == '-' || '0' <= r && r <= '9':
			return lexNumber(l)

		case r == '"':
			return lexString(l)

		case isAlphaNumeric(r):
			return lexIdentifier(l)

		// Parse the parenthesis, curly parens, and braces
		case groupers[r] != "tokenNone":
			l.emit(groupers[r])
			return lexProgram

		default:
			return l.errorf("Unknown Token")
		}
	}
}
