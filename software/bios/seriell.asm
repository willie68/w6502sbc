; seriell W65C51 adapter lib
;----- constants
CHR_DELAY .equ CLOCK / 125000
;----- macros
.macro SerMsgOut(msg)
	lda #>message_w6502sbc
	ldx #<message_w6502sbc
	jsr ser_strout 
	lda #" "
	jsr ser_chrout
	lda #>msg
	ldx #<msg
	jsr ser_strout
    SerCROut() 
.endmacro

.macro SerChrOut()
    sta  ACIA_TX        ;Write byte to ACIA transmit data register
@SEROUTL1:
    lda ACIA_STATUS    ;Read ACIA status register
    and #$10           ;Isolate transmit data register status bit
    beq @SEROUTL1      ;LOOP back to COUTL IF transmit data register is full
    lda #CHR_DELAY
    jsr do_delay
.endmacro

.macro SerCROut()
    pha
    lda #$0D
    SerChrOut()
    lda #$0A
    SerChrOut()
    pla
.endm

.macro SerStrOut(msg)
	lda #>msg
	ldx #<msg
	jsr ser_strout
.endmacro
;----- subroutines
do_serinit:
    pha
    ; ACIA setup
    stz ACIA_STATUS
    lda #%00011110			; 8-bit, 1 Stop bit, Baudrate, 9600 Baud
    sta ACIA_CONTROL
    lda #%00001011			; no parity, no echo, transmit irq disabled, no receiver irq, DTR High
    sta ACIA_COMMAND
    pla
    rts

ser_crout:
    SerCROut()
    rts
ser_chrout:
    pha                 ;Save ACCUMULATOR on STACK
    SerChrOut()
    pla
    rts

ser_chrin:
    lda #$08
@SERINL1:
    bit ACIA_STATUS     ; Check to see if the buffer is full
    beq @SERINL1
    lda ACIA_RX
    rts

ser_strout: ; output string, address of text hi: A, lo: X
	phy
    phx
    pha
    stx TEMP_VEC
	sta TEMP_VEC+1
  	ldy #0
@serprint:
  	lda (TEMP_VEC),y
  	beq @serreturn
    SerChrOut()
  	iny
  	jmp @serprint
@serreturn:
    pla
    plx
	ply
	rts
