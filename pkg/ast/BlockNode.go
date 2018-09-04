package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/geode-lang/llvm/ir/value"
)

// BlockNode is a block statement. A block statement is just an array of Nodes
// that run in sequence.
type BlockNode struct {
	NodeType
	TokenReference
	Nodes []Node
}

// NameString implements Node.NameString
func (n BlockNode) NameString() string { return "BlockNode" }

// Codegen implements Node.Codegen for BlockNode
func (n BlockNode) Codegen(prog *Program) (value.Value, error) {
	prog.Scope = prog.Scope.SpawnChild()

	for _, node := range n.Nodes {
		_, err := node.Codegen(prog)
		if err != nil {
			return nil, err
		}
	}
	if prog.Scope.Parent == nil {
		return nil, fmt.Errorf("attempt to step up a scope failed because the parent was nil %q", n.Token.FileInfo())
	}
	prog.Scope = prog.Scope.Parent
	return prog.Compiler.CurrentBlock(), nil
}

var blockindentdepth = 0

func (n BlockNode) String() string {

	buff := &bytes.Buffer{}

	fmt.Fprintf(buff, "{\n")

	blockindentdepth++

	for _, node := range n.Nodes {
		fmt.Fprintf(buff, "%s%s\n", strings.Repeat("\t", blockindentdepth), node)
	}

	blockindentdepth--

	fmt.Fprintf(buff, "%s}", strings.Repeat("\t", blockindentdepth))
	return buff.String()
}
