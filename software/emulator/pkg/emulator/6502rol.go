package emulator

import "fmt"

func asl(e *emu6502) string {
	e.cf = (e.a & 0x80) > 0
	e.a = e.a << 1
	e.setFlags(e.a, nil, nil)
	return "          asl"
}

func asl_abs(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	e.cf = (v & 0x80) > 0
	v = v << 1
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   als $%.4x", adr)
}

func asl_abs_x(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	e.cf = (v & 0x80) > 0
	v = v << 1
	e.setMemory(adr+uint16(e.x), v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   als $%.4x,X", adr)
}

func asl_zp(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	e.cf = (v & 0x80) > 0
	v = v << 1
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   als $%.2x", adr)
}

func asl_zp_x(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	e.cf = (v & 0x80) > 0
	v = v << 1
	e.setMemory(adr+uint16(e.x), v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   als $%.2x,X", adr)
}

func lsr(e *emu6502) string {
	e.cf = (e.a & 0x01) > 0
	e.a = e.a >> 1
	e.setFlags(e.a, nil, nil)
	return "          lsr"
}

func lsr_abs(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr)
	e.cf = (v & 0x01) > 0
	v = v >> 1
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   lsr $%.4x", adr)
}

func lsr_abs_x(e *emu6502) string {
	adr, str := e.getAddress()
	v := e.getMemory(adr + uint16(e.x))
	e.cf = (v & 0x01) > 0
	v = v >> 1
	e.setMemory(adr+uint16(e.x), v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   lsr $%.4x,X", adr)
}

func lsr_zp(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr)
	e.cf = (v & 0x01) > 0
	v = v >> 1
	e.setMemory(adr, v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   lsr $%.2x", adr)
}

func lsr_zp_x(e *emu6502) string {
	adr, str := e.getZPAddress()
	v := e.getMemory(adr + uint16(e.x))
	e.cf = (v & 0x01) > 0
	v = v >> 1
	e.setMemory(adr+uint16(e.x), v)
	e.setFlags(v, nil, nil)
	return str + fmt.Sprintf("   lsr $%.2x,X", adr)
}

func rol(e *emu6502) string {
	tmp := e.cf
	e.cf = (e.a & 0x80) > 0
	e.a = e.a << 1
	if tmp {
		e.a = e.a | 0x01
	}
	e.setFlags(e.a, nil, nil)
	return "          rol"
}

func ror(e *emu6502) string {
	tmp := e.cf
	e.cf = (e.a & 0x01) > 0
	e.a = e.a >> 1
	if tmp {
		e.a = e.a | 0x80
	}
	e.setFlags(e.a, nil, nil)
	return "          ror"
}
