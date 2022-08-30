do_serout:
    pha                 ;Save ACCUMULATOR on STACK
@SEROUTL1:
    lda  ACIA_STATUS    ;Read ACIA status register
    and  #$10           ;Isolate transmit data register status bit
    beq  @SEROUTL1      ;LOOP back to COUTL IF transmit data register is full
    pla                 ;ELSE, restore ACCUMULATOR from STACK
    sta  ACIA_TX        ;Write byte to ACIA transmit data register
    rts
