package ast

import (
	"fmt"
	"strings"

	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/metadata"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// GeodeBinaryInstr implements ir.Instruction
// This will be easier than doing individual work
// per operation
type GeodeBinaryInstr struct {
	// Parent basic block.
	Parent *ir.BasicBlock
	// Name of the local variable associated with the instruction.
	Name string
	// Operator string.
	Operator string
	// Operands.
	X, Y value.Value
	// Map from metadata identifier (e.g. !dbg) to metadata associated with the
	// instruction.
	Metadata map[string]*metadata.Metadata
}

// NewGeodeBinaryInstr returns a new geode binary instruction
func NewGeodeBinaryInstr(op string, x, y value.Value) *GeodeBinaryInstr {
	return &GeodeBinaryInstr{
		X:        x,
		Y:        y,
		Operator: op,
		Metadata: make(map[string]*metadata.Metadata),
	}
}

func escapeIdent(s string) string {
	tail := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz$-._0123456789"
	// Check if a replacement is required.
	extra := 0
	for i := 0; i < len(s); i++ {
		if strings.IndexByte(tail, s[i]) == -1 {
			// Two extra bytes are required for each invalid byte; e.g.
			//    "#" -> `\23`
			//    "世" -> `\E4\B8\96`
			extra += 2
		}
	}
	if extra == 0 {
		return s
	}

	// Replace invalid characters.
	const hextable = "0123456789ABCDEF"
	buf := make([]byte, len(s)+extra)
	j := 0
	for i := 0; i < len(s); i++ {
		b := s[i]
		if strings.IndexByte(tail, b) != -1 {
			buf[j] = b
			j++
			continue
		}
		buf[j] = '\\'
		buf[j+1] = hextable[b>>4]
		buf[j+2] = hextable[b&0x0F]
		j += 3
	}
	// Add surrounding quotes.
	return `"` + string(buf) + `"`
}

// Type returns the type of the instruction.
func (inst *GeodeBinaryInstr) Type() types.Type {
	return inst.X.Type()
}

// Ident returns the identifier associated with the instruction.
func (inst *GeodeBinaryInstr) Ident() string {
	return "%" + escapeIdent(inst.Name)
}

// GetName returns the name of the local variable associated with the
// instruction.
func (inst *GeodeBinaryInstr) GetName() string {
	return inst.Name
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *GeodeBinaryInstr) SetName(name string) {
	inst.Name = name
}

// String returns the LLVM syntax representation of the instruction.
func (inst *GeodeBinaryInstr) String() string {
	return fmt.Sprintf("%s = %s %s %s, %s",
		inst.Ident(),
		inst.Operator,
		inst.Type(),
		inst.X.Ident(),
		inst.Y.Ident())
}

// GetParent returns the parent basic block of the instruction.
func (inst *GeodeBinaryInstr) GetParent() *ir.BasicBlock {
	return inst.Parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *GeodeBinaryInstr) SetParent(parent *ir.BasicBlock) {
	inst.Parent = parent
}
