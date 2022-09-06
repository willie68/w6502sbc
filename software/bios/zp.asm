; ZERO Page registers $0000.. $00ff
COUNTER .equ $20 ; counter for different things
HNIBBLE .equ $21
LNIBBLE .equ $22
TEMPBYTE .equ $23
RAMTOP .equ $30 ; store the page of the last RAM ($30 is the low adress)
TEMP_VEC .equ $32 ; store a temporary vector, like the address to the string to output, $32 low, $33 hi

IN_READ .equ $80
IN_WRITE .equ $81

JTIME .equ $A0 ; to $A2 three bytes jiffy time, higher two bytes of the 3 bytes of the 1/50 secs of a day. 24h * 60m * 60s * 50, 4.320.000 ticks per day

SPI_WRITEB .equ $B0 ; write single byte buffer for SPI
SPI_READB .equ $B1 ; read single byte buffer for SPI