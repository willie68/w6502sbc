package emulator

import "fmt"

func and_direct(e *emu6502) string {
	v := e.getMnemonic()
	e.a = e.a & v
	e.setFlags(e.a, nil, nil)
	return fmt.Sprintf("           and #$%.2x", v)
}

func and_abs(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	e.a = e.a & v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("   and $%.4x", adr)
}

func and_abs_x(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	e.a = e.a & v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("   and $%.4x,X", adr)
}

func and_abs_y(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.y))
	e.a = e.a & v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("   and $%.4x,Y", adr)
}

func and_zp(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	e.a = e.a & v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("   and $%.2x", adr)
}

func and_zp_x(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	e.a = e.a & v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("    and $%.2x,X", adr)
}

func and_ind_x(e *emu6502) string {
	zp, str := e.getZPAddress()
	zpx := zp + uint16(e.x)
	adr := e.readVector(zpx)
	v := e.getMemory(adr)
	e.a = e.a & v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("        and ($%.2x,X)", adr)
}

func and_ind_y(e *emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	v := e.getMemory(adr + uint16(e.y))
	e.a = e.a & v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("        and ($%.2x),Y", adr)
}

func ora_direct(e *emu6502) string {
	v := e.getMnemonic()
	e.a = e.a | v
	e.setFlags(e.a, nil, nil)
	return fmt.Sprintf("           ora #$%.2x", v)
}

func ora_abs(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	e.a = e.a | v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("   ora $%.4x", adr)
}

func ora_abs_x(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	e.a = e.a | v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("   ora $%.4x,X", adr)
}

func ora_abs_y(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.y))
	e.a = e.a | v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("   ora $%.4x,Y", adr)
}

func ora_zp(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	e.a = e.a | v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("   ora $%.2x", adr)
}

func ora_zp_x(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	e.a = e.a | v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("    ora $%.2x,X", adr)
}

func ora_ind_x(e *emu6502) string {
	zp, str := e.getZPAddress()
	zpx := zp + uint16(e.x)
	adr := e.readVector(zpx)
	v := e.getMemory(adr)
	e.a = e.a | v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("        ora ($%.2x,X)", adr)
}

func ora_ind_y(e *emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	v := e.getMemory(adr + uint16(e.y))
	e.a = e.a | v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("        ora ($%.2x),Y", adr)
}

func eor_direct(e *emu6502) string {
	v := e.getMnemonic()
	e.a = e.a ^ v
	e.setFlags(e.a, nil, nil)
	return fmt.Sprintf("           eor #$%.2x", v)
}

func eor_abs(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	e.a = e.a ^ v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("   eor $%.4x", adr)
}

func eor_abs_x(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	e.a = e.a ^ v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("   eor $%.4x,X", adr)
}

func eor_abs_y(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.y))
	e.a = e.a ^ v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("   eor $%.4x,Y", adr)
}

func eor_zp(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	e.a = e.a ^ v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("   eor $%.2x", adr)
}

func eor_zp_x(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	e.a = e.a ^ v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("    eor $%.2x,X", adr)
}

func eor_ind_x(e *emu6502) string {
	zp, str := e.getZPAddress()
	zpx := zp + uint16(e.x)
	adr := e.readVector(zpx)
	v := e.getMemory(adr)
	e.a = e.a ^ v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("        eor ($%.2x,X)", adr)
}

func eor_ind_y(e *emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	v := e.getMemory(adr + uint16(e.y))
	e.a = e.a ^ v
	e.setFlags(e.a, nil, nil)
	return str + fmt.Sprintf("        eor ($%.2x),Y", adr)
}
