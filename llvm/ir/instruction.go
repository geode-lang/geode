// === [ Instructions ] ========================================================
//
// References:
//    http://llvm.org/docs/LangRef.html#instruction-reference

package ir

import "fmt"

// An Instruction represents a non-branching LLVM IR instruction.
//
// Instructions which produce results may be referenced from other instructions,
// and are thus considered LLVM IR values. Note, not all instructions produce
// results (e.g. store).
//
// Instruction may have one of the following underlying types.
//
// Binary instructions
//
// http://llvm.org/docs/LangRef.html#binary-operations
//
//    *ir.InstAdd    (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstAdd)
//    *ir.InstFAdd   (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstFAdd)
//    *ir.InstSub    (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstSub)
//    *ir.InstFSub   (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstFSub)
//    *ir.InstMul    (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstMul)
//    *ir.InstFMul   (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstFMul)
//    *ir.InstUDiv   (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstUDiv)
//    *ir.InstSDiv   (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstSDiv)
//    *ir.InstFDiv   (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstFDiv)
//    *ir.InstURem   (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstURem)
//    *ir.InstSRem   (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstSRem)
//    *ir.InstFRem   (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstFRem)
//
// Bitwise instructions
//
// http://llvm.org/docs/LangRef.html#bitwise-binary-operations
//
//    *ir.InstShl    (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstShl)
//    *ir.InstLShr   (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstLShr)
//    *ir.InstAShr   (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstAShr)
//    *ir.InstAnd    (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstAnd)
//    *ir.InstOr     (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstOr)
//    *ir.InstXor    (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstXor)
//
// Vector instructions
//
// http://llvm.org/docs/LangRef.html#vector-operations
//
//    *ir.InstExtractElement   (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstExtractElement)
//    *ir.InstInsertElement    (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstInsertElement)
//    *ir.InstShuffleVector    (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstShuffleVector)
//
// Aggregate instructions
//
// http://llvm.org/docs/LangRef.html#aggregate-operations
//
//    *ir.InstExtractValue   (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstExtractValue)
//    *ir.InstInsertValue    (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstInsertValue)
//
// Memory instructions
//
// http://llvm.org/docs/LangRef.html#memory-access-and-addressing-operations
//
//    *ir.InstAlloca          (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstAlloca)
//    *ir.InstLoad            (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstLoad)
//    *ir.InstStore           (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstStore)
//    *ir.InstGetElementPtr   (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstGetElementPtr)
//
// Conversion instructions
//
// http://llvm.org/docs/LangRef.html#conversion-operations
//
//    *ir.InstTrunc           (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstTrunc)
//    *ir.InstZExt            (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstZExt)
//    *ir.InstSExt            (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstSExt)
//    *ir.InstFPTrunc         (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstFPTrunc)
//    *ir.InstFPExt           (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstFPExt)
//    *ir.InstFPToUI          (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstFPToUI)
//    *ir.InstFPToSI          (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstFPToSI)
//    *ir.InstUIToFP          (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstUIToFP)
//    *ir.InstSIToFP          (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstSIToFP)
//    *ir.InstPtrToInt        (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstPtrToInt)
//    *ir.InstIntToPtr        (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstIntToPtr)
//    *ir.InstBitCast         (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstBitCast)
//    *ir.InstAddrSpaceCast   (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstAddrSpaceCast)
//
// Other instructions
//
// http://llvm.org/docs/LangRef.html#other-operations
//
//    *ir.InstICmp     (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstICmp)
//    *ir.InstFCmp     (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstFCmp)
//    *ir.InstPhi      (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstPhi)
//    *ir.InstSelect   (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstSelect)
//    *ir.InstCall     (https://godoc.org/github.com/geode-lang/geode/llvm/ir#InstCall)
type Instruction interface {
	fmt.Stringer
	// GetParent returns the parent basic block of the instruction.
	GetParent() *BasicBlock
	// SetParent sets the parent basic block of the instruction.
	SetParent(parent *BasicBlock)
}
