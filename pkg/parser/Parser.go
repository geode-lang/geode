package parser

import (
	"fmt"
	"time"
)

// Parser - stores the current parsing state
type Parser struct {
	source   string
	filename string
	index    int
	time     time.Time
	tokens   []Token
	ast      *AST

	parens   TokenStack
	inSring  bool
	haserror bool
}

// Tokenize - takes the source and tokenize it into the struct's storage
func (p *Parser) Tokenize() error {
	fmt.Println("TOKENIZE")
	// Default, no error occurred
	return nil
}

// Run the parser. Tokenize then Build AST
func (p *Parser) Run() error {
	tokerr := p.Tokenize()
	if tokerr != nil {
		return tokerr
	}
	return nil
}

// New will generate a new parser state
// for the tokenizer and ast generator to use
func New(src string, filename string) *Parser {
	p := Parser{}
	p.source = src
	p.filename = filename
	p.time = time.Now()
	return &p
}
