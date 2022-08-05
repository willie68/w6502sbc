.format "bin"
.include "../include/io.asm"
// for the input and outputs of the TPS we will take port B
	.org $E000
	.memory "fill", $E000, $2000, $ff

do_reset:	
// setting up the 65C22 VIA
	LDA #$FF
	STA VIA_DDRA
	LDA #$AA
	STA VIA_ORA
mainloop:
	lda #12
	TAY
	ROR
	ROR
	ROR
	ROR
	TAX
	JMP (op_code,x)

.macro delay(ms)
			ldy #ms
delayl2:	ldx #20
delayl1: 	dex          ; (2 cycles)
        	bne  delayl1   ; (3 cycles in loop, 2 cycles at end)
        	dey          ; (2 cycles)
        	bne  delayl2   ; (3 cycles in loop, 2 cycles at end)
.endmacro

op_nop:
    jmp mainloop
op_port:
    jmp mainloop
op_delay:
    jmp mainloop
op_rjmp:
    jmp mainloop
op_lda:
    jmp mainloop
op_xeqa:
    jmp mainloop
op_aeqe:
    jmp mainloop
op_math:
    jmp mainloop
op_page:
    jmp mainloop
op_jump:
    jmp mainloop
op_cloop:
    jmp mainloop
op_dloop:
    jmp mainloop
op_skip:
    jmp mainloop
op_call:
    jmp mainloop
op_ret:
    jmp mainloop
op_byte:
    jmp mainloop

do_nmi: NOP
		RTI
	 
do_irq: NOP
		RTI


.org  $FF00

// jump table
op_code: 
	.word op_nop
	.word op_port
	.word op_delay
	.word op_rjmp
	.word op_lda
	.word op_xeqa
	.word op_aeqe
	.word op_math
	.word op_page
	.word op_jump
	.word op_cloop
	.word op_dloop
	.word op_skip
	.word op_call
	.word op_ret
	.word op_byte

	.org  $FFFA
	.word   do_nmi
	.word   do_reset
	.word   do_irq
