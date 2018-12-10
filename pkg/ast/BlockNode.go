package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/llir/llvm/ir/value"
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
	prog.ScopeDown(n.Token)

	for _, node := range n.Nodes {

		_, err := node.Codegen(prog)
		if err != nil {
			return nil, err
		}

		if _, isReturn := node.(ReturnNode); isReturn {
			break
		}
	}

	if err := prog.ScopeUp(); err != nil {
		return nil, err
	}
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
