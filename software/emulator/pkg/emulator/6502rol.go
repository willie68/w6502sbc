package emulator

import "fmt"

func i_asl(e *emu6502, v uint8) uint8 {
	e.cf = (v & 0x80) > 0
	v = v << 1
	return v
}

func asl(e *emu6502) string {
	e.a = i_asl(e, e.a)
	e.setFlags(e.a, nil, nil)
	return "          asl"
}

func asl_abs(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	v = i_asl(e, v)
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   als $%.4x", adr)
}

func asl_abs_x(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	v = i_asl(e, v)
	e.setMemory(adr+uint16(e.x), v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   als $%.4x,X", adr)
}

func asl_zp(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	v = i_asl(e, v)
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   als $%.2x", adr)
}

func asl_zp_x(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	v = i_asl(e, v)
	e.setMemory(adr+uint16(e.x), v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   als $%.2x,X", adr)
}

func i_lsr(e *emu6502, v uint8) uint8 {
	e.cf = (v & 0x01) > 0
	v = v >> 1
	return v
}

func lsr(e *emu6502) string {
	e.a = i_lsr(e, e.a)
	e.setFlags(e.a, nil, nil)
	return "          lsr"
}

func lsr_abs(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	v = i_lsr(e, v)
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   lsr $%.4x", adr)
}

func lsr_abs_x(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	v = i_lsr(e, v)
	e.setMemory(adr+uint16(e.x), v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   lsr $%.4x,X", adr)
}

func lsr_zp(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	v = i_lsr(e, v)
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   lsr $%.2x", adr)
}

func lsr_zp_x(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	v = i_lsr(e, v)
	e.setMemory(adr+uint16(e.x), v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   lsr $%.2x,X", adr)
}

func i_rol(e *emu6502, v uint8) uint8 {
	tmp := e.cf
	e.cf = (v & 0x80) > 0
	v = v << 1
	if tmp {
		v = v | 0x01
	}
	return v
}

func rol(e *emu6502) string {
	e.a = i_rol(e, e.a)
	e.setFlags(e.a, nil, nil)
	return "          rol"
}

func rol_abs(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	v = i_rol(e, v)
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   rol $%.4x", adr)
}

func rol_abs_x(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	v = i_rol(e, v)
	e.setMemory(adr+uint16(e.x), v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   rol $%.4x,X", adr)
}

func rol_zp(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	v = i_rol(e, v)
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   rol $%.2x", adr)
}

func rol_zp_x(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	v = i_rol(e, v)
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   rol $%.2x,X", adr)
}

func i_ror(e *emu6502, v uint8) uint8 {
	tmp := e.cf
	e.cf = (v & 0x01) > 0
	v = v >> 1
	if tmp {
		e.a = e.a | 0x80
	}
	return v
}

func ror(e *emu6502) string {
	e.a = i_ror(e, e.a)
	e.setFlags(e.a, nil, nil)
	return "          ror"
}

func ror_abs(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	v = i_ror(e, v)
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   ror $%.4x", adr)
}

func ror_abs_x(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	v = i_ror(e, v)
	e.setMemory(adr+uint16(e.x), v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   ror $%.4x,X", adr)
}

func ror_zp(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	v = i_ror(e, v)
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   ror $%.2x", adr)
}

func ror_zp_x(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	v = i_ror(e, v)
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   ror $%.2x,X", adr)
}
