.format "bin"
.target "65C02"

.memory "fill", $E000, $2000, $ea
.org $E000
.include "io.asm" 

do_reset:	
// setting up the 65C22 VIA
l1:
	lda #$ff
	clc
	bcs l1
	bcc l2
	lda #$80
l2:
	lda #$40
	rts

do_nmi: NOP
		pha
		
		pla
		RTI
	 
do_irq: NOP
		RTI

	.org  $FFFA
	.word   do_nmi
	.word   do_reset
	.word   do_irq
