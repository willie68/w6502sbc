package emulator

import "fmt"

func lda_direct(e *Emu6502) string {
	v := e.getMnemonic()
	e.setFlags(v, nil, nil)
	e.a = v
	return fmt.Sprintf("%.2x      lda #$%.2x", e.a, e.a)
}

func lda_abs(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	e.setFlags(v, nil, nil)
	e.a = v
	return str + fmt.Sprintf("   lda $%.4x", adr)
}

func lda_abs_x(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	e.setFlags(v, nil, nil)
	e.a = v
	return str + fmt.Sprintf("   lda $%.4x,X", adr)
}

func lda_abs_y(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.y))
	e.setFlags(v, nil, nil)
	e.a = v
	return str + fmt.Sprintf("   lda $%.4x,Y", adr)
}

func lda_zp(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	e.setFlags(v, nil, nil)
	e.a = v
	return str + fmt.Sprintf("   lda $%.2x", adr)
}

func lda_zp_x(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	e.setFlags(v, nil, nil)
	e.a = v
	return str + fmt.Sprintf("   lda $%.2x,X", adr)
}

func lda_ind_x(e *Emu6502) string {
	zp, str := e.getZPAddress()
	zpx := zp + uint16(e.x)
	adr := e.readVector(zpx)
	v := e.getMemory(adr)
	e.setFlags(v, nil, nil)
	e.a = v
	return str + fmt.Sprintf("        lda ($%.2x,X)", adr)
}

func lda_ind_y(e *Emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	v := e.getMemory(adr + uint16(e.y))
	e.setFlags(v, nil, nil)
	e.a = v
	return str + fmt.Sprintf("        lda ($%.2x),Y", adr)
}

func ldx_direct(e *Emu6502) string {
	v := e.getMnemonic()
	e.setFlags(v, nil, nil)
	e.x = v
	return fmt.Sprintf("%.2x        ldx #$%.2x", e.x, e.x)
}

func ldx_abs(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	e.setFlags(v, nil, nil)
	e.x = v
	return str + fmt.Sprintf("   ldx $%.4x", adr)
}

func ldx_abs_y(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.y))
	e.setFlags(v, nil, nil)
	e.x = v
	return str + fmt.Sprintf("   ldx $%.4x,Y", adr)
}

func ldx_zp(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	e.setFlags(v, nil, nil)
	e.x = v
	return str + fmt.Sprintf("   ldx $%.2x", adr)
}

func ldx_zp_y(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.y))
	e.setFlags(v, nil, nil)
	e.x = v
	return str + fmt.Sprintf("   ldx $%.2x,Y", adr)
}

func ldy_direct(e *Emu6502) string {
	v := e.getMnemonic()
	e.setFlags(v, nil, nil)
	e.y = v
	return fmt.Sprintf("%.2x      ldy #$%.2x", e.a, e.a)
}

func ldy_abs(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	e.setFlags(v, nil, nil)
	e.y = v
	return str + fmt.Sprintf("   ldy $%.4x", adr)
}

func ldy_abs_x(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	e.setFlags(v, nil, nil)
	e.y = v
	return str + fmt.Sprintf("   ldy $%.4x,X", adr)
}

func ldy_zp(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	e.setFlags(v, nil, nil)
	e.y = v
	return str + fmt.Sprintf("   ldy $%.2x", adr)
}

func ldy_zp_x(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	e.setFlags(v, nil, nil)
	e.y = v
	return str + fmt.Sprintf("   ldy $%.2x,X", adr)
}

func sta_abs(e *Emu6502) string {
	adr, str := e.getAddress()
	e.setMemory(adr, e.a)
	return str + fmt.Sprintf("     sta $%.4x", adr)
}

func sta_abs_x(e *Emu6502) string {
	adr, str := e.getAddress()
	e.setMemory(adr+uint16(e.x), e.a)
	return str + fmt.Sprintf("     sta $%.4x,X", adr)
}

func sta_abs_y(e *Emu6502) string {
	adr, str := e.getAddress()
	e.setMemory(adr+uint16(e.y), e.a)
	return str + fmt.Sprintf("     sta $%.4x,Y", adr)
}

func sta_zp(e *Emu6502) string {
	adr, str := e.getZPAddress()
	e.setMemory(adr, e.a)
	return str + fmt.Sprintf("     sta $%.2x", adr)
}

func sta_zp_x(e *Emu6502) string {
	adr, str := e.getZPAddress()
	e.setMemory(adr+uint16(e.x), e.a)
	return str + fmt.Sprintf("    sta $%.2x,X", adr)
}

func sta_ind_x(e *Emu6502) string {
	zp, str := e.getZPAddress()
	zpx := zp + uint16(e.x)
	adr := e.readVector(zpx)
	e.setMemory(adr, e.a)
	return str + fmt.Sprintf("     sta ($%.2x,X)", adr)
}

func sta_ind_y(e *Emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	e.setMemory(adr+uint16(e.y), e.a)
	return str + fmt.Sprintf("     sta ($%.2x),Y", adr)
}

func stx_abs(e *Emu6502) string {
	adr, str := e.getAddress()
	e.setMemory(adr, e.x)
	return str + fmt.Sprintf("     stx $%.4x", adr)
}

func stx_zp(e *Emu6502) string {
	adr, str := e.getZPAddress()
	e.setMemory(adr, e.x)
	return str + fmt.Sprintf("     stx $%.2x", adr)
}

func stx_zp_y(e *Emu6502) string {
	adr, str := e.getZPAddress()
	e.setMemory(adr+uint16(e.y), e.x)
	return str + fmt.Sprintf("    stx $%.2x,Y", adr)
}

func sty_abs(e *Emu6502) string {
	adr, str := e.getAddress()
	e.setMemory(adr, e.y)
	return str + fmt.Sprintf("     sty $%.4x", adr)
}

func sty_zp(e *Emu6502) string {
	adr, str := e.getZPAddress()
	e.setMemory(adr, e.y)
	return str + fmt.Sprintf("     sty $%.2x", adr)
}

func sty_zp_x(e *Emu6502) string {
	adr, str := e.getZPAddress()
	e.setMemory(adr+uint16(e.x), e.y)
	return str + fmt.Sprintf("    sty $%.2x,X", adr)
}

func tax(e *Emu6502) string {
	e.x = e.a
	e.setFlags(e.a, nil, nil)
	return "        tax"
}

func tay(e *Emu6502) string {
	e.y = e.a
	e.setFlags(e.a, nil, nil)
	return "        tay"
}

func txa(e *Emu6502) string {
	e.a = e.x
	e.setFlags(e.a, nil, nil)
	return "        txa"
}

func tya(e *Emu6502) string {
	e.a = e.y
	e.setFlags(e.a, nil, nil)
	return "        tya"
}

func tsx(e *Emu6502) string {
	e.x = e.sp
	e.setFlags(e.x, nil, nil)
	return "        tsx"
}

func txs(e *Emu6502) string {
	e.sp = e.x
	e.setFlags(e.x, nil, nil)
	return "        txs"
}
