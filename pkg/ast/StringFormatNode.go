package ast

import (
	"bytes"
	"fmt"

	"github.com/geode-lang/geode/llvm/ir/value"
)

// StringFormatNode -
type StringFormatNode struct {
	NodeType
	TokenReference

	Format StringNode
	Args   []Node
}

// NameString implements Node.NameString
func (n StringFormatNode) NameString() string { return "StringFormatNode" }

// Codegen implements Node.Codegen for StringFormatNode
func (n StringFormatNode) Codegen(prog *Program) (value.Value, error) {
	str, err := n.Format.Codegen(prog)
	if err != nil {
		return nil, err
	}

	vals := make([]value.Value, 0, len(n.Args))

	vals = append(vals, str)

	for _, arg := range n.Args {
		val, err := arg.Codegen(prog)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}

	res, err := prog.NewRuntimeFunctionCall("__runtime_str_format", vals...)
	return res, err
}

// GenAccess implements Accessable.GenAccess
func (n StringFormatNode) GenAccess(prog *Program) (value.Value, error) {
	return n.Codegen(prog)
}

func (n StringFormatNode) String() string {
	buff := &bytes.Buffer{}

	fmt.Fprintf(buff, "%s(", n.Format)
	for i, arg := range n.Args {
		fmt.Fprintf(buff, "%s", arg)
		if i < len(n.Args)-1 {
			fmt.Fprintf(buff, ", ")
		}
	}

	fmt.Fprintf(buff, ")")
	return buff.String()
}
