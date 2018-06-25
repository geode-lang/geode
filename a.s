	.section	__TEXT,__text,regular,pure_instructions
	.macosx_version_min 10, 13
	.intel_syntax noprefix
	.globl	_fib                    ## -- Begin function fib
	.p2align	4, 0x90
_fib:                                   ## @fib
## %bb.0:                               ## %entry_1
	push	r14
	push	rbx
	push	rax
	mov	rbx, rdi
	cmp	rbx, 1
	jg	LBB0_3
## %bb.1:                               ## %if_0_then_2
	mov	rax, rbx
	add	rsp, 8
	pop	rbx
	pop	r14
	ret
LBB0_3:                                 ## %if_0_merge_4
	lea	rdi, [rbx - 1]
	call	_fib
	mov	r14, rax
	add	rbx, -2
	mov	rdi, rbx
	call	_fib
	add	rax, r14
	add	rsp, 8
	pop	rbx
	pop	r14
	ret
                                        ## -- End function
	.globl	_main                   ## -- Begin function main
	.p2align	4, 0x90
_main:                                  ## @main
## %bb.0:                               ## %entry_5
	push	rax
	mov	edi, 30
	call	_fib
	mov	rcx, rax
	lea	rdi, [rip + _.str_6]
	xor	eax, eax
	mov	rsi, rcx
	call	_printf
	xor	eax, eax
	pop	rcx
	ret
                                        ## -- End function
	.section	__TEXT,__const
	.globl	_.str_6                 ## @.str_6
_.str_6:
	.asciz	"%d\n"


.subsections_via_symbols
