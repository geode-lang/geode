package ast

import (
	"fmt"

	"github.com/geode-lang/llvm/ir/value"
)

// GenericSymbol is a way to represent a generic in a function/struct/class definition
type GenericSymbol struct {
	Name string
}

func (g *GenericSymbol) String() string {
	return g.Name
}

// GenericNodeWrapper -
type GenericNodeWrapper struct {
	NodeType
	TokenReference

	generics []GenericSymbol
	node     Node
}

func (n GenericNodeWrapper) String() string {
	return fmt.Sprintf("%s<%s>", n.node.NameString(), n.generics)
}

// Codegen implements Node.Codegen
func (n GenericNodeWrapper) Codegen(scope *Scope, c *Compiler) value.Value {
	return nil
}
