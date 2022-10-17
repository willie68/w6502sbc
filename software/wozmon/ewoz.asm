;EWoz 1.0
;by fsafstrom » Mar Wed 14, 2007 12:23 pm
;http://www.brielcomputers.com/phpBB3/viewtopic.php?f=9&t=197#p888
;via http://jefftranter.blogspot.co.uk/2012/05/woz-mon.html
;
;The EWoz 1.0 is just the good old Woz mon with a few improvements and extensions so to say. 
;
;It's using ACIA @ 19200 Baud. 
;It prints a small welcome message when started. 
;All key strokes are converted to uppercase. 
;The backspace works so the _ is no longer needed. 
;When you run a program, it's called with an jsr so if the program ends with an rts, you will be taken back to the monitor. 
;You can load Intel HEX format files and it keeps track of the checksum. 
;To load an Intel Hex file, just type L and hit return. 
;Now just send a Text file that is in the Intel HEX Format just as you would send a text file for the Woz mon. 
;You can abort the transfer by hitting ESC. 
;
;The reason for implementing a loader for HEX files is the 6502 Assembler @ http://home.pacbell.net/michal_k/6502.html 
;This assembler saves the code as Intel HEX format. 
;
;In the future I might implement XModem, that is if anyone would have any use for it...  
;
;Enjoy...
; ... eleminating binary ...
;
; EWOZ Extended Woz Monitor.
; Just a few mods to the original monitor.

.format "bin"
.target "65C02"
.include "io.asm" 

IN          .equ $0200          ;*Input buffer
XAML        .equ $24            ;*Index pointers
XAMH        .equ $25
STL         .equ $26
STH         .equ $27
L           .equ $28
H           .equ $29
YSAV        .equ $2A
MODE        .equ $2B
MSGL        .equ $2C
MSGH        .equ $2D
COUNTER     .equ $2E
CRC         .equ $2F
CRCCHECK    .equ $30

.org $E000

do_reset:   cld             ;Clear decimal arithmetic mode.
            cli
            lda #$1F        ;* Init ACIA to 19200 Baud.
            sta ACIA_CONTROL
            lda #$0B        ;* No Parity.
            sta ACIA_COMMAND
            lda #$0D
            jsr echo        ;* New line.
            lda #<msg1
            sta MSGL
            lda #>msg1
            sta MSGH
            jsr shwmsg      ;* Show Welcome.
            lda #$0D
            jsr echo        ;* New line.
softreset:  lda #$9B        ;* Auto escape.
notcr:      cmp #$88        ;"<-"? * Note this was chaged to $88 which is the back space key.
            beq backspace   ;Yes.
            cmp #$9B        ;ESC?
            beq escape      ;Yes.
            iny             ;Advance text index.
            bpl nextchar    ;Auto ESC if >127.
escape:     lda #$DC        ;"\"
            jsr echo        ;Output it.
getline:    lda #$8D        ;CR.
            jsr echo        ;Output it.
            ldy #$01        ;Initiallize text index.
backspace:  dey             ;Backup text index.
            bmi getline     ;Beyond start of line, reinitialize.
            lda #$A0        ;*Space, overwrite the backspaced char.
            jsr echo
            lda #$88        ;*backspace again to get to correct pos.
            jsr echo
nextchar:   lda ACIA_STATUS ;*See if we got an incoming char
            and #$08        ;*Test bit 3
            beq nextchar    ;*wait for character
            lda ACIA_RX     ;*Load char
            cmp #$60        ;*Is it Lower case
            bmi convert     ;*Nope, just convert it
            and #$5F        ;*If lower case, convert to Upper case
convert:    ora #$80        ;*convert it to "ASCII Keyboard" Input
            sta IN,y        ;Add to text buffer.
            jsr echo        ;Display character.
            cmp #$8D        ;CR?
            bne notcr       ;No.
            ldy #$FF        ;Reset text index.
            lda #$00        ;For XAM mode.
            tax             ;0->X.
setstore:   asl             ;Leaves $7B if setting STOR mode.
setmode:    sta MODE        ;$00 = XAM, $7B = STOR, $AE = BLOK XAM.
blskip:     iny             ;Advance text index.
nextitem:   lda IN,y        ;Get character.
            cmp #$8D        ;CR?
            beq getline     ;Yes, done this line.
            cmp #$AE        ;"."?
            bcc blskip      ;Skip delimiter.
            beq setmode     ;Set BLOCK XAM mode.
            cmp #$BA        ;":"?
            beq setstore     ;Yes, set STOR mode.
            cmp #$D2        ;"R"?
            beq run         ;Yes, run user program.
            cmp #$CC        ;* "L"?
            beq loadint     ;* Yes, Load Intel Code.
            stx L           ;$00->L.
            stx H           ; and H.
            sty YSAV        ;Save Y for comparison.
nexthex:    lda IN,y        ;Get character for hex test.
            eor #$B0        ;Map digits to $0-9.
            cmp #$0A        ;digit?
            bcc dig         ;Yes.
            adc #$88        ;Map letter "A"-"F" to $FA-FF.
            cmp #$FA        ;Hex letter?
            bcc nothex      ;No, character not hex.
dig:        asl
            asl             ;Hex digit to MSD of A.
            asl
            asl
            ldx #$04        ;Shift count.
hexshift:   asl             ;Hex digit left MSB to carry.
            rol L           ;Rotate into LSD.
            rol H           ;Rotate into MSD's.
            dex             ;Done 4 shifts?
            bne hexshift    ;No, loop.
            iny             ;Advance text index.
            bne nexthex     ;Always taken. Check next character for hex.
nothex:     cpy YSAV        ;Check if L, H empty (no hex digits).
            bne noescape    ;* Branch out of range, had to improvise...
            jmp escape      ;Yes, generate ESC sequence.

run:        jsr actrun      ;* jsr to the Address we want to run.
            jmp softreset   ;* When returned for the program, reset EWOZ.
actrun:     jmp (XAML)      ;run at current XAM index.

loadint:    jsr loadintel   ;* Load the Intel code.
            jmp softreset   ;* When returned from the program, reset EWOZ.

noescape:   bit MODE        ;Test MODE byte.
            bvc notstore     ;B6=0 for STOR, 1 for XAM and BLOCK XAM
            lda L           ;LSD's of hex data.
            sta (STL, x)    ;Store at current "store index".
            inc STL         ;Increment store index.
            bne nextitem    ;Get next item. (no carry).
            inc STH         ;Add carry to 'store index' high order.
tonextitem: jmp nextitem    ;Get next command item.
notstore:   bmi xamnext     ;B7=0 for XAM, 1 for BLOCK XAM.
            ldx #$02        ;Byte count.
setadr:     lda L-1,x       ;Copy hex data to
            sta STL-1,x     ;"store index".
            sta XAML-1,x    ;and to "XAM index'.
            dex             ;Next of 2 bytes.
            bne setadr      ;Loop unless X = 0.
nxtprnt:    bne prdata      ;NE means no address to print.
            lda #$8D        ;CR.
            jsr echo        ;Output it.
            lda XAMH        ;'Examine index' high-order byte.
            jsr prbyte      ;Output it in hex format.
            lda XAML        ;Low-order "examine index" byte.
            jsr prbyte      ;Output it in hex format.
            lda #$BA        ;":".
            jsr echo        ;Output it.
prdata:     lda #$A0        ;Blank.
            jsr echo        ;Output it.
            lda (XAML,x)    ;Get data byte at 'examine index".
            jsr prbyte      ;Output it in hex format.
xamnext:    stx MODE        ;0-> MODE (XAM mode).
            lda XAML
            cmp L           ;Compare 'examine index" to hex data.
            lda XAMH
            sbc H
            bcs tonextitem  ;Not less, so no more data to output.
            inc XAML
            bne mod8chk     ;Increment 'examine index".
            inc XAMH
mod8chk:    lda XAML        ;Check low-order 'exainine index' byte
            and #$0F        ;For MOD 8=0 ** changed to $0F to get 16 values per row **
            bpl nxtprnt     ;Always taken.
prbyte:     pha             ;Save A for LSD.
            lsr
            lsr
            lsr             ;MSD to LSD position.
            lsr
            jsr prhex       ;Output hex digit.
            pla             ;Restore A.
prhex:      and #$0F        ;Mask LSD for hex print.
            ora #$B0        ;Add "0".
            cmp #$BA        ;digit?
            bcc echo        ;Yes, output it.
            adc #$06        ;Add offset for letter.
echo:       pha             ;*Save A
            and #$7F        ;*Change to "standard ASCII"
            sta ACIA_TX     ;*Send it.
@wait:      lda ACIA_STATUS ;*Load status register for ACIA
            and #$10        ;*Mask bit 4.
            beq @wait       ;*ACIA not done yet, wait.
            pla             ;*Restore A
            rts             ;*Done, over and out...

shwmsg:     ldy #$0
@print:     lda (MSGL),y
            beq @DONE
            jsr echo
            iny 
            bne @print
@DONE:      rts 


; Load an program in Intel Hex Format.
loadintel:  lda #$0D
            jsr echo        ;New line.
            lda #<msg2
            sta MSGL
            lda #>msg2
            sta MSGH
            jsr shwmsg      ;Show start Transfer.
            lda #$0D
            jsr echo        ;New line.
            ldy #$00
            sty CRCCHECK   ;If CRCCHECK=0, all is good.
intelline:  jsr getchar    ;Get char
            sta IN,y       ;Store it
            iny            ;Next
            cmp   #$1B     ;escape ?
            beq inteldone  ;Yes, abort.
            cmp #$0D       ;Did we find a new line ?
            bne intelline  ;Nope, continue to scan line.
            ldy #$FF       ;Find (:)
findcol:    iny
            lda IN,y
            cmp #$3A       ; Is it Colon ?
            bne findcol    ; Nope, try next.
            iny            ; Skip colon
            ldx   #$00     ; Zero in X
            stx   CRC      ; Zero Check sum
            jsr gethex     ; Get Number of bytes.
            sta COUNTER    ; Number of bytes in Counter.
            clc            ; Clear carry
            adc CRC        ; Add CRC
            sta CRC        ; Store it
            jsr gethex     ; Get Hi byte
            sta STH        ; Store it
            clc            ; Clear carry
            adc CRC        ; Add CRC
            sta CRC        ; Store it
            jsr gethex     ; Get Lo byte
            sta STL        ; Store it
            clc            ; Clear carry
            adc CRC        ; Add CRC
            sta CRC        ; Store it
            lda #$2E       ; Load "."
            jsr echo       ; print it to indicate activity.
nodot:      jsr gethex     ; Get Control byte.
            cmp   #$01     ; Is it a Termination record ?
            beq inteldone  ; Yes, we are done.
            clc            ; Clear carry
            adc CRC        ; Add CRC
            sta CRC        ; Store it
intelstore: jsr gethex     ; Get Data Byte
            sta (STL,x)    ; Store it
            clc            ; Clear carry
            adc CRC        ; Add CRC
            sta CRC        ; Store it
            inc STL        ; Next Address
            bne testcount  ; Test to see if Hi byte needs inc
            inc STH        ; If so, inc it.
testcount:  dec COUNTER    ; Count down.
            bne intelstore ; Next byte
            jsr gethex     ; Get Checksum
            ldy #$00       ; Zero Y
            clc            ; Clear carry
            adc CRC        ; Add CRC
            beq intelline  ; Checksum OK.
            lda #$01       ; Flag CRC error.
            sta CRCCHECK   ; Store it
            jmp intelline  ; Process next line.

inteldone:  lda CRCCHECK   ; Test if everything is OK.
            beq okmess     ; Show OK message.
            lda #$0D
            jsr echo       ;New line.
            lda #<msg4     ; Load Error Message
            sta MSGL
            lda #>msg4
            sta MSGH
            jsr shwmsg     ;Show Error.
            lda #$0D
            jsr echo       ;New line.
            rts

okmess:     lda #$0D
            jsr echo       ;New line.
            lda #<msg3     ;Load OK Message.
            sta MSGL
            lda #>msg3
            sta MSGH
            jsr shwmsg     ;Show Done.
            lda #$0D
            jsr echo       ;New line.
            rts

gethex:     lda IN,y       ;Get first char.
            eor #$30
            cmp #$0A
            bcc donefirst
            adc #$08
donefirst:  asl
            asl
            asl
            asl
            sta L
            iny
            lda IN,y       ;Get next char.
            eor #$30
            cmp #$0A
            bcc donesecond
            adc #$08
donesecond: and #$0F
            ora L
            iny
            rts

getchar:    lda ACIA_STATUS ;See if we got an incoming char
            and #$08        ;Test bit 3
            beq getchar     ;wait for character
            lda ACIA_RX     ;Load char
            rts

msg1:  .asciiz "Welcome to EWOZ 1.0."
msg2:  .asciiz "start Intel Hex code Transfer."
msg3:  .asciiz "Intel Hex Imported OK."
msg4:  .asciiz "Intel Hex Imported with checksum error."

do_nmi: ; nmi service routine 
do_irq: ; irq service routine
	rti

;----- cpu vectors -----
	.org  $FFFA
	.word   do_nmi
	.word   do_reset
	.word   do_irq