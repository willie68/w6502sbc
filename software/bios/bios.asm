.format "bin"

	.org $E000
	.memory "fill", $E000, $2000, $ff
	VIA .equ $D000
	VIA_ORB .equ VIA
	VIA_ORA .equ VIA+1
	VIA_DDRB .equ VIA+2
	VIA_DDRA .equ VIA+3
	VIA_T1Cl .equ VIA+4
	VIA_T1CH .equ VIA+5
	VIA_T1LL .equ VIA+6
	VIA_T1LH .equ VIA+7
	VIA_T2CL .equ VIA+8
	VIA_T2CH .equ VIA+9
	VIA_SR .equ VIA+$A
	VIA_ACR .equ VIA+$B
	VIA_PCR .equ VIA+$C
	VIA_IFR .equ VIA+$D
	VIA_IER .equ VIA+$E
	VIA_IRA .equ VIA+$F
	ACIA .equ $D100

do_reset:	
// setting up the 65C22 VIA
	LDA #$FF
	STA VIA_DDRA
	LDA #$AA
	STA VIA_ORA
blinkloop:
    delay(250)
	ROR
	STA VIA_ORA
	jmp blinkloop

.macro delay(ms)
			ldy #ms
delayl2:	ldx #20
delayl1: 	dex          ; (2 cycles)
        	bne  delayl1   ; (3 cycles in loop, 2 cycles at end)
        	dey          ; (2 cycles)
        	bne  delayl2   ; (3 cycles in loop, 2 cycles at end)
.endmacro

do_nmi: NOP
		RTI
	 
do_irq: NOP
		RTI

	.org  $FFFA
	.word   do_nmi
	.word   do_reset
	.word   do_irq
