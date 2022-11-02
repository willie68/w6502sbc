;----- macros -----
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
