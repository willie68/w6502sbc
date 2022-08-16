.format "bin"

	.memory "fill", $E000, $2000, $ea
	.org $E000
	IOBASE .equ $B000
	VIA .equ IOBASE
	VIA_ORB .equ VIA
	VIA_ORA .equ VIA+1
	VIA_DDRB .equ VIA+2
	VIA_DDRA .equ VIA+3
	VIA_T1CL .equ VIA+4
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

	;constants for board specifig
	JIFFY_VIA_TIMER_LOAD .equ 20000   ; this is the value for 1MHZ / 50 ticks per second

	; ZERO Page registers $0000.. $00ff
	RAMTOP .equ $31 ; store the page of the last RAM ($30 is the low adress)
	JTIME .equ $A0 ; to $A2 three bytes jiffy time
	IN_READ .equ $80
	IN_WRITE .equ $81

	; Stack  $0100.. $01ff
	SPAGE .equ $0100
	; Bios data
	BIOSPAGE .equ $0200
	IRQ_SRV .equ  $0214    ; $0214 LOW byte, $0215 HIGH byte for a external irq service routine
	NMI_SRV .equ  $0216    ; $0216 LOW byte, $0217 HIGH byte for a external nmi service routine
	RTI_SRV .equ  $0218    ; every user irq or nmi routine should call this for returning, jmp (RTI_SRV)
	IN_BUF_LEN .equ $0F    ; length of input buffer
	IN_BUFFER .equ $0280   ; 16 bytes of input buffer

	; BASIC data
	BASICPAGE .equ $0300
	; RAM start
	RAMSTART .equ $0400

do_reset:
    ldx #$ff ; set the stack pointer 
   	txs 

	jsr do_ioinit
	jsr do_scinit
	jsr do_ramtas
	jsr do_srvinit

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
delayl2:	
	ldx #20
delayl1: 	
	dex          ; (2 cycles)
    bne  delayl1   ; (3 cycles in loop, 2 cycles at end)
    dey          ; (2 cycles)
    bne  delayl2   ; (3 cycles in loop, 2 cycles at end)
.endmacro

do_scinit:
	rts

do_ioinit:
    sei
	; disable all interrupts
  	stz VIA_IER  
	; setting free run mode with interrupts enabled
	lda #%01000000
  	sta VIA_ACR     
	lda #%11000000
  	sta VIA_IER  ; enable interrupt for timer 1
	; setting the vias timer 1 in free run mode, jiffy timer, load with 
	lda #<JIFFY_VIA_TIMER_LOAD
	sta VIA_T1LL 
	lda #>JIFFY_VIA_TIMER_LOAD
	sta VIA_T1LH 
	cli
	rts

do_ramtas: 
	sei
	lda #$00
	tay
	; clear memory on zeropage, stack, biospage, basicpage
ramtas_l1:
	sta $0000, y
	sta SPAGE, y
	sta BIOSPAGE, y
	sta BASICPAGE, y
	iny
	bne ramtas_l1
	; checking every page 0 byte to get the last RAM Page
	tay
	sta RAMTOP-1     ; put a 0 into $30 for later indirect acces to $30 $31 for RAM Test adress
	lda >RAMSTART-1
	sta RAMTOP
ramtas_l2:
	inc RAMTOP
	lda #$55         ; test with 01010101
	sta (RAMTOP), y
	cmp (RAMTOP), y
	bne ramtas_ramtop
	rol				 ; test with 10101010
	sta (RAMTOP), y
	cmp (RAMTOP), y
	bne ramtas_ramtop
	jmp ramtas_l2
ramtas_ramtop:
	dec RAMTOP  ; found the last RAM page at adress one page before 
	cli
	rts

do_restor:
    jsr do_srvinit
	stz IRQ_SRV
	stz IRQ_SRV+1
	stz NMI_SRV
	stz NMI_SRV+1
	rts

do_srvinit:
	; saving the isr return adress to kernel page
	lda <isr_end
	sta RTI_SRV
	lda >isr_end
	sta RTI_SRV+1
	rts

do_putin:
	; adding something to the input ring buffer , its always overwriting
	ldx IN_WRITE
	sta IN_BUFFER, x
	dec IN_WRITE
	bpl putin_end  
	ldx #IN_BUF_LEN
	stx IN_WRITE
putin_end:
	rts

do_getin:
	; ring buffer to get a char
	lda IN_READ
	cmp IN_WRITE
	bne getin_jmp1
	lda #$00 ; this means nothing in the buffer
	rts
getin_jmp1:
	tax
	lda IN_BUFFER, x
	dec IN_READ
	bpl getin_end  
	ldx #IN_BUF_LEN
	stx IN_READ
getin_end:
	rts

do_iobase:
	ldx IOBASE & $00ff
	ldy IOBASE >> 8 & $00ff
	rts

	; higher two bytes of the 3 bytes of the 1/50 secs of a day. 24h * 60m * 60s * 50, 4.320.000 ticks per day
    jiffyday .equ $E1EB  

do_settim:
	sei
	sta JTIME+2
	stx JTIME+1
	sty JTIME
	cli
	rts

do_rdtim:
	sei
	lda JTIME+2
	ldx JTIME+1
	ldy JTIME
	cli
	rts

do_udtim:
	inc JTIME+2     ; increment Low-Byte 
    bne jiend     ; =0?
    inc JTIME+1   ; ja, dann Überlauf 255->0 und High-Byte auch erhöhen
    bne jtest     ; =0?
    inc JTIME   ; ja, dann Überlauf 255->0 und High-Byte auch erhöhen
jtest:
	; test on E1EB00 (4.320.000 ticks or 24 Hours) reset to null
	sec
	lda JTIME+1 
	sbc #<jiffyday 
	lda JTIME
	sbc #>jiffyday 
	bcc jiend
	; reset to null
	stz JTIME +1
	stz JTIME +2
jiend:
	rts


do_nmi: 
	pha
	phx
	phy
	; look if nmi service routine is set
	lda #$00
	cmp >NMI_SRV
	beq nmi_end
	jmp (NMI_SRV)
nmi_end:
	ply
	plx
	pla
	rti
	 
do_irq: 
	pha
	phx
	phy
	; testing for timer 1, jiffy timer interrupt
	bit VIA_IFR          ; Bit 6 copied to overflow flag
  	bvc isr_no_timer1
	lda VIA_T1CL         ; Clears the interrupt
	jsr do_udtim
isr_no_timer1:
	; here do other isr stuff
	; look if irq service routine is set
	lda #$00
	cmp >IRQ_SRV
	beq isr_end
	jmp (IRQ_SRV)
isr_end:
	ply
	plx
	pla
	rti

end_of_kernel:

	.org $FF81 ; SCINIT Initialize "Screen output", (here only the serial monitor)
	jmp do_scinit
	
	.org $FF84 ; IOINIT Initialize VIA & IRQ
	jmp do_ioinit

	.org $FF87 ; RAMTAS RAM test and search memory end
	jmp do_ramtas

	.org $FF8A ; RESTOR restore default kernel vectors
	jmp do_restor 

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
