	.section	__TEXT,__text,regular,pure_instructions
	.macosx_version_min 10, 13
	.intel_syntax noprefix
	.globl	_exp                    ## -- Begin function exp
	.p2align	4, 0x90
_exp:                                   ## @exp
	.cfi_startproc
## %bb.0:                               ## %entry_1
	sub	rsp, 24
	.cfi_def_cfa_offset 32
	mov	qword ptr [rsp + 16], rdi
	mov	qword ptr [rsp + 8], rsi
	mov	rsi, qword ptr [rsp + 8]
	test	rsi, rsi
	jne	LBB0_2
	jmp	LBB0_1
LBB0_1:                                 ## %if_0_then_2
	mov	eax, 1
                                        ## kill: def %rax killed %eax
	add	rsp, 24
	ret
LBB0_2:                                 ## %if_0_else_3
	jmp	LBB0_3
LBB0_3:                                 ## %if_0_merge_4
	mov	rax, qword ptr [rsp + 16]
	mov	rdi, qword ptr [rsp + 16]
	mov	rcx, qword ptr [rsp + 8]
	sub	rcx, 1
	mov	rsi, rcx
	mov	qword ptr [rsp], rax    ## 8-byte Spill
	call	_exp
	mov	rcx, qword ptr [rsp]    ## 8-byte Reload
	imul	rcx, rax
	mov	rax, rcx
	add	rsp, 24
	ret
	.cfi_endproc
                                        ## -- End function

.subsections_via_symbols
