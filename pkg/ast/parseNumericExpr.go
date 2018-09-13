package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/geode-lang/geode/llvm/ir/types"
)

func (p *Parser) parseNumericExpr() Node {

	n, _ := GetNumberNodeFromString(p.token.Value)
	p.Next()

	return n
}

func inferNumberType(str string) (types.Type, interface{}) {
	intval, intErr := strconv.ParseInt(str, 10, 64)
	if intErr == nil {
		return types.I64, intval
	}

	floatval, floatErr := strconv.ParseFloat(str, 64)
	if floatErr == nil {
		return types.Double, floatval
	}

	return nil, nil
}

// GetNumberNodeFromString returns the number node for a string
func GetNumberNodeFromString(str string) (Node, error) {
	t, val := inferNumberType(str)
	// Parse Hex Literals
	if strings.Contains(str, "x") {
		if !strings.Contains(str, "0x") {
			return nil, fmt.Errorf("hex Literal must be of the following format: 0x___")
		}
		n := IntNode{}
		n.NodeType = nodeInt
		parsed, e := strconv.ParseInt(strings.TrimPrefix(str, "0x"), 16, 64)
		if e != nil {
			return nil, fmt.Errorf("error decoding hex token")
		}
		n.Value = parsed
		return n, nil
	}

	// Parse Binary Literals
	if strings.Contains(str, "b") {
		if !strings.Contains(str, "0b") {
			fmt.Errorf("binary Literal must be of the following format: 0b___")
		} else {
			n := IntNode{}
			n.NodeType = nodeInt
			parsed, e := strconv.ParseInt(strings.TrimPrefix(str, "0b"), 2, 64)
			if e != nil {
				return nil, fmt.Errorf("error decoding binary token")
			}
			n.Value = parsed
			return n, nil
		}
	}

	if types.Equal(t, types.I64) {
		n := IntNode{}
		n.NodeType = nodeInt
		n.Value = val.(int64)
		return n, nil
	}

	if types.Equal(t, types.Double) {
		n := FloatNode{}
		n.NodeType = nodeFloat
		n.Value = val.(float64)
		return n, nil
	}

	if types.Equal(t, types.I8) {
		n := CharNode{}
		n.NodeType = nodeChar
		n.Value = val.(rune)
		return n, nil
	}
	return nil, fmt.Errorf("unable to parse number to node")
}
