	.section	__TEXT,__text,regular,pure_instructions
	.macosx_version_min 10, 13
	.intel_syntax noprefix
	.globl	_main                   ## -- Begin function main
	.p2align	4, 0x90
_main:                                  ## @main
	.cfi_startproc
## %bb.0:                               ## %entry_5
	sub	rsp, 24
	.cfi_def_cfa_offset 32
	mov	qword ptr [rsp + 16], 0
LBB0_1:                                 ## %for_0_cond_6
                                        ## =>This Inner Loop Header: Depth=1
	mov	rax, qword ptr [rsp + 16]
	sub	rax, 25000
	mov	qword ptr [rsp + 8], rax ## 8-byte Spill
	jg	LBB0_3
	jmp	LBB0_2
LBB0_2:                                 ## %for_0_body_7
                                        ##   in Loop: Header=BB0_1 Depth=1
	lea	rdi, [rip + _.str_8]
	mov	rsi, qword ptr [rsp + 16]
	mov	al, 0
	call	_printf
	mov	rsi, qword ptr [rsp + 16]
	add	rsi, 1
	mov	qword ptr [rsp + 16], rsi
	jmp	LBB0_1
LBB0_3:                                 ## %for_0_merge_9
	mov	eax, 1
                                        ## kill: def %rax killed %eax
	add	rsp, 24
	ret
	.cfi_endproc
                                        ## -- End function
	.section	__TEXT,__const
	.globl	_.str_8                 ## @.str_8
_.str_8:
	.asciz	"%d\n"


.subsections_via_symbols
