package ast

import (
	"fmt"

	"github.com/geode-lang/geode/pkg/lexer"
)

// TestNewParser -
func TestNewParser() {
	p := NewQuickParser(`foo.bar[0].zoop().zap[1]`)
	chain, err := p.parseCompoundExpression(false)
	if err != nil {
		fmt.Println("error in parsing:")
		fmt.Println(err)
	} else {
		fmt.Println(chain)
		fmt.Println()
		node, err := chain.(*BaseComponent).ConstructNode(nil)
		if err != nil {
			fmt.Println("error in construction:")
			fmt.Println(err)
		} else {
			fmt.Println("Node created:")
			fmt.Println(node)
		}

	}

}

func (p *Parser) parseExpression(allowdecl bool) Node {
	lhs := p.parseUnary(allowdecl)
	if lhs == nil {
		return nil
	}
	if p.token.Is(lexer.TokAs) {
		return p.parseCastExpr(lhs)
	}
	defer p.globTerminator()
	return p.parseBinaryOpRHS(-100, lhs)
}

// =========================== parseCompoundExpression ===========================

func (p *Parser) parseCompoundExpression(allowdecl bool) (ExpComponent, error) {
	var err error
	chain := &BaseComponent{}

	switch p.token.Type {

	case lexer.TokIdent:
		err = p.parseIdentifierComponent(chain, allowdecl)
	case lexer.TokNumber:
		err = p.parseNumberComponent(chain)
	case lexer.TokLeftBrace:
		err = p.parseArrayComponent(chain)
	case lexer.TokString:
		err = p.parseStringComponent(chain)
	case lexer.TokLeftParen:
		err = p.parseParenthesisComponent(chain)
	case lexer.TokBool:
		err = p.parseBooleanComponent(chain)
	case lexer.TokChar:
		err = p.parseCharComponent(chain)
	case lexer.TokInfo:
		err = p.parseTypeInfoComponent(chain)
	default:
		return nil, p.Errorf("Failed to parse expression: %s", p.token.FileInfo())
	}
	if err != nil {
		return nil, err
	}
	return chain, nil
}

// =========================== parseOperatorComponent ===========================

func (p *Parser) parseOperatorComponent(base *BaseComponent) error {
	var err error
	switch p.token.Type {
	case lexer.TokLeftParen:
		err = p.parseCallComponent(base)
	case lexer.TokLeftBrace:
		err = p.parseSubscriptComponent(base)
	case lexer.TokDot:
		err = p.parseDotComponent(base)
	}
	if err != nil {
		return nil
	}
	return err
}

// =========================== parseIdentifierComponent ===========================

// parseIdentifierComponent has two roles. The first of which is obvious, parsing plain
// identifiers. ex: "foo" where foo is a variable, function, etc... The other role of
// this funciton is to parse variable declarations, which it tries first by forking and
// backtracing if it fails.
func (p *Parser) parseIdentifierComponent(base *BaseComponent, allowdecl bool) error {

	fk := p.Fork()
	if allowdecl && fk.parseIdentDeclComponent(base) == nil {
		p.Join(fk)
		return nil
	}

	n := &IdentComponent{}
	n.token = p.token
	name, err := p.parseName()
	if err != nil {
		return err
	}
	n.Value = name
	base.Add(n)

	fork := p.Fork()
	err = fork.parseOperatorComponent(base)
	if err == nil {
		p.Join(fork)
	}

	return nil
}

// =========================== parseIdentDeclComponent ===========================

func (p *Parser) parseIdentDeclComponent(base *BaseComponent) error {

	n := &IdentDeclComponent{}
	n.token = p.token

	if !p.atType() {
		return p.Errorf("parser not at type")
	}

	n.Type = p.parseType()

	if !p.token.Is(lexer.TokIdent) {
		return p.Errorf("ident not found after type in declaration statement")
	}

	// we now know the token is an ident, so we pull the value from it.
	n.Name.Value = p.token.Value
	p.Next()

	base.Add(n)

	// Declarations have no rhs or anything after them because that is handled
	// by the binary operator parser

	return nil
}

// =========================== parseCallComponent ===========================

func (p *Parser) parseCallComponent(base *BaseComponent) error {
	n := &CallComponent{}
	n.token = p.token

	for p.Next(); p.token.Type != lexer.TokRightParen; {
		switch p.token.Type {
		case lexer.TokComma:
			p.Next()
		default:

			arg := p.parseExpression(false)

			if arg == nil {
				return p.Errorf("invalid call syntax")
			}
			n.Args = append(n.Args, arg)
		}
	}

	p.Next()

	base.Add(n)

	fork := p.Fork()
	err := fork.parseOperatorComponent(base)
	if err == nil {
		p.Join(fork)
	}

	return nil
}

// =========================== parseArrayComponent ===========================

func (p *Parser) parseArrayComponent(base *BaseComponent) error {
	n := &ArrayComponent{}
	n.token = p.token

	for p.Next(); p.token.Type != lexer.TokRightBrace; {
		switch p.token.Type {
		case lexer.TokComma:
			p.Next()
		default:

			val := p.parseExpression(false)

			if val == nil {
				return p.Errorf("invalid call syntax")
			}
			n.Values = append(n.Values, val)
		}
	}

	p.Next()

	base.Add(n)

	fork := p.Fork()
	err := fork.parseOperatorComponent(base)
	if err == nil {
		p.Join(fork)
	}
	return nil
}

// =========================== parseNumberComponent ===========================

func (p *Parser) parseNumberComponent(base *BaseComponent) error {
	n := &NumberComponent{}
	n.token = p.token

	n.Value = p.token.Value
	p.Next()
	base.Add(n)
	return nil
}

// =========================== parseSubscriptComponent ===========================

func (p *Parser) parseSubscriptComponent(base *BaseComponent) error {
	n := &SubscriptComponent{}
	var err error

	p.Next()

	n.Value, err = p.parseCompoundExpression(false)
	if err != nil {
		return err
	}

	if !p.token.Is(lexer.TokRightBrace) {
		return p.Errorf("malformed array subscript %s", p.token.FileInfo())
	}
	p.Next()
	base.Add(n)

	fork := p.Fork()
	err = fork.parseOperatorComponent(base)
	if err == nil {
		p.Join(fork)
	}

	return nil
}

// =========================== parseDotComponent ===========================

func (p *Parser) parseDotComponent(base *BaseComponent) error {
	n := &DotComponent{}
	n.token = p.token

	p.Next()
	n.Value = p.token.Value
	p.Next()
	base.Add(n)

	fork := p.Fork()
	err := fork.parseOperatorComponent(base)
	if err == nil {
		p.Join(fork)
	}

	return nil
}

// =========================== parseStringComponent ===========================

func (p *Parser) parseStringComponent(base *BaseComponent) error {
	n := &StringComponent{}
	n.token = p.token

	n.Value = p.token.Value
	p.Next()
	base.Add(n)

	fork := p.Fork()
	err := fork.parseOperatorComponent(base)
	if err == nil {
		p.Join(fork)
	}

	return nil
}

// =========================== parseParenthesisComponent ===========================

func (p *Parser) parseParenthesisComponent(base *BaseComponent) error {
	n := &ParenthesisComponent{}
	n.token = p.token

	if !p.token.Is(lexer.TokLeftParen) {
		return p.Errorf("parseParenthesisComponent expects a left paren to start")
	}

	p.Next()

	n.Value = p.parseExpression(false)

	if !p.token.Is(lexer.TokRightParen) {
		p.token.SyntaxError()
		return p.Errorf("invalid parenthesis syntax")
	}

	p.Next()

	base.Add(n)

	return nil
}

// =========================== parseBooleanComponent ===========================

func (p *Parser) parseBooleanComponent(base *BaseComponent) error {
	n := &BooleanComponent{}
	n.token = p.token

	if !p.token.Is(lexer.TokBool) {
		return p.Errorf("parseBooleanComponent expects a left paren to start")
	}

	n.Value = p.token.Value

	p.Next()

	base.Add(n)

	return nil
}

// =========================== parseCharComponent ===========================

func (p *Parser) parseCharComponent(base *BaseComponent) error {
	n := &CharComponent{}
	n.token = p.token

	if !p.token.Is(lexer.TokChar) {
		return p.Errorf("parseCharComponent expects a left paren to start")
	}

	n.Value = p.token.Value

	p.Next()

	base.Add(n)

	return nil
}

// =========================== parseTypeInfoComponent ===========================

func (p *Parser) parseTypeInfoComponent(base *BaseComponent) error {
	n := &TypeInfoComponent{}
	n.token = p.token

	p.Next()

	if !p.token.Is(lexer.TokLeftParen) {
		return p.Errorf("invalid call to info")
	}
	p.Next()

	n.Type = p.parseType()

	if !p.token.Is(lexer.TokRightParen) {
		return p.Errorf("invalid call to info")
	}

	p.Next()
	base.Add(n)

	fork := p.Fork()
	err := fork.parseOperatorComponent(base)
	if err == nil {
		p.Join(fork)
	}

	return nil
}
