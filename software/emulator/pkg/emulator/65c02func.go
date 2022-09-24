package emulator

import "fmt"

var cfunc = []func(*emu6502) string{
	0x04: tsb_zp, 0x0c: tsb_abs, 0x0f: bbr0,
	0x12: ora_zp_ind, 0x14: trb_zp, 0x1a: inc, 0x1c: trb_abs, 0x1f: bbr1,
	0x2f: bbr2,
	0x32: and_zp_ind, 0x34: bit_zp_x, 0x3a: dec, 0x3c: bit_abs_x, 0x3f: bbr3,
	0x4f: bbr4,
	0x52: eor_zp_ind, 0x5a: phy, 0x5f: bbr5,
	0x64: stz_zp, 0x6f: bbr6,
	0x72: adc_zp_ind, 0x74: stz_zp_x, 0x7a: ply, 0x7c: jmp_abs_x, 0x7f: bbr7,
	0x80: bra, 0x89: bit_direct, 0x8f: bbs0,
	0x92: sta_zp_ind, 0x9c: stz_abs, 0x9e: stz_abs_x, 0x9f: bbs1,
	0xaf: bbs2,
	0xb2: lda_zp_ind, 0xbf: bbs3,
	0xcf: bbs4,
	0xd2: cmp_zp_ind, 0xda: phx, 0xdf: bbs5,
	0xef: bbs6,
	0xf2: sbc_zp_ind, 0xfa: plx, 0xff: bbs7,
}

func rmb0(e *emu6502) string {
	zp, str := e.getZPAddress()
	v := e.getMemory(zp)
	v = v & 0xfe
	e.setMemory(zp, v)
	return str + fmt.Sprintf("      rmb0 $%.2x", zp)
}

func smb0(e *emu6502) string {
	zp, str := e.getZPAddress()
	v := e.getMemory(zp)
	v = v | 0x01
	e.setMemory(zp, v)
	return str + fmt.Sprintf("      smb0 $%.2x", zp)
}

func i_bbx(e *emu6502, bit uint8, c bool) (uint16, string) {
	zp, str := e.getZPAddress()
	rel := e.getMnemonic()
	v := e.getMemory(zp)
	if c {
		if (v & bit) == 0 {
			badr(e, rel)
		}
	} else {
		if (v & bit) > 0 {
			badr(e, rel)
		}
	}
	return zp, str
}

func bbr0(e *emu6502) string {
	zp, str := i_bbx(e, 0x01, true)
	return str + fmt.Sprintf("   bbr0 $%.2x, $%.4x", zp, e.address)
}

func bbs0(e *emu6502) string {
	zp, str := i_bbx(e, 0x01, false)
	return str + fmt.Sprintf("   bbs0 $%.2x, $%.4x", zp, e.address)
}

func bbr1(e *emu6502) string {
	zp, str := i_bbx(e, 0x02, true)
	return str + fmt.Sprintf("   bbr1 $%.2x, $%.4x", zp, e.address)
}

func bbs1(e *emu6502) string {
	zp, str := i_bbx(e, 0x02, false)
	return str + fmt.Sprintf("   bbs1 $%.2x, $%.4x", zp, e.address)
}

func bbr2(e *emu6502) string {
	zp, str := i_bbx(e, 0x04, true)
	return str + fmt.Sprintf("   bbr2 $%.2x, $%.4x", zp, e.address)
}

func bbs2(e *emu6502) string {
	zp, str := i_bbx(e, 0x04, false)
	return str + fmt.Sprintf("   bbs2 $%.2x, $%.4x", zp, e.address)
}

func bbr3(e *emu6502) string {
	zp, str := i_bbx(e, 0x08, true)
	return str + fmt.Sprintf("   bbr3 $%.2x, $%.4x", zp, e.address)
}

func bbs3(e *emu6502) string {
	zp, str := i_bbx(e, 0x08, false)
	return str + fmt.Sprintf("   bbs3 $%.2x, $%.4x", zp, e.address)
}

func bbr4(e *emu6502) string {
	zp, str := i_bbx(e, 0x10, true)
	return str + fmt.Sprintf("   bbr4 $%.2x, $%.4x", zp, e.address)
}

func bbs4(e *emu6502) string {
	zp, str := i_bbx(e, 0x10, false)
	return str + fmt.Sprintf("   bbs4 $%.2x, $%.4x", zp, e.address)
}

func bbr5(e *emu6502) string {
	zp, str := i_bbx(e, 0x20, true)
	return str + fmt.Sprintf("   bbr5 $%.2x, $%.4x", zp, e.address)
}

func bbs5(e *emu6502) string {
	zp, str := i_bbx(e, 0x20, false)
	return str + fmt.Sprintf("   bbs5 $%.2x, $%.4x", zp, e.address)
}

func bbr6(e *emu6502) string {
	zp, str := i_bbx(e, 0x40, true)
	return str + fmt.Sprintf("   bbr6 $%.2x, $%.4x", zp, e.address)
}

func bbs6(e *emu6502) string {
	zp, str := i_bbx(e, 0x40, false)
	return str + fmt.Sprintf("   bbs6 $%.2x, $%.4x", zp, e.address)
}

func bbr7(e *emu6502) string {
	zp, str := i_bbx(e, 0x80, true)
	return str + fmt.Sprintf("   bbr7 $%.2x, $%.4x", zp, e.address)
}

func bbs7(e *emu6502) string {
	zp, str := i_bbx(e, 0x80, false)
	return str + fmt.Sprintf("   bbs7 $%.2x, $%.4x", zp, e.address)
}

func i_tsb(e *emu6502, v uint8) uint8 {
	i_bit(e, v)
	return (e.a | v)
}

func tsb_abs(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	e.setMemory(adr, i_tsb(e, v))
	return str + fmt.Sprintf("     tsb $%.4x", adr)
}

func tsb_zp(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	e.setMemory(adr, i_tsb(e, v))
	return str + fmt.Sprintf("     tsb $%.2x", adr)
}

func i_trb(e *emu6502, v uint8) uint8 {
	i_bit(e, v)
	return (e.a ^ 0xff) & v
}

func trb_abs(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	e.setMemory(adr, i_trb(e, v))
	return str + fmt.Sprintf("     trb $%.4x", adr)
}

func trb_zp(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	e.setMemory(adr, i_trb(e, v))
	return str + fmt.Sprintf("     trb $%.2x", adr)
}

func stz_abs(e *emu6502) string {
	adr, str := e.getAddress()
	e.setMemory(adr, 0)
	return str + fmt.Sprintf("     stz $%.4x", adr)
}

func stz_abs_x(e *emu6502) string {
	adr, str := e.getAddress()
	e.setMemory(adr+uint16(e.x), 0)
	return str + fmt.Sprintf("     stz $%.4x,X", adr)
}

func stz_zp(e *emu6502) string {
	adr, str := e.getZPAddress()
	e.setMemory(adr, 0)
	return str + fmt.Sprintf("     stz $%.2x", adr)
}

func stz_zp_x(e *emu6502) string {
	adr, str := e.getZPAddress()
	e.setMemory(adr+uint16(e.x), 0)
	return str + fmt.Sprintf("     stz $%.2x,X", adr)
}

func plx(e *emu6502) string {
	e.x = e.pop()
	return "          plx"
}

func phx(e *emu6502) string {
	e.push(e.x)
	return "          phx"
}

func ply(e *emu6502) string {
	e.a = e.pop()
	return "          ply"
}

func phy(e *emu6502) string {
	e.push(e.y)
	return "          phy"
}

func bra(e *emu6502) string {
	v := e.getMnemonic()
	badr(e, v)
	return fmt.Sprintf("%.2x        bra $%.4x", v, e.address)
}

func inc(e *emu6502) string {
	e.a++
	e.setFlags(e.a, nil, nil)
	return "          inc"
}

func dec(e *emu6502) string {
	e.a--
	e.setFlags(e.a, nil, nil)
	return "          dec"
}

func bit_direct(e *emu6502) string {
	v := e.getMnemonic()
	i_bit(e, v)
	return fmt.Sprintf("%.2x        bit #$%.2x", v, v)
}

func bit_zp_x(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	i_bit(e, v)
	return str + fmt.Sprintf("   bit $%.2x,X", adr)
}

func bit_abs_x(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	i_bit(e, v)
	return str + fmt.Sprintf("   bit $%.4x,X", adr)
}

func ora_zp_ind(e *emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	v := e.getMemory(adr)
	e.a = e.a | v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("        ora ($%.2x)", adr)
}

func and_zp_ind(e *emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	v := e.getMemory(adr)
	e.a = e.a & v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("        and ($%.2x)", adr)
}

func eor_zp_ind(e *emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	v := e.getMemory(adr)
	e.a = e.a ^ v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("        eor ($%.2x)", adr)
}

func adc_zp_ind(e *emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	v := e.getMemory(adr)
	adc(e, v)
	return str + fmt.Sprintf("        adc ($%.2x)", adr)
}

func sta_zp_ind(e *emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	e.setMemory(adr, e.a)
	return str + fmt.Sprintf("     sta ($%.2x)", adr)
}

func lda_zp_ind(e *emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	v := e.getMemory(adr)
	e.setFlags(v, nil, nil)
	e.a = v
	return str + fmt.Sprintf("        lda ($%.2x)", adr)
}

func cmp_zp_ind(e *emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	v := e.getMemory(adr)
	i_cmp(e, e.a, v)
	return str + fmt.Sprintf("        cmp ($%.2x)", adr)
}

func sbc_zp_ind(e *emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	v := e.getMemory(adr)
	sbc(e, v)
	return str + fmt.Sprintf("     sbc ($%.2x)", adr)
}

func jmp_abs_x(e *emu6502) string {
	adr, str := e.getAddress()
	e.address = e.readVector(adr + uint16(e.x))
	return str + fmt.Sprintf("     jmp ($%.4x,X)", adr)
}
