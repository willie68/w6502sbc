.format "bin"

	.org $E000
	.memory "fill", $E000, $2000, $ff
	_VIA .equ $D000
	_ORB .equ _VIA
	_ORA .equ _VIA+1
	_DDRB .equ _VIA+2
	_DDRA .equ _VIA+3
	_T1Cl .equ _VIA+4
	_T1CH .equ _VIA+5
	_T1LL .equ _VIA+6
	_T1LH .equ _VIA+7
	_T2CL .equ _VIA+8
	_T2CH .equ _VIA+9
	_SR .equ _VIA+$A
	_ACR .equ _VIA+$B
	_PCR .equ _VIA+$C
	_IFR .equ _VIA+$D
	_IER .equ _VIA+$E
	_IRA .equ _VIA+$F
	_ACIA .equ $D100
do_reset:	
// setting up the 65C22 VIA
	LDA #$FF
	STA _DDRA
	LDA #$AA
	STA _ORA
blinkloop:
    delay(250)
	ROR
	STA _ORA
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
