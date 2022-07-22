.format "bin"

	.org $E000
	.memory "fill", $E000, $2000, $ff
do_start:	LDA #12
			STA 12
loop:
	jmp loop

do_nmi: NOP
		RTI
	 
do_irq: NOP
		RTI

	.org  $FFFA
	.word   do_nmi
	.word   do_start
	.word   do_irq
