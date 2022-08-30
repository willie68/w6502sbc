;----- macros -----
.macro msg_out(msg)
	lda #>message_w6502sbc
	ldx #<message_w6502sbc
	jsr do_strout 
	lda #" "
	jsr do_chrout
	lda #>msg
	ldx #<msg
	jsr do_strout 
.endmacro

.macro swn() 
    asl
    adc  #$80
    rol
    asl
    adc  #$80
    rol  
.endmacro

.macro toggle_a()
	pha
	lda #$ff
	sta VIA_ORA
	stz VIA_ORA
	pla
.endmacro
