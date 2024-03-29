package emulator

import "fmt"

func adc(e *Emu6502, v uint8) {
	e.vf = (((uint16(e.a) ^ uint16(v)) & 0x80) == 0)
	t := uint16(e.a) + uint16(v)
	if e.cf {
		t++
	}
	e.a = uint8(t)
	cf := t > 0x00ff
	vf := e.vf
	if t > 0x00ff {
		if t > 0x017f {
			vf = false
		}
	} else {
		if t < 0x0080 {
			vf = false
		}
	}
	e.setFlags(e.a, &cf, &vf)
}

func adc_direct(e *Emu6502) string {
	v := e.getMnemonic()
	adc(e, v)
	return fmt.Sprintf("%.2x        adc #$%.2x", v, v)
}

func adc_abs(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	adc(e, v)
	return str + fmt.Sprintf("   adc $%.4x", adr)
}

func adc_abs_x(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	adc(e, v)
	return str + fmt.Sprintf("   adc $%.4x,X", adr)
}

func adc_abs_y(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.y))
	adc(e, v)
	return str + fmt.Sprintf("   adc $%.4x,Y", adr)
}

func adc_zp(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	adc(e, v)
	return str + fmt.Sprintf("   adc $%.2x", adr)
}

func adc_zp_x(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	adc(e, v)
	return str + fmt.Sprintf("    adc $%.2x,X", adr)
}

func adc_ind_x(e *Emu6502) string {
	zp, str := e.getZPAddress()
	zpx := zp + uint16(e.x)
	adr := e.readVector(zpx)
	v := e.getMemory(adr)
	adc(e, v)
	return str + fmt.Sprintf("        adc ($%.2x,X)", adr)
}

func adc_ind_y(e *Emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	v := e.getMemory(adr + uint16(e.y))
	adc(e, v)
	return str + fmt.Sprintf("        adc ($%.2x),Y", adr)
}

func sbc(e *Emu6502, v uint8) {
	e.vf = (((uint16(e.a) ^ uint16(v)) & 0x80) != 0)
	t := uint16(0xff) + uint16(e.a) - uint16(v)
	if e.cf {
		t++
	}
	e.a = uint8(t)
	cf := t > 0x00ff
	vf := e.vf
	if t > 0x00ff {
		if t > 0x017f {
			vf = false
		}
	} else {
		if t < 0x0080 {
			vf = false
		}
	}
	e.setFlags(e.a, &cf, &vf)
}

func sbc_direct(e *Emu6502) string {
	v := e.getMnemonic()
	sbc(e, v)
	return fmt.Sprintf("%.2x        sbc #$%.2x", v, v)
}

func sbc_abs(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	sbc(e, v)
	return str + fmt.Sprintf("   sbc $%.4x", adr)
}

func sbc_abs_x(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	sbc(e, v)
	return str + fmt.Sprintf("   sbc $%.4x,X", adr)
}

func sbc_abs_y(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.y))
	sbc(e, v)
	return str + fmt.Sprintf("   sbc $%.4x,Y", adr)
}

func sbc_zp(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	sbc(e, v)
	return str + fmt.Sprintf("   sbc $%.2x", adr)
}

func sbc_zp_x(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	sbc(e, v)
	return str + fmt.Sprintf("    sbc $%.2x,X", adr)
}

func sbc_ind_x(e *Emu6502) string {
	zp, str := e.getZPAddress()
	zpx := zp + uint16(e.x)
	adr := e.readVector(zpx)
	v := e.getMemory(adr)
	sbc(e, v)
	return str + fmt.Sprintf("        sbc ($%.2x,X)", adr)
}

func sbc_ind_y(e *Emu6502) string {
	zp, str := e.getZPAddress()
	adr := e.readVector(zp)
	v := e.getMemory(adr + uint16(e.y))
	sbc(e, v)
	return str + fmt.Sprintf("        sbc ($%.2x),Y", adr)
}

func inc_abs(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	v++
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   inc $%.4x", adr)
}

func inc_abs_x(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	v++
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   inc $%.4x,X", adr)
}

func inc_zp(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	v++
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   inc $%.2x", adr)
}

func inc_zp_x(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	v++
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("    inc $%.2x,X", adr)
}

func dec_abs(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	v--
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   dec $%.4x", adr)
}

func dec_abs_x(e *Emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	v--
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   dec $%.4x,X", adr)
}

func dec_zp(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	v--
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   dec $%.2x", adr)
}

func dec_zp_x(e *Emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	v--
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("    dec $%.2x,X", adr)
}

func inx(e *Emu6502) string {
	e.x++
	e.setFlags(e.x, nil, nil)
	return "           inx"
}

func iny(e *Emu6502) string {
	e.y++
	e.setFlags(e.y, nil, nil)
	return "           iny"
}

func dex(e *Emu6502) string {
	e.x--
	e.setFlags(e.x, nil, nil)
	return "           dex"
}

func dey(e *Emu6502) string {
	e.y--
	e.setFlags(e.y, nil, nil)
	return "           dey"
}
