// generated by gen.go using 'go generate'; DO NOT EDIT.

// === [ Conversion instructions ] =============================================
//
// References:
//    http://llvm.org/docs/LangRef.html#conversion-operations

package ir

import (
	"fmt"

	"github.com/geode-lang/geode/llvm/enc"
	"github.com/geode-lang/geode/llvm/ir/metadata"
	"github.com/geode-lang/geode/llvm/ir/types"
	"github.com/geode-lang/geode/llvm/ir/value"
)

// --- [ trunc ] ---------------------------------------------------------------

// InstTrunc represents a truncation instruction.
//
// References:
//    http://llvm.org/docs/LangRef.html#trunc-instruction
type InstTrunc struct {
	// Parent basic block.
	Parent *BasicBlock
	// Name of the local variable associated with the instruction.
	Name string
	// Value before conversion.
	From value.Value
	// Type after conversion.
	To types.Type
	// Map from metadata identifier (e.g. !dbg) to metadata associated with the
	// instruction.
	Metadata map[string]*metadata.Metadata
}

// NewTrunc returns a new trunc instruction based on the given source value and target type.
func NewTrunc(from value.Value, to types.Type) *InstTrunc {
	return &InstTrunc{
		From:     from,
		To:       to,
		Metadata: make(map[string]*metadata.Metadata),
	}
}

// Type returns the type of the instruction.
func (inst *InstTrunc) Type() types.Type {
	return inst.To
}

// Ident returns the identifier associated with the instruction.
func (inst *InstTrunc) Ident() string {
	return enc.Local(inst.Name)
}

// GetName returns the name of the local variable associated with the
// instruction.
func (inst *InstTrunc) GetName() string {
	return inst.Name
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *InstTrunc) SetName(name string) {
	inst.Name = name
}

// String returns the LLVM syntax representation of the instruction.
func (inst *InstTrunc) String() string {
	md := metadataString(inst.Metadata, ",")
	return fmt.Sprintf("%s = trunc %s %s to %s%s",
		inst.Ident(),
		inst.From.Type(),
		inst.From.Ident(),
		inst.To,
		md)
}

// GetParent returns the parent basic block of the instruction.
func (inst *InstTrunc) GetParent() *BasicBlock {
	return inst.Parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *InstTrunc) SetParent(parent *BasicBlock) {
	inst.Parent = parent
}

// --- [ zext ] ----------------------------------------------------------------

// InstZExt represents a zero extension instruction.
//
// References:
//    http://llvm.org/docs/LangRef.html#zext-instruction
type InstZExt struct {
	// Parent basic block.
	Parent *BasicBlock
	// Name of the local variable associated with the instruction.
	Name string
	// Value before conversion.
	From value.Value
	// Type after conversion.
	To types.Type
	// Map from metadata identifier (e.g. !dbg) to metadata associated with the
	// instruction.
	Metadata map[string]*metadata.Metadata
}

// NewZExt returns a new zext instruction based on the given source value and target type.
func NewZExt(from value.Value, to types.Type) *InstZExt {
	return &InstZExt{
		From:     from,
		To:       to,
		Metadata: make(map[string]*metadata.Metadata),
	}
}

// Type returns the type of the instruction.
func (inst *InstZExt) Type() types.Type {
	return inst.To
}

// Ident returns the identifier associated with the instruction.
func (inst *InstZExt) Ident() string {
	return enc.Local(inst.Name)
}

// GetName returns the name of the local variable associated with the
// instruction.
func (inst *InstZExt) GetName() string {
	return inst.Name
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *InstZExt) SetName(name string) {
	inst.Name = name
}

// String returns the LLVM syntax representation of the instruction.
func (inst *InstZExt) String() string {
	md := metadataString(inst.Metadata, ",")
	return fmt.Sprintf("%s = zext %s %s to %s%s",
		inst.Ident(),
		inst.From.Type(),
		inst.From.Ident(),
		inst.To,
		md)
}

// GetParent returns the parent basic block of the instruction.
func (inst *InstZExt) GetParent() *BasicBlock {
	return inst.Parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *InstZExt) SetParent(parent *BasicBlock) {
	inst.Parent = parent
}

// --- [ sext ] ----------------------------------------------------------------

// InstSExt represents a sign extension instruction.
//
// References:
//    http://llvm.org/docs/LangRef.html#sext-instruction
type InstSExt struct {
	// Parent basic block.
	Parent *BasicBlock
	// Name of the local variable associated with the instruction.
	Name string
	// Value before conversion.
	From value.Value
	// Type after conversion.
	To types.Type
	// Map from metadata identifier (e.g. !dbg) to metadata associated with the
	// instruction.
	Metadata map[string]*metadata.Metadata
}

// NewSExt returns a new sext instruction based on the given source value and target type.
func NewSExt(from value.Value, to types.Type) *InstSExt {
	return &InstSExt{
		From:     from,
		To:       to,
		Metadata: make(map[string]*metadata.Metadata),
	}
}

// Type returns the type of the instruction.
func (inst *InstSExt) Type() types.Type {
	return inst.To
}

// Ident returns the identifier associated with the instruction.
func (inst *InstSExt) Ident() string {
	return enc.Local(inst.Name)
}

// GetName returns the name of the local variable associated with the
// instruction.
func (inst *InstSExt) GetName() string {
	return inst.Name
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *InstSExt) SetName(name string) {
	inst.Name = name
}

// String returns the LLVM syntax representation of the instruction.
func (inst *InstSExt) String() string {
	md := metadataString(inst.Metadata, ",")
	return fmt.Sprintf("%s = sext %s %s to %s%s",
		inst.Ident(),
		inst.From.Type(),
		inst.From.Ident(),
		inst.To,
		md)
}

// GetParent returns the parent basic block of the instruction.
func (inst *InstSExt) GetParent() *BasicBlock {
	return inst.Parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *InstSExt) SetParent(parent *BasicBlock) {
	inst.Parent = parent
}

// --- [ fptrunc ] -------------------------------------------------------------

// InstFPTrunc represents a floating-point truncation instruction.
//
// References:
//    http://llvm.org/docs/LangRef.html#fptrunc-instruction
type InstFPTrunc struct {
	// Parent basic block.
	Parent *BasicBlock
	// Name of the local variable associated with the instruction.
	Name string
	// Value before conversion.
	From value.Value
	// Type after conversion.
	To types.Type
	// Map from metadata identifier (e.g. !dbg) to metadata associated with the
	// instruction.
	Metadata map[string]*metadata.Metadata
}

// NewFPTrunc returns a new fptrunc instruction based on the given source value and target type.
func NewFPTrunc(from value.Value, to types.Type) *InstFPTrunc {
	return &InstFPTrunc{
		From:     from,
		To:       to,
		Metadata: make(map[string]*metadata.Metadata),
	}
}

// Type returns the type of the instruction.
func (inst *InstFPTrunc) Type() types.Type {
	return inst.To
}

// Ident returns the identifier associated with the instruction.
func (inst *InstFPTrunc) Ident() string {
	return enc.Local(inst.Name)
}

// GetName returns the name of the local variable associated with the
// instruction.
func (inst *InstFPTrunc) GetName() string {
	return inst.Name
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *InstFPTrunc) SetName(name string) {
	inst.Name = name
}

// String returns the LLVM syntax representation of the instruction.
func (inst *InstFPTrunc) String() string {
	md := metadataString(inst.Metadata, ",")
	return fmt.Sprintf("%s = fptrunc %s %s to %s%s",
		inst.Ident(),
		inst.From.Type(),
		inst.From.Ident(),
		inst.To,
		md)
}

// GetParent returns the parent basic block of the instruction.
func (inst *InstFPTrunc) GetParent() *BasicBlock {
	return inst.Parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *InstFPTrunc) SetParent(parent *BasicBlock) {
	inst.Parent = parent
}

// --- [ fpext ] ---------------------------------------------------------------

// InstFPExt represents a floating-point extension instruction.
//
// References:
//    http://llvm.org/docs/LangRef.html#fpext-instruction
type InstFPExt struct {
	// Parent basic block.
	Parent *BasicBlock
	// Name of the local variable associated with the instruction.
	Name string
	// Value before conversion.
	From value.Value
	// Type after conversion.
	To types.Type
	// Map from metadata identifier (e.g. !dbg) to metadata associated with the
	// instruction.
	Metadata map[string]*metadata.Metadata
}

// NewFPExt returns a new fpext instruction based on the given source value and target type.
func NewFPExt(from value.Value, to types.Type) *InstFPExt {
	return &InstFPExt{
		From:     from,
		To:       to,
		Metadata: make(map[string]*metadata.Metadata),
	}
}

// Type returns the type of the instruction.
func (inst *InstFPExt) Type() types.Type {
	return inst.To
}

// Ident returns the identifier associated with the instruction.
func (inst *InstFPExt) Ident() string {
	return enc.Local(inst.Name)
}

// GetName returns the name of the local variable associated with the
// instruction.
func (inst *InstFPExt) GetName() string {
	return inst.Name
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *InstFPExt) SetName(name string) {
	inst.Name = name
}

// String returns the LLVM syntax representation of the instruction.
func (inst *InstFPExt) String() string {
	md := metadataString(inst.Metadata, ",")
	return fmt.Sprintf("%s = fpext %s %s to %s%s",
		inst.Ident(),
		inst.From.Type(),
		inst.From.Ident(),
		inst.To,
		md)
}

// GetParent returns the parent basic block of the instruction.
func (inst *InstFPExt) GetParent() *BasicBlock {
	return inst.Parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *InstFPExt) SetParent(parent *BasicBlock) {
	inst.Parent = parent
}

// --- [ fptoui ] --------------------------------------------------------------

// InstFPToUI represents a floating-point to unsigned integer conversion instruction.
//
// References:
//    http://llvm.org/docs/LangRef.html#fptoui-instruction
type InstFPToUI struct {
	// Parent basic block.
	Parent *BasicBlock
	// Name of the local variable associated with the instruction.
	Name string
	// Value before conversion.
	From value.Value
	// Type after conversion.
	To types.Type
	// Map from metadata identifier (e.g. !dbg) to metadata associated with the
	// instruction.
	Metadata map[string]*metadata.Metadata
}

// NewFPToUI returns a new fptoui instruction based on the given source value and target type.
func NewFPToUI(from value.Value, to types.Type) *InstFPToUI {
	return &InstFPToUI{
		From:     from,
		To:       to,
		Metadata: make(map[string]*metadata.Metadata),
	}
}

// Type returns the type of the instruction.
func (inst *InstFPToUI) Type() types.Type {
	return inst.To
}

// Ident returns the identifier associated with the instruction.
func (inst *InstFPToUI) Ident() string {
	return enc.Local(inst.Name)
}

// GetName returns the name of the local variable associated with the
// instruction.
func (inst *InstFPToUI) GetName() string {
	return inst.Name
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *InstFPToUI) SetName(name string) {
	inst.Name = name
}

// String returns the LLVM syntax representation of the instruction.
func (inst *InstFPToUI) String() string {
	md := metadataString(inst.Metadata, ",")
	return fmt.Sprintf("%s = fptoui %s %s to %s%s",
		inst.Ident(),
		inst.From.Type(),
		inst.From.Ident(),
		inst.To,
		md)
}

// GetParent returns the parent basic block of the instruction.
func (inst *InstFPToUI) GetParent() *BasicBlock {
	return inst.Parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *InstFPToUI) SetParent(parent *BasicBlock) {
	inst.Parent = parent
}

// --- [ fptosi ] --------------------------------------------------------------

// InstFPToSI represents a floating-point to signed integer conversion instruction.
//
// References:
//    http://llvm.org/docs/LangRef.html#fptosi-instruction
type InstFPToSI struct {
	// Parent basic block.
	Parent *BasicBlock
	// Name of the local variable associated with the instruction.
	Name string
	// Value before conversion.
	From value.Value
	// Type after conversion.
	To types.Type
	// Map from metadata identifier (e.g. !dbg) to metadata associated with the
	// instruction.
	Metadata map[string]*metadata.Metadata
}

// NewFPToSI returns a new fptosi instruction based on the given source value and target type.
func NewFPToSI(from value.Value, to types.Type) *InstFPToSI {
	return &InstFPToSI{
		From:     from,
		To:       to,
		Metadata: make(map[string]*metadata.Metadata),
	}
}

// Type returns the type of the instruction.
func (inst *InstFPToSI) Type() types.Type {
	return inst.To
}

// Ident returns the identifier associated with the instruction.
func (inst *InstFPToSI) Ident() string {
	return enc.Local(inst.Name)
}

// GetName returns the name of the local variable associated with the
// instruction.
func (inst *InstFPToSI) GetName() string {
	return inst.Name
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *InstFPToSI) SetName(name string) {
	inst.Name = name
}

// String returns the LLVM syntax representation of the instruction.
func (inst *InstFPToSI) String() string {
	md := metadataString(inst.Metadata, ",")
	return fmt.Sprintf("%s = fptosi %s %s to %s%s",
		inst.Ident(),
		inst.From.Type(),
		inst.From.Ident(),
		inst.To,
		md)
}

// GetParent returns the parent basic block of the instruction.
func (inst *InstFPToSI) GetParent() *BasicBlock {
	return inst.Parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *InstFPToSI) SetParent(parent *BasicBlock) {
	inst.Parent = parent
}

// --- [ uitofp ] --------------------------------------------------------------

// InstUIToFP represents an unsigned integer to floating-point conversion instruction.
//
// References:
//    http://llvm.org/docs/LangRef.html#uitofp-instruction
type InstUIToFP struct {
	// Parent basic block.
	Parent *BasicBlock
	// Name of the local variable associated with the instruction.
	Name string
	// Value before conversion.
	From value.Value
	// Type after conversion.
	To types.Type
	// Map from metadata identifier (e.g. !dbg) to metadata associated with the
	// instruction.
	Metadata map[string]*metadata.Metadata
}

// NewUIToFP returns a new uitofp instruction based on the given source value and target type.
func NewUIToFP(from value.Value, to types.Type) *InstUIToFP {
	return &InstUIToFP{
		From:     from,
		To:       to,
		Metadata: make(map[string]*metadata.Metadata),
	}
}

// Type returns the type of the instruction.
func (inst *InstUIToFP) Type() types.Type {
	return inst.To
}

// Ident returns the identifier associated with the instruction.
func (inst *InstUIToFP) Ident() string {
	return enc.Local(inst.Name)
}

// GetName returns the name of the local variable associated with the
// instruction.
func (inst *InstUIToFP) GetName() string {
	return inst.Name
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *InstUIToFP) SetName(name string) {
	inst.Name = name
}

// String returns the LLVM syntax representation of the instruction.
func (inst *InstUIToFP) String() string {
	md := metadataString(inst.Metadata, ",")
	return fmt.Sprintf("%s = uitofp %s %s to %s%s",
		inst.Ident(),
		inst.From.Type(),
		inst.From.Ident(),
		inst.To,
		md)
}

// GetParent returns the parent basic block of the instruction.
func (inst *InstUIToFP) GetParent() *BasicBlock {
	return inst.Parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *InstUIToFP) SetParent(parent *BasicBlock) {
	inst.Parent = parent
}

// --- [ sitofp ] --------------------------------------------------------------

// InstSIToFP represents a signed integer to floating-point conversion instruction.
//
// References:
//    http://llvm.org/docs/LangRef.html#sitofp-instruction
type InstSIToFP struct {
	// Parent basic block.
	Parent *BasicBlock
	// Name of the local variable associated with the instruction.
	Name string
	// Value before conversion.
	From value.Value
	// Type after conversion.
	To types.Type
	// Map from metadata identifier (e.g. !dbg) to metadata associated with the
	// instruction.
	Metadata map[string]*metadata.Metadata
}

// NewSIToFP returns a new sitofp instruction based on the given source value and target type.
func NewSIToFP(from value.Value, to types.Type) *InstSIToFP {
	return &InstSIToFP{
		From:     from,
		To:       to,
		Metadata: make(map[string]*metadata.Metadata),
	}
}

// Type returns the type of the instruction.
func (inst *InstSIToFP) Type() types.Type {
	return inst.To
}

// Ident returns the identifier associated with the instruction.
func (inst *InstSIToFP) Ident() string {
	return enc.Local(inst.Name)
}

// GetName returns the name of the local variable associated with the
// instruction.
func (inst *InstSIToFP) GetName() string {
	return inst.Name
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *InstSIToFP) SetName(name string) {
	inst.Name = name
}

// String returns the LLVM syntax representation of the instruction.
func (inst *InstSIToFP) String() string {
	md := metadataString(inst.Metadata, ",")
	return fmt.Sprintf("%s = sitofp %s %s to %s%s",
		inst.Ident(),
		inst.From.Type(),
		inst.From.Ident(),
		inst.To,
		md)
}

// GetParent returns the parent basic block of the instruction.
func (inst *InstSIToFP) GetParent() *BasicBlock {
	return inst.Parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *InstSIToFP) SetParent(parent *BasicBlock) {
	inst.Parent = parent
}

// --- [ ptrtoint ] ------------------------------------------------------------

// InstPtrToInt represents a pointer to integer conversion instruction.
//
// References:
//    http://llvm.org/docs/LangRef.html#ptrtoint-instruction
type InstPtrToInt struct {
	// Parent basic block.
	Parent *BasicBlock
	// Name of the local variable associated with the instruction.
	Name string
	// Value before conversion.
	From value.Value
	// Type after conversion.
	To types.Type
	// Map from metadata identifier (e.g. !dbg) to metadata associated with the
	// instruction.
	Metadata map[string]*metadata.Metadata
}

// NewPtrToInt returns a new ptrtoint instruction based on the given source value and target type.
func NewPtrToInt(from value.Value, to types.Type) *InstPtrToInt {
	return &InstPtrToInt{
		From:     from,
		To:       to,
		Metadata: make(map[string]*metadata.Metadata),
	}
}

// Type returns the type of the instruction.
func (inst *InstPtrToInt) Type() types.Type {
	return inst.To
}

// Ident returns the identifier associated with the instruction.
func (inst *InstPtrToInt) Ident() string {
	return enc.Local(inst.Name)
}

// GetName returns the name of the local variable associated with the
// instruction.
func (inst *InstPtrToInt) GetName() string {
	return inst.Name
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *InstPtrToInt) SetName(name string) {
	inst.Name = name
}

// String returns the LLVM syntax representation of the instruction.
func (inst *InstPtrToInt) String() string {
	md := metadataString(inst.Metadata, ",")
	return fmt.Sprintf("%s = ptrtoint %s %s to %s%s",
		inst.Ident(),
		inst.From.Type(),
		inst.From.Ident(),
		inst.To,
		md)
}

// GetParent returns the parent basic block of the instruction.
func (inst *InstPtrToInt) GetParent() *BasicBlock {
	return inst.Parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *InstPtrToInt) SetParent(parent *BasicBlock) {
	inst.Parent = parent
}

// --- [ inttoptr ] ------------------------------------------------------------

// InstIntToPtr represents an integer to pointer conversion instruction.
//
// References:
//    http://llvm.org/docs/LangRef.html#inttoptr-instruction
type InstIntToPtr struct {
	// Parent basic block.
	Parent *BasicBlock
	// Name of the local variable associated with the instruction.
	Name string
	// Value before conversion.
	From value.Value
	// Type after conversion.
	To types.Type
	// Map from metadata identifier (e.g. !dbg) to metadata associated with the
	// instruction.
	Metadata map[string]*metadata.Metadata
}

// NewIntToPtr returns a new inttoptr instruction based on the given source value and target type.
func NewIntToPtr(from value.Value, to types.Type) *InstIntToPtr {
	return &InstIntToPtr{
		From:     from,
		To:       to,
		Metadata: make(map[string]*metadata.Metadata),
	}
}

// Type returns the type of the instruction.
func (inst *InstIntToPtr) Type() types.Type {
	return inst.To
}

// Ident returns the identifier associated with the instruction.
func (inst *InstIntToPtr) Ident() string {
	return enc.Local(inst.Name)
}

// GetName returns the name of the local variable associated with the
// instruction.
func (inst *InstIntToPtr) GetName() string {
	return inst.Name
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *InstIntToPtr) SetName(name string) {
	inst.Name = name
}

// String returns the LLVM syntax representation of the instruction.
func (inst *InstIntToPtr) String() string {
	md := metadataString(inst.Metadata, ",")
	return fmt.Sprintf("%s = inttoptr %s %s to %s%s",
		inst.Ident(),
		inst.From.Type(),
		inst.From.Ident(),
		inst.To,
		md)
}

// GetParent returns the parent basic block of the instruction.
func (inst *InstIntToPtr) GetParent() *BasicBlock {
	return inst.Parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *InstIntToPtr) SetParent(parent *BasicBlock) {
	inst.Parent = parent
}

// --- [ bitcast ] -------------------------------------------------------------

// InstBitCast represents a bitcast instruction.
//
// References:
//    http://llvm.org/docs/LangRef.html#bitcast-instruction
type InstBitCast struct {
	// Parent basic block.
	Parent *BasicBlock
	// Name of the local variable associated with the instruction.
	Name string
	// Value before conversion.
	From value.Value
	// Type after conversion.
	To types.Type
	// Map from metadata identifier (e.g. !dbg) to metadata associated with the
	// instruction.
	Metadata map[string]*metadata.Metadata
}

// NewBitCast returns a new bitcast instruction based on the given source value and target type.
func NewBitCast(from value.Value, to types.Type) *InstBitCast {
	return &InstBitCast{
		From:     from,
		To:       to,
		Metadata: make(map[string]*metadata.Metadata),
	}
}

// Type returns the type of the instruction.
func (inst *InstBitCast) Type() types.Type {
	return inst.To
}

// Ident returns the identifier associated with the instruction.
func (inst *InstBitCast) Ident() string {
	return enc.Local(inst.Name)
}

// GetName returns the name of the local variable associated with the
// instruction.
func (inst *InstBitCast) GetName() string {
	return inst.Name
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *InstBitCast) SetName(name string) {
	inst.Name = name
}

// String returns the LLVM syntax representation of the instruction.
func (inst *InstBitCast) String() string {
	md := metadataString(inst.Metadata, ",")
	return fmt.Sprintf("%s = bitcast %s %s to %s%s",
		inst.Ident(),
		inst.From.Type(),
		inst.From.Ident(),
		inst.To,
		md)
}

// GetParent returns the parent basic block of the instruction.
func (inst *InstBitCast) GetParent() *BasicBlock {
	return inst.Parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *InstBitCast) SetParent(parent *BasicBlock) {
	inst.Parent = parent
}

// --- [ addrspacecast ] -------------------------------------------------------

// InstAddrSpaceCast represents an address space cast instruction.
//
// References:
//    http://llvm.org/docs/LangRef.html#addrspacecast-instruction
type InstAddrSpaceCast struct {
	// Parent basic block.
	Parent *BasicBlock
	// Name of the local variable associated with the instruction.
	Name string
	// Value before conversion.
	From value.Value
	// Type after conversion.
	To types.Type
	// Map from metadata identifier (e.g. !dbg) to metadata associated with the
	// instruction.
	Metadata map[string]*metadata.Metadata
}

// NewAddrSpaceCast returns a new addrspacecast instruction based on the given source value and target type.
func NewAddrSpaceCast(from value.Value, to types.Type) *InstAddrSpaceCast {
	return &InstAddrSpaceCast{
		From:     from,
		To:       to,
		Metadata: make(map[string]*metadata.Metadata),
	}
}

// Type returns the type of the instruction.
func (inst *InstAddrSpaceCast) Type() types.Type {
	return inst.To
}

// Ident returns the identifier associated with the instruction.
func (inst *InstAddrSpaceCast) Ident() string {
	return enc.Local(inst.Name)
}

// GetName returns the name of the local variable associated with the
// instruction.
func (inst *InstAddrSpaceCast) GetName() string {
	return inst.Name
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *InstAddrSpaceCast) SetName(name string) {
	inst.Name = name
}

// String returns the LLVM syntax representation of the instruction.
func (inst *InstAddrSpaceCast) String() string {
	md := metadataString(inst.Metadata, ",")
	return fmt.Sprintf("%s = addrspacecast %s %s to %s%s",
		inst.Ident(),
		inst.From.Type(),
		inst.From.Ident(),
		inst.To,
		md)
}

// GetParent returns the parent basic block of the instruction.
func (inst *InstAddrSpaceCast) GetParent() *BasicBlock {
	return inst.Parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *InstAddrSpaceCast) SetParent(parent *BasicBlock) {
	inst.Parent = parent
}
