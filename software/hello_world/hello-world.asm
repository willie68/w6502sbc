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
  sta VIA_ORB
  lda #0         ; Clear RS/RW/E bits
  sta VIA_ORA
  lda #E         ; Set E bit to send instruction
  sta VIA_ORA
  lda #0         ; Clear RS/RW/E bits
  sta VIA_ORA

  lda #%00000001 ; clear Display;
  sta VIA_ORB
  lda #0         ; Clear RS/RW/E bits
  sta VIA_ORA
  lda #E         ; Set E bit to send instruction
  sta VIA_ORA
  lda #0         ; Clear RS/RW/E bits
  sta VIA_ORA

  lda #%00000010 ; return home;
  sta VIA_ORB
  lda #0         ; Clear RS/RW/E bits
  sta VIA_ORA
  lda #E         ; Set E bit to send instruction
  sta VIA_ORA
  lda #0         ; Clear RS/RW/E bits
  sta VIA_ORA

  lda #%00001110 ; Display on; cursor on; blink off
  sta VIA_ORB
  lda #0         ; Clear RS/RW/E bits
  sta VIA_ORA
  lda #E         ; Set E bit to send instruction
  sta VIA_ORA
  lda #0         ; Clear RS/RW/E bits
  sta VIA_ORA

  lda #%00000110 ; Increment and shift cursor; don't shift display
  sta VIA_ORB
  lda #0         ; Clear RS/RW/E bits
  sta VIA_ORA
  lda #E         ; Set E bit to send instruction
  sta VIA_ORA
  lda #0         ; Clear RS/RW/E bits
  sta VIA_ORA

  lda #"H"
  sta VIA_ORB
  lda #RS         ; Set RS; Clear RW/E bits
  sta VIA_ORA
  lda #(RS | E)   ; Set E bit to send instruction
  sta VIA_ORA
  lda #RS         ; Clear E bits
  sta VIA_ORA

  lda #"e"
  sta VIA_ORB
  lda #RS         ; Set RS; Clear RW/E bits
  sta VIA_ORA
  lda #(RS | E)   ; Set E bit to send instruction
  sta VIA_ORA
  lda #RS         ; Clear E bits
  sta VIA_ORA

  lda #"l"
  sta VIA_ORB
  lda #RS         ; Set RS; Clear RW/E bits
  sta VIA_ORA
  lda #(RS | E)   ; Set E bit to send instruction
  sta VIA_ORA
  lda #RS         ; Clear E bits
  sta VIA_ORA

  lda #"l"
  sta VIA_ORB
  lda #RS         ; Set RS; Clear RW/E bits
  sta VIA_ORA
  lda #(RS | E)   ; Set E bit to send instruction
  sta VIA_ORA
  lda #RS         ; Clear E bits
  sta VIA_ORA

  lda #"o"
  sta VIA_ORB
  lda #RS         ; Set RS; Clear RW/E bits
  sta VIA_ORA
  lda #(RS | E)   ; Set E bit to send instruction
  sta VIA_ORA
  lda #RS         ; Clear E bits
  sta VIA_ORA

  lda #","
  sta VIA_ORB
  lda #RS         ; Set RS; Clear RW/E bits
  sta VIA_ORA
  lda #(RS | E)   ; Set E bit to send instruction
  sta VIA_ORA
  lda #RS         ; Clear E bits
  sta VIA_ORA

  lda #" "
  sta VIA_ORB
  lda #RS         ; Set RS; Clear RW/E bits
  sta VIA_ORA
  lda #(RS | E)   ; Set E bit to send instruction
  sta VIA_ORA
  lda #RS         ; Clear E bits
  sta VIA_ORA

  lda #"w"
  sta VIA_ORB
  lda #RS         ; Set RS; Clear RW/E bits
  sta VIA_ORA
  lda #(RS | E)   ; Set E bit to send instruction
  sta VIA_ORA
  lda #RS         ; Clear E bits
  sta VIA_ORA

  lda #"o"
  sta VIA_ORB
  lda #RS         ; Set RS; Clear RW/E bits
  sta VIA_ORA
  lda #(RS | E)   ; Set E bit to send instruction
  sta VIA_ORA
  lda #RS         ; Clear E bits
  sta VIA_ORA

  lda #"r"
  sta VIA_ORB
  lda #RS         ; Set RS; Clear RW/E bits
  sta VIA_ORA
  lda #(RS | E)   ; Set E bit to send instruction
  sta VIA_ORA
  lda #RS         ; Clear E bits
  sta VIA_ORA

  lda #"l"
  sta VIA_ORB
  lda #RS         ; Set RS; Clear RW/E bits
  sta VIA_ORA
  lda #(RS | E)   ; Set E bit to send instruction
  sta VIA_ORA
  lda #RS         ; Clear E bits
  sta VIA_ORA

  lda #"d"
  sta VIA_ORB
  lda #RS         ; Set RS; Clear RW/E bits
  sta VIA_ORA
  lda #(RS | E)   ; Set E bit to send instruction
  sta VIA_ORA
  lda #RS         ; Clear E bits
  sta VIA_ORA

  lda #"!"
  sta VIA_ORB
  lda #RS         ; Set RS; Clear RW/E bits
  sta VIA_ORA
  lda #(RS | E)   ; Set E bit to send instruction
  sta VIA_ORA
  lda #RS         ; Clear E bits
  sta VIA_ORA

main_loop:
  jmp main_loop

do_nmi: NOP
		RTI
	 
do_irq: NOP
		RTI

	.org  $FFFA
	.word   do_nmi
	.word   do_reset
	.word   do_irq
