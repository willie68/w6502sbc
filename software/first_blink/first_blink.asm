.format "bin"
.target "65C02"

.memory "fill", $E000, $2000, $ea
.org $E000
.include "io.asm" 

do_reset:	
// setting up the 65C22 VIA
	LDA #$FF
	STA VIA_DDRA
	LDA #$AA
	STA VIA_ORA
blinkloop:
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
