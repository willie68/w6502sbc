package emulator

import "fmt"

var functions = []func(*emu6502) string{
	0x01: ora_ind_x, 0x05: ora_zp,
	0x09: ora_direct, 0x0a: asl_a, 0x0d: ora_abs,
	0x11: ora_ind_y, 0x15: ora_zp_x,
	0x19: ora_abs_y, 0x1d: ora_abs_x,
	0x21: and_ind_x, 0x25: and_zp,
	0x29: and_direct, 0x2a: rol_a, 0x2d: and_abs,
	0x31: and_ind_y, 0x35: and_zp_x,
	0x39: and_abs_y, 0x3d: and_abs_x,
	0x41: eor_ind_x, 0x45: eor_zp,
	0x49: eor_direct, 0x4c: jmp_abs, 0x4a: lsr_a, 0x4d: eor_abs,
	0x51: eor_ind_y, 0x55: eor_zp_x,
	0x59: eor_abs_y, 0x5d: eor_abs_x,
	0x61: adc_ind_x, 0x65: adc_zp,
	0x69: adc_direct, 0x6a: ror_a, 0x6d: adc_abs,
	0x71: adc_ind_y, 0x75: adc_abs_x,
	0x79: adc_abs_y, 0x7d: adc_abs_x,
	0x81: sta_ind_x, 0x84: sty_zp, 0x85: sta_zp, 0x86: stx_zp,
	0x8a: txa, 0x8c: sty_abs, 0x8d: sta_abs, 0x8e: stx_abs,
	0x91: sta_ind_y, 0x94: sty_zp_x, 0x95: sta_zp_x, 0x96: stx_zp_y,
	0x98: tya, 0x99: sta_abs_y, 0x9a: txs, 0x9d: sta_abs_x,
	0xa0: ldy_direct, 0xa1: lda_ind_x, 0xa2: ldx_direct, 0xa4: ldy_zp, 0xa5: lda_zp, 0xa6: ldx_zp,
	0xa8: tay, 0xa9: lda_direct, 0xaa: tax, 0xac: ldy_abs, 0xad: lda_abs, 0xae: ldx_abs,
	0xb1: lda_ind_y, 0xb4: ldy_zp_x, 0xb5: lda_zp_x, 0xb6: ldx_zp_y,
	0xb9: lda_abs_y, 0xba: tsx, 0xbc: ldy_abs_x, 0xbd: lda_abs_x, 0xbe: ldx_abs_y,
	0xe1: sbc_ind_x, 0xe5: sbc_zp,
	0xe9: sbc_direct, 0xea: nop, 0xed: sbc_abs,
	0xf1: sbc_ind_y, 0xf5: sbc_zp_x,
	0xf9: sbc_abs_y, 0xfd: sbc_abs_x, 0xff: nil,
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
