.format "bin"
.target "65C02"

.memory "fill", $E000, $2000, $ea
.org $E000
.include "io.asm" 

;----- constant definitions -----
;constants for board specifig
JIFFY_VIA_TIMER_LOAD .equ 20000   ; this is the value for 1MHZ / 50 ticks per second
; constants for LCD
LCD_E  .equ %10000000
LCD_RW .equ %01000000
LCD_RS .equ %00100000


; ZERO Page registers $0000.. $00ff
COUNTER .equ $20 ; counter for different things
HNIBBLE .equ $21
LNIBBLE .equ $22
RAMTOP .equ $30 ; store the page of the last RAM ($30 is the low adress)

JTIME .equ $A0 ; to $A2 three bytes jiffy time, higher two bytes of the 3 bytes of the 1/50 secs of a day. 24h * 60m * 60s * 50, 4.320.000 ticks per day

IN_READ .equ $80
IN_WRITE .equ $81
TEMP_VEC .equ $32 ; store a temporary vector, like the address to the string to output, $32 low, $33 hi
; Stack  $0100.. $01ff
SPAGE .equ $0100
; Bios data
BIOSPAGE .equ $0200
IRQ_SRV .equ  $0214    ; $0214 LOW byte, $0215 HIGH byte for a external irq service routine
NMI_SRV .equ  $0216    ; $0216 LOW byte, $0217 HIGH byte for a external nmi service routine
RTI_SRV .equ  $0218    ; every user irq or nmi routine should call this for returning from interrupt, jmp (RTI_SRV)
IN_BUF_LEN .equ $0F    ; length of input buffer
IN_BUFFER .equ $0280   ; 16 bytes of input buffer

; BASIC data
BASICPAGE .equ $0300
; RAM start for testing
RAMSTART .equ $0400

;----- macros -----
.macro msg_out(msg)
	lda #>msg
	ldx #<msg
	jsr do_strout 
.endmacro

.macro toggle_a()
	pha
	lda #$ff
	sta VIA_ORA
	stz VIA_ORA
	pla
.endmacro
;----- bios code -----
do_reset: ; bios reset routine 
	sei
    ldx #$ff ; set the stack pointer 
   	txs 

	jsr do_ioinit
	jsr do_scinit
	
	;jsr lcd_clear
	msg_out(message_w6502sbc)

;	jsr lcd_clear
;	msg_out(message_ramtas)
;	jsr do_ramtas

	jsr lcd_clear
	msg_out(message_srvinit)
	jsr do_srvinit

;----- main -----
	jsr lcd_clear
	msg_out(message_ready)
	jsr lcd_secondrow
	msg_out(message_britta)

	cli
main_loop:
	lda #$00 ; 34,7ms
	toggle_a()
	jsr do_delay
	toggle_a()

	lda #$01  ;160us
	toggle_a()
	jsr do_delay
	toggle_a()

	lda #$02 ; 303us
	toggle_a()
	jsr do_delay
	toggle_a()

	lda #$04 ; 575us
	toggle_a()
	jsr do_delay
	toggle_a()

	lda #$08 ; 1,1 ms
	toggle_a()
	jsr do_delay
	toggle_a()

	lda #$10 ; 2,2 ms
	toggle_a()
	jsr do_delay
	toggle_a()

	lda #$20 ; 4,4 ms
	toggle_a()
	jsr do_delay
	toggle_a()

	lda #$40 ; 8,7ms
	toggle_a()
	jsr do_delay
	toggle_a()

	lda #$80 ; 17,4
	toggle_a()
	jsr do_delay
	toggle_a()

	lda #$ff ; 34,7
	toggle_a()
	jsr do_delay
	toggle_a()
	jmp main_loop

do_ioinit: ; initialise the timer for the jiffy clock
/*    sei
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
	sta SPAGE, y
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
	lda <isr_end
	sta RTI_SRV
	lda >isr_end
	sta RTI_SRV+1
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

; ---- Display routines ----
do_scinit: 		; initialise LC-Display on port B
	; D4..D7 on Port pins PB0..3
	; RS; R/W and E on Port pins PB5, PB6, PB7
	lda #$ff 	; Set all pins on port B to output
  	sta VIA_DDRB
	lda #0 		; all pins low
	sta VIA_ORB

	; reset the display, wait at least 15ms
	lda #$58
	jsr do_delay

	; send 3 times the reset...
  	lda #(%00000011 | LCD_E) ; 1. RESET
	sta VIA_ORB
	eor #LCD_E
	sta VIA_ORB
	lda #$1f
	jsr do_delay

  	lda #(%00000011 | LCD_E) ; 1. RESET
	sta VIA_ORB
	eor #LCD_E
	sta VIA_ORB
	lda #$01
	jsr do_delay

  	lda #(%00000011 | LCD_E) ; 1. RESET
	sta VIA_ORB
	eor #LCD_E
	sta VIA_ORB
	lda #$01
	jsr do_delay

  	lda #(%00000010 | LCD_E) ; Set 4-bit mode; 
	sta VIA_ORB
	eor #LCD_E
	sta VIA_ORB
	lda #$01
	jsr do_delay
	lda #$01
	jsr do_delay
	
  	lda #%00101000 ; 2-line display; 5x8 font
  	jsr lcd_instruction

  	lda #%00001110 ; Display on; cursor on; blink off
  	jsr lcd_instruction

  	lda #%00000110 ; Increment and shift cursor; don't shift display
  	jsr lcd_instruction

  	lda #%00000010 ; Return home
  	jsr lcd_instruction

  	lda #%00000001 ; Clear display
  	jsr lcd_instruction
	rts

lcd_wait: ; wait until the LCD is not busy
	pha
	lda #%11110000 ;set PORTB pins 0 - 3 as input
	sta VIA_DDRB
@lcdbusy:
	lda #LCD_RW
	sta VIA_ORB
	ora #LCD_E
	sta VIA_ORB
	; loding high nibble with busy flag
	lda VIA_ORB
	sta HNIBBLE
	lda #LCD_RW
	sta VIA_ORB
	ora #LCD_E
	sta VIA_ORB
	; getting the low nibble, address counter
	lda VIA_ORB
	sta LNIBBLE
	lda #LCD_RW
	sta VIA_ORB
	lda HNIBBLE
	and #%00001000 ; mask the busy flag
	bne @lcdbusy
	lda #$FF ; setting port to output again
	sta VIA_DDRB
	pla
	rts

lcd_instruction: ; sending A as an instruction to LCD
;	jsr lcd_wait
	pha
	pha
	lsr
	lsr
	lsr
	lsr
	ora #LCD_E
	sta VIA_ORB
	eor #LCD_E
	sta VIA_ORB
	pla
	and #$0f
	ora #LCD_E
	sta VIA_ORB
	eor #LCD_E
	sta VIA_ORB
;	jsr do_delay
	pla
	rts

lcd_secondrow: ; move cursor to second row
	pha
  	;jsr lcd_wait
  	lda #%10000000 + $40
  	jsr lcd_instruction
	pla
  	rts
lcd_home:; move cursor to first row
	pha
	;jsr lcd_wait
	lda #%10000000 + $00
	jsr lcd_instruction
	pla
	rts
lcd_clear: ; clear entire LCD
	pha
	;jsr lcd_wait
	lda #$00000001 ; Clear display
  	jsr lcd_instruction
	pla
	rts

do_strout: ; output string, address of text hi: A, lo: X
	phy
    stx TEMP_VEC
	sta TEMP_VEC+1
  	ldy #0
strprint:
  	lda (TEMP_VEC),y
  	beq strreturn
  	jsr do_chrout
  	iny
  	jmp strprint
strreturn:
	ply
	rts

do_chrout: ; output a single char to LCD, char in A
	jsr lcd_wait
	pha
	; sending high nibble
	lsr
	lsr
	lsr
	lsr
	ora #(LCD_RS | LCD_E)
	sta VIA_ORB
	eor #LCD_E
	sta VIA_ORB

	pla 
	and #$0F
	ora #(LCD_RS | LCD_E)
	sta VIA_ORB
	eor #LCD_E
	sta VIA_ORB
	rts

;------------------------------------------------------------------------------
; provides about 100uS delay for each OUTER loop. 
; The count of outloops will be used from A.
; for 1MHz we had a cycle with 1us. if A = 1 we had 20 + 20 clks, which means a minimum of 200us
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
do_nmi: ; nmi service routine
	pha
	phx
	phy
	; look if an external nmi service routine is set
	lda #$00
	cmp >NMI_SRV
	beq isr_end
	jmp (NMI_SRV)
	 
do_irq: ; irq service routine
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
	; look if an external irq service routine is set
	lda #$00
	cmp >IRQ_SRV
	beq isr_end
	jmp (IRQ_SRV)
isr_end: ; this is the ending for all interrupt service routines
	ply
	plx
	pla
	rti


;----- Messages of the bios -----
	message_w6502sbc: .asciiz "W6502SBC Welcome"
	message_ramtas: .asciiz "W6502SBC RAMTAS"
	message_srvinit: .asciiz "W6502SBC SRV INIT"
	message_ready: .asciiz "W6502SBC ready"
	message_britta: .asciiz "Hallo Britta"

;----- jump table for bios routines -----
jump_table: 

	.org $FF00 ; STROUT output string, A = high, X = low
	jmp do_strout

	.org $FF81 ; SCINIT Initialize "Screen output", (here only the serial monitor)
	jmp do_scinit
	
	.org $FF84 ; IOINIT Initialize VIA & IRQ
	jmp do_ioinit

	.org $FF87 ; RAMTAS RAM test and search memory end
	jmp do_ramtas

	.org $FF8A ; RESTOR restore default kernel vectors
	jmp do_restor 

	.org $FFD2 ; CHROUT Output an character
	jmp do_chrout

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

;----- cpu vectors -----
	.org  $FFFA
	.word   do_nmi
	.word   do_reset
	.word   do_irq
