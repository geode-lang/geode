package ast

import (
	"bytes"
	"encoding/hex"

	"github.com/geode-lang/geode/pkg/util/log"
)

const (
	hexChars = "0123456789abcdefABCDEF"
)

func isHex(r rune) bool {
	for _, c := range hexChars {
		if r == c {
			return true
		}
	}
	return false
}

// UnescapeString UTF-8 string
// e.g. convert "\u0e27\u0e23\u0e0d\u0e32" to "วรญา"
func UnescapeString(s string) (string, error) {
	// out := make([]rune, 0)
	buff := bytes.NewBufferString("")
	sr := []rune(s)

	escapes := map[rune]rune{
		'a':  0x07,
		'b':  0x08,
		'f':  0x0C,
		'n':  0x0A,
		'r':  0x0D,
		't':  0x09,
		'v':  0x0B,
		'\\': 0x5C,
		'\'': 0x27,
		'"':  0x22,
		'?':  0x3F,
	}

	for i := 0; i < len(sr); i++ {
		if sr[i] == '\\' {
			i++

			if sr[i] == 'x' {
				i++
				hexStr := make([]rune, 0)
				for ; i < len(sr) && isHex(sr[i]); i++ {
					hexStr = append(hexStr, sr[i])
				}

				bts, _ := hex.DecodeString(string(hexStr))
				for _, b := range bts {
					i--
					buff.WriteByte(b)
				}
				continue
			}

			esc, ok := escapes[sr[i]]
			if !ok {
				log.Fatal("Unknown escape: '\\%c'\n", sr[i])
			}
			buff.WriteRune(esc)
		} else {
			buff.WriteRune(sr[i])
		}
	}
	return buff.String(), nil
}

func (p *Parser) parseStringExpr() Node {
	n := StringNode{}
	n.TokenReference.Token = p.token
	n.NodeType = nodeString

	val := p.token.Value[1 : len(p.token.Value)-1]
	escaped, _ := UnescapeString(val)

	n.Value = escaped
	p.Next()
	return n
}

func (p *Parser) parseCharExpr() Node {
	n := CharNode{}

	n.TokenReference.Token = p.token
	n.NodeType = nodeChar

	val := p.token.Value[1 : len(p.token.Value)-1]

	escaped, _ := UnescapeString(val)
	n.Value = []rune(escaped)[0]
	p.Next()
	return n
}
