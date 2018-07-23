package ast

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/geode-lang/geode/pkg/util/log"
)

// UnescapeString UTF-8 string
// e.g. convert "\u0e27\u0e23\u0e0d\u0e32" to "วรญา"
func UnescapeString(str string) string {
	var replaced []rune
	s := []byte(str)
	r := 0
	for r < len(s) {
		if s[r] != '\\' {
			c, size := utf8.DecodeRune(s[r:])
			r += size
			replaced = append(replaced, c)
			continue
		}
		r++
		if r >= len(s) {
			log.Fatal("Escape sequence at end of string.")
		}
		switch s[r] {
		default:
			log.Fatal("Expected valid escape code after \\, but got %q.", s[r])
		case 'b':
			replaced = append(replaced, rune(0x0008))
			r++
		case 't':
			replaced = append(replaced, rune(0x0009))
			r++
		case 'n':
			replaced = append(replaced, rune(0x000A))
			r++
		case 'f':
			replaced = append(replaced, rune(0x000C))
			r++
		case 'r':
			replaced = append(replaced, rune(0x000D))
			r++
		case '"':
			replaced = append(replaced)
			r++
		case '\\':
			replaced = append(replaced, rune(0x005C))
			r++
		case 'u':
			// At this point, we know we have a Unicode escape of the form
			// `uXXXX` at [r, r+5). (Because the lexer guarantees this
			// for us.)
			escaped := asciiEscapeToUnicode(s[r+1 : r+5])
			replaced = append(replaced, escaped)
			r += 5
		case 'U':
			// At this point, we know we have a Unicode escape of the form
			// `uXXXX` at [r, r+9). (Because the lexer guarantees this
			// for us.)
			escaped := asciiEscapeToUnicode(s[r+1 : r+9])
			replaced = append(replaced, escaped)
			r += 9
		}
	}
	return string(replaced)
}

func asciiEscapeToUnicode(bs []byte) rune {
	s := string(bs)
	hex, err := strconv.ParseUint(strings.ToLower(s), 16, 32)
	fmt.Println(hex, err)
	if err != nil {
		log.Fatal("Could not parse '%s' as a hexadecimal number, but the "+
			"lexer claims it's OK: %s\n", s, err)
	}
	if !utf8.ValidRune(rune(hex)) {
		log.Fatal("Escaped character '\\u%s' is not valid UTF-8.", s)
	}
	return rune(hex)
}

func (p *Parser) parseStringExpr() Node {
	n := StringNode{}
	n.TokenReference.Token = p.token
	n.NodeType = nodeString
	escaped, err := strconv.Unquote(p.token.Value)
	if err != nil {
		panic(err)
	}

	// untrimmed := p.token.Value

	// trimmed := untrimmed[1 : len(untrimmed)-1]
	// fmt.Println("untrimmed", untrimmed)
	// fmt.Println("trimmed", trimmed)

	// escaped := strings.Trim(UnescapeString(p.token.Value), "\"")
	// fmt.Println(p.token.Value, escaped)
	n.Value = escaped
	p.next()
	return n
}
