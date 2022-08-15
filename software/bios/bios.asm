.format "bin"

	.memory "fill", $E000, $2000, $ea
	.org $E000
	IOBASE .equ $B000
	VIA .equ IOBASE
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
	ACIA .equ IOBASE + $0100

	JTIME .equ $A0 ; to $A2 three bytes jiffy time

do_reset:
    ldx #$ff ; set the stack pointer 
   	txs 

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

do_scinit:
	rts

do_ioinit:
	rts

do_ramtas:
	rts

do_getin:
	rts

do_iobase:
	ldx IOBASE & $00ff
	ldy IOBASE >> 8 & $00ff
	rts

do_settim:
	rts

do_rdtim:
	rts

do_udtim:
	inc JTIME     ; increment Low-Byte 
    bne jend          ; =0?
    inc JTIME+1   ; ja, dann Überlauf 255->0 und High-Byte auch erhöhen
    bne jend          ; =0?
    inc JTIME+2   ; ja, dann Überlauf 255->0 und High-Byte auch erhöhen
	; TODO on E1EB00 (4.320.000 ticks or 24 Hours) reset to null
jend:
	rts

do_nmi: NOP
		RTI
	 
do_irq: NOP
		RTI

	.org $FF81 ; SCINIT Initialize "Screen output", (here only the serial monitor)
	jmp do_scinit
	
	.org $FF84 ; IOINIT Initialize VIA & IRQ
	jmp do_ioinit

	.org $FF87 ; RAMTAS RAM test and search memory end
	jmp do_ramtas

	.org $FFDB ; SETTIM Set the Jiffy Clock
	jmp do_settim

	.org $FFDE ; RDTIM read the Jiffy Clock
	jmp do_rdtim

	.org $FFEA ; UDTIM Tick the Jiffy Clock
	jmp do_udtim

	.org $FFE4 ; GETIN Read a byte from the input channel
	jmp do_getin

	.org $FFF3 ; IOBASE	Read the base address of I/O chips
	jmp do_iobase

	.org  $FFFA
	.word   do_nmi
	.word   do_reset
	.word   do_irq
