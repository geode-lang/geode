package ast

import (
	"fmt"
	"strings"

	"github.com/geode-lang/geode/pkg/lexer"
)

type namestack struct {
	parts []string
}

func (n *namestack) push(part string) {
	n.parts = append(n.parts, part)
}

func (n *namestack) last() string {
	return n.parts[len(n.parts)-1]
}

func (n *namestack) join() string {
	return strings.Join(n.parts, "")
}

func (n *namestack) String() string {
	return fmt.Sprintf("[%s]", strings.Join(n.parts, ", "))
}

func (p *Parser) parseName() (string, error) {
	stack := &namestack{}

	if !p.token.Is(lexer.TokIdent, lexer.TokType) {
		return "", fmt.Errorf("Invalid Name Reference")
	}

	for {
		if p.token.Is(lexer.TokIdent, lexer.TokType) {
			stack.push(p.token.Value)
		} else {
			return "", fmt.Errorf("Invalid Name Reference")
		}
		p.Next()
		break
	}

	return stack.join(), nil
}
