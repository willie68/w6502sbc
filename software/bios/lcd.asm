; constants for LCD
LCD_E  .equ %00010000
LCD_RW .equ %01000000
LCD_RS .equ %00100000

; ---- Macros for lc display
.macro LcdOutput(col)
	lda #$01
	ldx #col
	jsr lcd_goto 

	lda #col
	jsr do_nhexout
.endmacro

.macro LcdMsgOut(msg)
	lda #>message_w6502sbc
	ldx #<message_w6502sbc
	jsr lcd_strout 
	lda #" "
	jsr lcd_chrout
	lda #>msg
	ldx #<msg
	jsr lcd_strout 
.endmacro
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
	jsr lcd_wait
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
lcd_goto: ; move cursor to A row, X Column
	pha
  	cmp #$00
	beq @docol
  	lda #$3F
@docol:
	stx TEMPBYTE 
	adc TEMPBYTE
	ora #%10000000	
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

lcd_strout: ; output string, address of text hi: A, lo: X
	phy
	phx
	pha
    stx TEMP_VEC
	sta TEMP_VEC+1
  	ldy #0
@lcdstrprint:
  	lda (TEMP_VEC),y
  	beq @lcdstrreturn
  	jsr lcd_chrout
  	iny
  	jmp @lcdstrprint
@lcdstrreturn:
	pla
	plx
	ply
	rts

lcd_chrout: ; output a single char to LCD, char in A
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

lcd_numout: ; output a number in decimal, A: hi byte, X lo byte
            ; Output 16-bit unsigned integer to stdout
            ; by Michael T. Barry 2017.07.07. Free to
            ; copy, use and modify, but without warranty
@iout:
    stx LNIBBLE ; low-order half
    sta HNIBBLE ; high-order half
    lda #0 ; null delimiter for print
    pha ; repeat {
@iout2: ; divide by 10
    lda #0 ; remainder
    ldx #16 ; loop counter
@iout3:
    cmp #5 ; partial remainder >= 10 (/2)?
    bcc @iout4
    sbc #5 ; yes: update partial
; remainder, set carry
@iout4:
    rol LNIBBLE ; gradually replace dividend
    rol HNIBBLE ; with the quotient
    rol ; A is gradually replaced
    dex ; with the remainder
    bne @iout3 ; loop 16 times
    ora #$30 ; convert remainder to ASCII
    pha ; stack digits in ascending
    lda LNIBBLE ; order ('0' for zero)
    ora HNIBBLE
    bne @iout2 ; } until quotient is 0
    pla
@iout5:
    jsr lcd_chrout ; print digits in descending
    pla ; order until delimiter is
    bne @iout5 ; encountered
    rts 
lcd_hexout: ; output a number in hex with leading $, A: hi byte, X lo byte 
    pha
    pha
    lda #"$"
    jsr lcd_chrout
    pla             ; output hi byte
    jsr lcd_bhexout
    txa             ; output hi byte
    jsr lcd_bhexout
    pla
    rts
lcd_bhexout:
    pha             ; output hi nibble
    lsr
	lsr
	lsr
	lsr
	jsr lcd_nhexout
    pla             ; output lo nibble of hi byte
	jsr lcd_nhexout
    rts

lcd_nhexout: ; output the lower nibble as hex
    pha             ; output hi nibble
    and #$0F
    ora #$30        ; add "0"
	cmp #$3A
	bmi @bh2
	adc #$06
@bh2:
    jsr lcd_chrout
    pla
    rts
