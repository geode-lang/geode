package ast

import (
	"encoding/json"
	"strconv"
)

// UnescapeString UTF-8 string
// e.g. convert "\u0e27\u0e23\u0e0d\u0e32" to "วรญา"
func UnescapeString(str string) []byte {
	bs := make([]byte, 0)
	json.Unmarshal([]byte(`"`+str+`"`), &bs)
	return bs
}

func (p *Parser) parseStringExpr() Node {
	n := StringNode{}
	n.NodeType = nodeString
	escaped, err := strconv.Unquote(p.token.Value)
	if err != nil {
		panic(err)
	}
	n.Value = string(escaped)
	p.next()
	return n
}
