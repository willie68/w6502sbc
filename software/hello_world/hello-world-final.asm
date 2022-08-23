.format "bin"
.target "65C02"

.memory "fill", $E000, $2000, $ea
.org $E000
.include "io.asm" 
	
	E  = %10000000
	RW = %01000000
	RS = %00100000

do_reset:
  ldx #$ff ;init stack pointer
  txs 

  lda #%11111111 ; Set all pins on port B to output
  sta VIA_DDRB
  lda #%11100000 ; Set top 3 pins on port A to output
  sta VIA_DDRA

  lda #%00111000 ; Set 8-bit mode; 2-line display; 5x8 font
  jsr lcd_instruction
  lda #%00001110 ; Display on; cursor on; blink off
  jsr lcd_instruction
  lda #%00000110 ; Increment and shift cursor; don't shift display
  jsr lcd_instruction
  lda #$00000010 ; Return home
  jsr lcd_instruction
  lda #$00000001 ; Clear display
  jsr lcd_instruction

  ldx #0
print:
  lda message,x
  beq main_loop
  jsr print_char
  inx
  jmp print

main_loop:
  jmp main_loop

message: .asciiz "Hello, willie!"

lcd_wait:
  pha
  lda #%00000000  ; Port B is input
  sta VIA_DDRB
lcdbusy:
  lda #RW
  sta VIA_ORA
  lda #(RW | E)
  sta VIA_ORA
  lda VIA_ORB
  and #%10000000
  bne lcdbusy

  lda #RW
  sta VIA_ORA
  lda #%11111111  ; Port B is output
  sta VIA_DDRB
  pla
  rts

lcd_instruction:
  jsr lcd_wait
  sta VIA_ORB
  lda #0         ; Clear RS/RW/E bits
  sta VIA_ORA
  lda #E         ; Set E bit to send instruction
  sta VIA_ORA
  lda #0         ; Clear RS/RW/E bits
  sta VIA_ORA
  rts

print_char:
  jsr lcd_wait
  sta VIA_ORB
  lda #RS         ; Set RS; Clear RW/E bits
  sta VIA_ORA
  lda #(RS | E)   ; Set E bit to send instruction
  sta VIA_ORA
  lda #RS         ; Clear E bits
  sta VIA_ORA
  rts

do_nmi: NOP
		RTI
	 
do_irq: NOP
		RTI

	.org  $FFFA
	.word   do_nmi
	.word   do_reset
	.word   do_irq
