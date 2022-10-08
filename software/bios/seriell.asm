do_serinit:
    pha
    ; ACIA setup
    lda #$00
    sta ACIA_STATUS
    lda #%00001011			; no parity, no echo, transmit irq disabled, no receiver irq, DTR High
    sta ACIA_COMMAND
    lda #%00011110			; 8-bit, 1 Stop bit, Baudrate, 9600 Baud
    sta ACIA_CONTROL
    pla
    rts

do_serout:
    pha                 ;Save ACCUMULATOR on STACK
    pha
@SEROUTL1:
    lda  ACIA_STATUS    ;Read ACIA status register
    and  #$10           ;Isolate transmit data register status bit
    beq  @SEROUTL1      ;LOOP back to COUTL IF transmit data register is full
    pla                 ;ELSE, restore ACCUMULATOR from STACK
    sta  ACIA_TX        ;Write byte to ACIA transmit data register
    lda #$01
    jsr do_delay
    pla
    rts

do_serin:
    lda #$08
@SERINL1:
    bit ACIA_STATUS     ; Check to see if the buffer is full
    beq @SERINL1
    lda ACIA_RX
    rts

do_serstrout: ; output string, address of text hi: A, lo: X
	phy
    stx TEMP_VEC
	sta TEMP_VEC+1
  	ldy #$00
@serprint:
  	lda (TEMP_VEC),y
  	beq @serreturn
  	jsr do_serout
  	iny
  	jmp @serprint
@serreturn:
	ply
	rts
