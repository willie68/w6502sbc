.format "bin"
.target "65C02"

.memory "fill", $E000, $2000, $ea
.org $E000
.include "io.asm" 
.include "zp.asm"
.include "macros.asm"

;----- constant definitions -----
;constants for board specifig
JIFFY_VIA_TIMER_LOAD .equ 20000   ; this is the value for 1MHZ / 50 ticks per second

; Stack  $0100.. $01ff
STACK .equ $0100
; Bios data
BIOSPAGE .equ $0200
IRQ_SRV .equ  $0214    ; $0214 LOW byte, $0215 HIGH byte for a external irq service routine
BRK_SRV .equ  $0216    ; $0216 LOW byte, $0217 HIGH byte for a external nmi service routine
NMI_SRV .equ  $0218    ; $0216 LOW byte, $0217 HIGH byte for a external nmi service routine
IN_BUF_LEN .equ $0F    ; length of input buffer
IN_BUFFER .equ $0280   ; 16 bytes of input buffer

; BASIC data
BASICPAGE .equ $0300
; RAM start for testing
RAMSTART .equ $0400

;----- bios start code -----
do_reset: ; bios reset routine 
    ldx #$ff ; set the stack pointer 
   	txs 

	lda #$00
	ldx #$00
	ldy #$00
	jsr do_settim

	jsr do_srvinit ; cleanup the interrupt registers
	jsr do_ioinit  ; initialise port A an timer of VIA
	jsr do_scinit


	;jsr lcd_clear
	msg_out(message_welcome)
;	jsr lcd_clear
;	msg_out(message_ramtas)
;	jsr do_ramtas

	jsr lcd_clear
	msg_out(message_ready)

/*	jsr lcd_secondrow
	lda #$04					; output 1059 = $0423
	ldx #$23
	jsr do_numout
	lda #" "
	jsr do_chrout
	lda #$04					; output 1059 = $0423
	jsr do_hexout
	lda #" "
	jsr do_chrout
*/
	stz COUNTER
;----- main -----
.macro output(col)

	lda #$01
	ldx #col
	jsr lcd_goto 

	lda #col
	jsr do_nhexout
.endmacro

CNT .var 0
.loop 16
	output(CNT)
	CNT = CNT + 1
.endloop


main_loop:
	jmp main_loop

do_ioinit: ; initialise the timer for the jiffy clock
	lda #%01111111
  	sta VIA_IER  ; disable all interrupts
	; setting free run mode with interrupts enabled
	lda #%01000011
  	sta VIA_ACR     
	lda #%11000000
  	sta VIA_IER  ; enable interrupt for timer 1
	; setting the vias timer 1 in free run mode, jiffy timer, load with 
	lda #<JIFFY_VIA_TIMER_LOAD
	sta VIA_T1LL 
	sta VIA_T1CL
	lda #>JIFFY_VIA_TIMER_LOAD
	sta VIA_T1LH
	sta VIA_T1CH

    ; ACIA setup
/*    lda #$00
    sta ACIA_STATUS
    lda #$0b
    sta ACIA_COMMAND
    lda #$1f
    sta ACIA_CONTROL
*/
	lda #$FF
	sta VIA_DDRA
	lda #$00
	sta VIA_ORA
	rts

do_ramtas: ; initialising memory page 0, stack, bios, basic, lokking for the address of the last RAM page, write it to RAMTOP
	lda #$00
	tay
	; clear memory on zeropage, stack, biospage, basicpage, determing the memory max page
ramtas_l1:
	sta $0000, y
	sta STACK, y
	sta BIOSPAGE, y
	sta BASICPAGE, y
	iny
	bne ramtas_l1
	; checking 0 byte from every page to get the last RAM Page
	tay
	stz RAMTOP     	; put a 0 into $30 for later indirect acces to $30 $31 for RAM Test adress
	lda >RAMSTART-1    	;
	sta RAMTOP+1 			;after this, the RAMTOP should be set to $0400
ramtas_l2:
	inc RAMTOP+1
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
	dec RAMTOP+1  ; found the last RAM page at adress one page before 
	rts

do_restor: ; restoring interrupt vectors to 0
    jsr do_srvinit
	stz IRQ_SRV
	stz IRQ_SRV+1
	stz NMI_SRV
	stz NMI_SRV+1
	rts

do_srvinit: ; saving the isr return adress to kernel page
	stz NMI_SRV
	stz IRQ_SRV
	stz BRK_SRV
	rts

do_putin: ; adding something to the input ring buffer , its always overwriting
	ldx IN_WRITE
	sta IN_BUFFER, x
	dec IN_WRITE
	bpl putin_end  
	ldx #IN_BUF_LEN
	stx IN_WRITE
putin_end:
	rts

do_getin: ; ring buffer to get a char
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

do_iobase: ; return the io base address lo: X, hi: Y
	ldx #<IOBASE
	ldy #>IOBASE
	rts

    jiffyday .equ $E1EB  

; ---- Jiffy clock ----
do_settim: ; setting the jiffy clock to the value, lo: Y, mid: X, hi: A
	sei
	sta JTIME+2
	stx JTIME+1
	sty JTIME
	cli
	rts

do_rdtim:; reading the actual jiffy clock, lo: Y, mid: X, hi: A
	sei
	lda JTIME+2
	ldx JTIME+1
	ldy JTIME
	cli
	rts

do_udtim: ; update the jiffy clock
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

.include "lcd.asm"

;------------------------------------------------------------------------------
; The count of outloops will be used from A.
; for 1MHz we had a cycle with 1us. if A = 1 we had 20 + 20 clks, which means a minimum of 200us,
; but the reality is somtime different. To get the 200us on my sbc there must be 32 inner loops.
; $01 = 200uS, $02= 360us, $04= 700uS, $08= 1,4ms, $10= 2,7ms, $20= 5,3ms, $40= 10,6ms, $80= 21,3ms, $FF=42,3ms
do_delay:
	phy			; 3 clk
@outer:    
	ldy  #$20	; 2 clk, this gives an inner loop of 5 cycles x 20 =  100uS        
@inner:
    dey			; 2 clk
	bne @inner	; 2 + 1 clk (for the jump back)
    sbc #$01	; 2 clk
    bne @outer	; 2 + 1 clk exit when COUNTER is less than 0
    ply			; 4 clk
    rts			; 6 clk 

;---- Interrupt service routines ----
do_brk;
	lda #$00
	cmp >BRK_SRV
	beq isr_end
	jmp (BRK_SRV)

do_nmi: ; nmi service routine
	pha
	phx
	phy
	; look if an external nmi service routine is set
	lda NMI_SRV
	cmp #$00
	beq isr_end
	jmp (NMI_SRV)
	 
do_irq: ; irq service routine
	pha
	phx
	phy
	; check for brk
/*	php					; put status to stack
	pla					; get staus in A
	and #$10
	beq @irq1
	jmp do_brk
*/
	; testing for timer 1, jiffy timer interrupt
@irq1:
	bit VIA_IFR          	; Bit 6 copied to overflow flag
  	bvc @isr_no_timer1

	bit VIA_T1CL         	; Clears the interrupt
	jsr do_udtim
	jmp isr_end
@isr_no_timer1:
	; here do other isr stuff
	; look if an external irq service routine is set
	lda IRQ_SRV
	cmp #$00
	beq isr_end
	jmp (IRQ_SRV)
isr_end: ; this is the ending for all interrupt service routines
	ply
	plx
	pla
	rti

do_setirqsrv: ; setting an external irq routine for checking, A hi, X lo
	stx IRQ_SRV
	sta IRQ_SRV+1
	rts
do_setbrksrv: ; setting an external irq routine for checking, A hi, X lo
	stx BRK_SRV
	sta BRK_SRV+1
	rts
do_setnmisrv: ; setting an external irq routine for checking, A hi, X lo
	stx NMI_SRV
	sta NMI_SRV+1
	rts

;----- Messages of the bios -----
	message_w6502sbc: .asciiz "W6502SBC"
	message_welcome: .asciiz "Welcome"
	message_ramtas: .asciiz "RAMTAS"
	message_srvinit: .asciiz "SRV INIT"
	message_ready: .asciiz "ready"
	message_showdec: .asciiz "show dec"
	message_britta: .asciiz "Hallo Britta"

;----- jump table for bios routines -----
jump_table: 

STROUT:
	.org $FF00 ; STROUT output string, A = high, X = low
	jmp do_strout
HEXOUT:
	.org $FF03 ; HEXOUT output a 16-bit value as hex with leading $, A = high, X = low
	jmp do_hexout
BHEXOUT:
	.org $FF06 ; BHEXOUT output a 8-bit value as hex without leading, A = value
	jmp do_bhexout
NUMOUT:
	.org $FF09 ; NUMOUT output a 16-bit value as decimal, A = high, X = low
	jmp do_numout
LCDCLEAR:
	.org $FF0C ; LCDCLEAR Clear display
	jmp lcd_clear
LCDHOME:
	.org $FF0F ; LCDHOME goto 1. line 1. column
	jmp lcd_home
LCDSECROW:
	.org $FF12 ; LCDSECROW goto 2. line 1. Column
	jmp lcd_secondrow
SCINIT:
	.org $FF81 ; SCINIT Initialize "Screen output", (here only the serial monitor)
	jmp do_scinit
IOINIT:
	.org $FF84 ; IOINIT Initialize VIA & IRQ
	jmp do_ioinit
RAMTAS:
	.org $FF87 ; RAMTAS RAM test and search memory end
	jmp do_ramtas
RESTORE:
	.org $FF8A ; RESTOR restore default kernel vectors
	jmp do_restor 
SETIRQ:
	.org $FFA0 ; SETIRQ set external irq vectors
	jmp do_setirqsrv 
SETBRK:
	.org $FFA3 ; SETBRK set external brk vectors
	jmp do_setbrksrv 
SETNMI:
	.org $FFA6 ; SETNMI set external nmi vectors
	jmp do_setnmisrv 
ISREND:
	.org $FFA9 ; ISREND return from irq, nmi, brk
	jmp isr_end 
CHROUT:
	.org $FFD2 ; CHROUT Output an character
	jmp do_chrout
SETTIM:
	.org $FFDB ; SETTIM Set the Jiffy Clock
	jmp do_settim
RDTIM:
	.org $FFDE ; RDTIM read the Jiffy Clock
	jmp do_rdtim
UDTIM:
	.org $FFEA ; UDTIM Tick the Jiffy Clock
	jmp do_udtim
GETIN:
	.org $FFE4 ; GETIN Read a byte from the input channel
	jmp do_getin
RDIOBS:
	.org $FFF3 ; RDIOBS	Read the base address of I/O chips
	jmp do_iobase

;----- cpu vectors -----
	.org  $FFFA
	.word   do_nmi
	.word   do_reset
	.word   do_irq
