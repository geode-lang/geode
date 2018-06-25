package gen

import (
	"github.com/llir/llvm/ir/types"
	"gitlab.com/nickwanninger/geode/pkg/lexer"
	"gitlab.com/nickwanninger/geode/pkg/typesystem"
)

func (p *Parser) parseType() (t types.Type, ptr bool) {
	p.requires(lexer.TokType)
	t = typesystem.GlobalTypeMap.GetType(p.token.Value)
	p.next()
	for p.token.Is(lexer.TokOper) && p.token.Value == "*" {
		ptr = true
		t = types.NewPointer(t)
		p.next()
	}
	return t, ptr
}
