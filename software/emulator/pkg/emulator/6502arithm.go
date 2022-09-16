package emulator

import "fmt"

func adc(e *emu6502, v uint8) {
	e.of = (((uint16(e.a) ^ uint16(v)) & 0x80) == 0)
	t := uint16(e.a) + uint16(v)
	if e.cf {
		t++
	}
	e.a = uint8(t)
	cf := t > 0x00ff
	of := e.of
	if t > 0x00ff {
		if t > 0x017f {
			of = false
		}
	} else {
		if t < 0x0080 {
			of = false
		}
	}
	e.setFlags(e.a, &cf, &of)
}

func adc_direct(e *emu6502) string {
	v := e.getMnemonic()
	adc(e, v)
	return "           adc"
}

func adc_abs(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	adc(e, v)
	return str + fmt.Sprintf("   adc $%.4x", adr)
}

func adc_abs_x(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	adc(e, v)
	return str + fmt.Sprintf("   adc $%.4x,X", adr)
}

func adc_abs_y(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.y))
	adc(e, v)
	return str + fmt.Sprintf("   adc $%.4x,Y", adr)
}

func adc_zp(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	adc(e, v)
	return str + fmt.Sprintf("   adc $%.2x", adr)
}

func adc_zp_x(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	adc(e, v)
	return str + fmt.Sprintf("    adc $%.2x,X", adr)
}

func adc_ind_x(e *emu6502) string {
	zp, str := e.getZPAddress()
	zpx := zp + uint16(e.x)
	adr := e.readVector(zpx)
	v := e.getMemory(adr)
	adc(e, v)
	return str + fmt.Sprintf("        adc ($%.2x,X)", adr)
}

func adc_ind_y(e *emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	v := e.getMemory(adr + uint16(e.y))
	adc(e, v)
	return str + fmt.Sprintf("        adc ($%.2x),Y", adr)
}

func sbc(e *emu6502, v uint8) {
	e.of = (((uint16(e.a) ^ uint16(v)) & 0x80) != 0)
	t := uint16(0xff) + uint16(e.a) - uint16(v)
	if e.cf {
		t++
	}
	e.a = uint8(t)
	cf := t > 0x00ff
	of := e.of
	if t > 0x00ff {
		if t > 0x017f {
			of = false
		}
	} else {
		if t < 0x0080 {
			of = false
		}
	}
	e.setFlags(e.a, &cf, &of)
}

func sbc_direct(e *emu6502) string {
	v := e.getMnemonic()
	sbc(e, v)
	return "           sbc"
}

func sbc_abs(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	sbc(e, v)
	return str + fmt.Sprintf("   sbc $%.4x", adr)
}

func sbc_abs_x(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	sbc(e, v)
	return str + fmt.Sprintf("   sbc $%.4x,X", adr)
}

func sbc_abs_y(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.y))
	sbc(e, v)
	return str + fmt.Sprintf("   sbc $%.4x,Y", adr)
}

func sbc_zp(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	sbc(e, v)
	return str + fmt.Sprintf("   sbc $%.2x", adr)
}

func sbc_zp_x(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	sbc(e, v)
	return str + fmt.Sprintf("    sbc $%.2x,X", adr)
}

func sbc_ind_x(e *emu6502) string {
	zp, str := e.getZPAddress()
	zpx := zp + uint16(e.x)
	adr := e.readVector(zpx)
	v := e.getMemory(adr)
	sbc(e, v)
	return str + fmt.Sprintf("        sbc ($%.2x,X)", adr)
}

func sbc_ind_y(e *emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	v := e.getMemory(adr + uint16(e.y))
	sbc(e, v)
	return str + fmt.Sprintf("        sbc ($%.2x),Y", adr)
}
