; constants for LCD
LCD_E  .equ %10000000
LCD_RW .equ %01000000
LCD_RS .equ %00100000

; ---- Display routines ----
do_scinit: 		; initialise LC-Display on port B
	; D4..D7 on Port pins PB0..3
	; RS; R/W and E on Port pins PB5, PB6, PB7
	lda #$ff 	; Set all pins on port B to output
  	sta VIA_DDRB
	stz VIA_ORB

	; reset the display, wait at least 15ms before sending something to the lcd
	lda #$58
	jsr do_delay

	; send 3 times the reset...
  	lda #(%00000011 | LCD_E) ; 1. RESET
	sta VIA_ORB
	eor #LCD_E
	sta VIA_ORB
    lda #$1f                 ; wait min 4.1ms
	jsr do_delay

  	lda #(%00000011 | LCD_E) ; 2. RESET
	sta VIA_ORB
	eor #LCD_E
	sta VIA_ORB
	lda #$01                 ; wait min 100us
	jsr do_delay

  	lda #(%00000011 | LCD_E) ; 3. RESET min 100us
	sta VIA_ORB
	eor #LCD_E
	sta VIA_ORB
	lda #$01                 ; wait min 100us
	jsr do_delay

  	lda #(%00000010 | LCD_E) ; Set 4-bit mode; 
	sta VIA_ORB
	eor #LCD_E
	sta VIA_ORB
	lda #$01                 ; wait min 37us
	jsr do_delay

	; after this command we can use the 4-Bit mode and we could use busy flag for former sync
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
	jsr lcd_busy
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
