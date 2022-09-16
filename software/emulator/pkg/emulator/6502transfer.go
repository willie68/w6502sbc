package emulator

import "fmt"

func lda_direct(e *emu6502) string {
	e.a = e.getMnemonic()
	e.zf = e.a == 0
	e.nf = (e.a & 0x80) > 0
	return fmt.Sprintf("%.2x        lda #$%.2x", e.a, e.a)
}

func lda_abs(e *emu6502) string {
	lo := e.getMnemonic()
	hi := e.getMnemonic()
	adr := uint16(hi)*256 + uint16(lo)
	e.a = e.getMemory(adr)
	e.zf = e.a == 0
	e.nf = (e.a & 0x80) > 0
	return fmt.Sprintf("%.2x %.2x   lda $%.4x", lo, hi, adr)
}

func lda_abs_x(e *emu6502) string {
	lo := e.getMnemonic()
	hi := e.getMnemonic()
	adr := uint16(hi)*256 + uint16(lo)
	e.a = e.getMemory(adr + uint16(e.x))
	e.zf = e.a == 0
	e.nf = (e.a & 0x80) > 0
	return fmt.Sprintf("%.2x %.2x   lda $%.4x,X", lo, hi, adr)
}

func lda_abs_y(e *emu6502) string {
	lo := e.getMnemonic()
	hi := e.getMnemonic()
	adr := uint16(hi)*256 + uint16(lo)
	e.a = e.getMemory(adr + uint16(e.y))
	e.zf = e.a == 0
	e.nf = (e.a & 0x80) > 0
	return fmt.Sprintf("%.2x %.2x   lda $%.4x,Y", lo, hi, adr)
}

func lda_zp(e *emu6502) string {
	lo := e.getMnemonic()
	adr := uint16(lo)
	e.a = e.getMemory(adr)
	e.zf = e.a == 0
	e.nf = (e.a & 0x80) > 0
	return fmt.Sprintf("%.2x        lda $%.2x", lo, lo)
}

func lda_zp_x(e *emu6502) string {
	lo := e.getMnemonic()
	adr := uint16(lo)
	e.a = e.getMemory(adr + uint16(e.x))
	e.zf = e.a == 0
	e.nf = (e.a & 0x80) > 0
	return fmt.Sprintf("%.2x        lda $%.2x,X", lo, adr)
}

func lda_ind_x(e *emu6502) string {
	zp := e.getMnemonic()
	zpx := zp + e.x
	lo := e.getMemory(uint16(zpx))
	hi := e.getMemory(uint16(zpx + 1))
	adr := uint16(hi)*256 + uint16(lo)
	e.a = e.getMemory(adr)
	e.zf = e.a == 0
	e.nf = (e.a & 0x80) > 0
	return fmt.Sprintf("%.2x        lda ($%.2x,X)", lo, adr)
}

func lda_ind_y(e *emu6502) string {
	zp := e.getMnemonic()
	lo := e.getMemory(uint16(zp))
	hi := e.getMemory(uint16(zp + 1))
	adr := uint16(hi)*256 + uint16(lo)
	e.a = e.getMemory(adr + uint16(e.y))
	e.zf = e.a == 0
	e.nf = (e.a & 0x80) > 0
	return fmt.Sprintf("%.2x        lda ($%.2x),Y", lo, adr)
}

func ldx_direct(e *emu6502) string {
	e.x = e.getMnemonic()
	e.zf = e.x == 0
	e.nf = (e.x & 0x80) > 0
	return fmt.Sprintf("%.2x        ldx #$%.2x", e.x, e.x)
}

func ldx_abs(e *emu6502) string {
	lo := e.getMnemonic()
	hi := e.getMnemonic()
	adr := uint16(hi)*256 + uint16(lo)
	e.x = e.getMemory(adr)
	e.zf = e.x == 0
	e.nf = (e.x & 0x80) > 0
	return fmt.Sprintf("%.2x %.2x   ldx $%.4x", lo, hi, adr)
}

func ldx_abs_y(e *emu6502) string {
	lo := e.getMnemonic()
	hi := e.getMnemonic()
	adr := uint16(hi)*256 + uint16(lo)
	e.x = e.getMemory(adr + uint16(e.y))
	e.zf = e.x == 0
	e.nf = (e.x & 0x80) > 0
	return fmt.Sprintf("%.2x %.2x   ldx $%.4x,Y", lo, hi, adr)
}

func ldx_zp(e *emu6502) string {
	lo := e.getMnemonic()
	adr := uint16(lo)
	e.x = e.getMemory(adr)
	e.zf = e.x == 0
	e.nf = (e.x & 0x80) > 0
	return fmt.Sprintf("%.2x        ldx $%.2x", lo, lo)
}

func ldx_zp_y(e *emu6502) string {
	lo := e.getMnemonic()
	adr := uint16(lo)
	e.x = e.getMemory(adr + uint16(e.y))
	e.zf = e.x == 0
	e.nf = (e.x & 0x80) > 0
	return fmt.Sprintf("%.2x        ldx $%.2x,Y", lo, adr)
}

func ldy_direct(e *emu6502) string {
	e.y = e.getMnemonic()
	e.zf = e.y == 0
	e.nf = (e.y & 0x80) > 0
	return fmt.Sprintf("%.2x        ldy #$%.2x", e.y, e.y)
}

func ldy_abs(e *emu6502) string {
	lo := e.getMnemonic()
	hi := e.getMnemonic()
	adr := uint16(hi)*256 + uint16(lo)
	e.y = e.getMemory(adr)
	e.zf = e.y == 0
	e.nf = (e.y & 0x80) > 0
	return fmt.Sprintf("%.2x %.2x   ldy $%.4x", lo, hi, adr)
}

func ldy_abs_x(e *emu6502) string {
	lo := e.getMnemonic()
	hi := e.getMnemonic()
	adr := uint16(hi)*256 + uint16(lo)
	e.y = e.getMemory(adr + uint16(e.x))
	e.zf = e.y == 0
	e.nf = (e.y & 0x80) > 0
	return fmt.Sprintf("%.2x %.2x   ldy $%.4x,X", lo, hi, adr)
}

func ldy_zp(e *emu6502) string {
	lo := e.getMnemonic()
	adr := uint16(lo)
	e.y = e.getMemory(adr)
	e.zf = e.y == 0
	e.nf = (e.y & 0x80) > 0
	return fmt.Sprintf("%.2x        ldy $%.2x", lo, lo)
}

func ldy_zp_x(e *emu6502) string {
	lo := e.getMnemonic()
	adr := uint16(lo)
	e.y = e.getMemory(adr + uint16(e.y))
	e.zf = e.y == 0
	e.nf = (e.y & 0x80) > 0
	return fmt.Sprintf("%.2x        ldy $%.2x,X", lo, adr)
}
