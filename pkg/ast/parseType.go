package ast

import (
	"github.com/llir/llvm/ir/types"
	"github.com/nickwanninger/geode/pkg/lexer"
	"github.com/nickwanninger/geode/pkg/typesystem"
)

func (p *Parser) parseType() (t types.Type, ptr bool) {
	p.requires(lexer.TokType)
	t = typesystem.GlobalTypeMap.GetType(p.token.Value)
	p.next()

	for {

		// if p.token.Is(lexer.TokLeftBrace) {
		// 	ptr = true
		// 	t = types.NewArray(t, 0)
		// 	p.next()
		// 	p.requires(lexer.TokRightBrace)
		// 	p.next()
		// 	continue
		// }

		if p.token.Is(lexer.TokOper) && p.token.Value == "*" {
			t = types.NewPointer(t)
			p.next()
			continue
		}

		break

	}
	return t, ptr
}
