package emulator

import "fmt"

var functions = []func(*emu6502) string{
	0x0a: asl_a,
	0x2a: rol_a,
	0x4c: jmp_abs,
	0x4a: lsr_a,
	0x6a: ror_a,
	0x8a: txa,
	0x8d: sta_abs,
	0x98: tya,
	0xa1: lda_ind_x,
	0xa5: lda_zp,
	0xa8: tay,
	0xa9: lda_direct,
	0xaa: tax,
	0xad: lda_abs,
	0xb1: lda_ind_y,
	0xb5: lda_zp_x,
	0xb9: lda_abs_y,
	0xbd: lda_abs_x,
	0xea: nop,
	0xff: nil,
}

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

func sta_abs(e *emu6502) string {
	lo := e.getMnemonic()
	hi := e.getMnemonic()
	adr := uint16(hi)*256 + uint16(lo)
	e.setMemory(adr, e.a)
	return fmt.Sprintf("%.2x %.2x     sta $%.4x", lo, hi, adr)
}

func tax(e *emu6502) string {
	e.x = e.a
	e.zf = e.a == 0
	e.nf = (e.a & 0x80) > 0
	return "           tax"
}

func tay(e *emu6502) string {
	e.y = e.a
	e.zf = e.a == 0
	e.nf = (e.a & 0x80) > 0
	return "           tay"
}

func txa(e *emu6502) string {
	e.a = e.x
	e.zf = e.a == 0
	e.nf = (e.a & 0x80) > 0
	return "           txa"
}

func tya(e *emu6502) string {
	e.a = e.y
	e.zf = e.a == 0
	e.nf = (e.a & 0x80) > 0
	return "           tya"
}

func ror_a(e *emu6502) string {
	tmp := e.cf
	e.cf = (e.a & 0x01) > 0
	e.a = e.a >> 1
	if tmp {
		e.a = e.a | 0x80
	}
	e.zf = e.a == 0
	e.nf = (e.a & 0x80) > 0
	return "           ror"
}

func lsr_a(e *emu6502) string {
	e.cf = (e.a & 0x01) > 0
	e.a = e.a >> 1
	e.zf = e.a == 0
	e.nf = (e.a & 0x80) > 0
	return "           lsr"
}

func rol_a(e *emu6502) string {
	tmp := e.cf
	e.cf = (e.a & 0x80) > 0
	e.a = e.a << 1
	if tmp {
		e.a = e.a | 0x01
	}
	e.zf = e.a == 0
	e.nf = (e.a & 0x80) > 0
	return "           rol"
}

func asl_a(e *emu6502) string {
	e.cf = (e.a & 0x80) > 0
	e.a = e.a << 1
	e.zf = e.a == 0
	e.nf = (e.a & 0x80) > 0
	return "           asl"
}

func jmp_abs(e *emu6502) string {
	lo := e.getMnemonic()
	hi := e.getMnemonic()
	adr := uint16(hi)*256 + uint16(lo)
	e.address = adr
	return fmt.Sprintf("%.2x %.2x     jmp $%.4x", lo, hi, adr)
}

func nop(e *emu6502) string {
	return "           nop"
}
