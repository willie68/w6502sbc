package emulator

import "fmt"

func i_cmp(e *Emu6502, v1 uint8, v2 uint8) {
	e.cf = v1 >= v2
	e.zf = v1 == v2
	e.nf = ((v1 - v2) & 0x80) != 0
}

func cmp_direct(e *Emu6502) string {
	v := e.getMnemonic()
	i_cmp(e, e.a, v)
	return fmt.Sprintf("%.2x       cmp #$%.2x", v, v)
}

func cmp_abs(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	i_cmp(e, e.a, v)
	return str + fmt.Sprintf("    cmp $%.4x", adr)
}

func cmp_abs_x(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	i_cmp(e, e.a, v)
	return str + fmt.Sprintf("    cmp $%.4x,X", adr)
}

func cmp_abs_y(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.y))
	i_cmp(e, e.a, v)
	return str + fmt.Sprintf("    cmp $%.4x,Y", adr)
}

func cmp_zp(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	i_cmp(e, e.a, v)
	return str + fmt.Sprintf("    cmp $%.2x", adr)
}

func cmp_zp_x(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	i_cmp(e, e.a, v)
	return str + fmt.Sprintf("    cmp $%.2x,X", adr)
}

func cmp_ind_x(e *Emu6502) string {
	zp, str := e.getZPAddress()
	zpx := zp + uint16(e.x)
	adr := e.readVector(zpx)
	v := e.getMemory(adr)
	i_cmp(e, e.a, v)
	return str + fmt.Sprintf("    cmp ($%.2x,X)", adr)
}

func cmp_ind_y(e *Emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	v := e.getMemory(adr + uint16(e.y))
	i_cmp(e, e.a, v)
	return str + fmt.Sprintf("    cmp ($%.2x),Y", adr)
}

func cpx_direct(e *Emu6502) string {
	v := e.getMnemonic()
	i_cmp(e, e.x, v)
	return fmt.Sprintf("%.2x       cpx #$%.2x", v, v)
}

func cpx_abs(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	i_cmp(e, e.x, v)
	return str + fmt.Sprintf("    cpx $%.4x", adr)
}

func cpx_zp(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	i_cmp(e, e.x, v)
	return str + fmt.Sprintf("    cpx $%.2x", adr)
}

func cpy_direct(e *Emu6502) string {
	v := e.getMnemonic()
	i_cmp(e, e.y, v)
	return fmt.Sprintf("%.2x       cpy #$%.2x", v, v)
}

func cpy_abs(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	i_cmp(e, e.y, v)
	return str + fmt.Sprintf("    cpy $%.4x", adr)
}

func cpy_zp(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	i_cmp(e, e.y, v)
	return str + fmt.Sprintf("    cpy $%.2x", adr)
}

func i_bit(e *Emu6502, v uint8) {
	e.nf = (v & 0x80) > 0
	e.vf = (v & 0x40) > 0
	e.zf = (v & e.a) == 0
}

func bit_abs(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	i_bit(e, v)
	return str + fmt.Sprintf("   bit $%.4x", adr)
}

func bit_zp(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	i_bit(e, v)
	return str + fmt.Sprintf("   bit $%.2x", adr)
}
