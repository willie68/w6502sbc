/* SPI CONNECTION: 6522 VIA PORT A to MAX3100
   ________                             ______
           |                           |
        PA7|<----- RECEIVE (NEW) --+--<|DOUT
        PA6|                           |
  6522  PA5|                           | MAX3100
  VIA   PA4|                           |  UART
        PA3|                           |
        PA2|>------ CHIP SELECT ------>|/CS
        PA1|>------- TRANSMIT -------->|DIN
        PA0|>--------- CLOCK --------->|SCLK
           |                           |
        CA1|<------- INTERRUPT --------|/IRQ
        CA2|                           |
   ________|                           |______

THE 6522 PARALLEL PORT CONTAINS 8 BIDIRECTIONAL DATA
PINS (PA0 TO PA7) AND 2 "CONTROL" PINS (CA1 AND CA2). 
ARROWS INDICATE THE DIRECTIONS WE ARE USING TO SEND
AND RECEIVE DATA.
Inspired by coronax SPI Implementation (coronax.wordpress.com)

*/
do_serinit:
    pha
    ; VIA setup for SPI max3100
	lda #%01111111  ; set every port pin to output except the recive PA7
	sta VIA_DDRA
	lda VIA_PCR
    and #%11111110  ; activate /IRQ on CA1
	sta VIA_PCR
    pla
    rts

;; spibyte
;; Sends the byte in accumulator and receives a
;; new byte into accumulator.
spibyte:
	sta SPI_WRITEB
.loop 8			        ; copy the next section 8 times
	lda #%01111000		; base DATAB value with chip select for
				        ; MAX3100 and a zero bit in the output
				        ; line.
	rol SPI_WRITEB
	bcc @writing_zero_bit
	ora #%00000010		; write a 1 bit to the output line.
@writing_zero_bit:
		
 	sta VIA_ORA    	    ; write data back to the port
	inc VIA_ORA    	    ; set clock high

	lda VIA_ORA		    ; Read input bit
	rol			        ; Shift input bit to carry flag
	rol SPI_READB	    ; Shift carry into readbuffer

	dec VIA_ORA    	    ; set clock low
.endloop
	lda SPI_READB	    ; result goes in A
	rts

do_serout:
    pha                 ;Save ACCUMULATOR on STACK
    pla
    rts

do_serin:
    pha                 ;Save ACCUMULATOR on STACK
    pla
    rts

do_serstrout:           ; output string, address of text hi: A, lo: X
	phy
    stx TEMP_VEC
	sta TEMP_VEC+1
  	ldy #0
@serprint:
  	lda (TEMP_VEC),y
  	beq @serreturn
  	jsr do_serout
  	iny
  	jmp @serprint
@serreturn:
	ply
	rts
