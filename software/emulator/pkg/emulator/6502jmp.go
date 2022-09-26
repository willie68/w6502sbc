package emulator

import "fmt"

func jmp_abs(e *Emu6502) string {
	adr, str := e.getAddress()
	e.address = adr
	return str + fmt.Sprintf("   jmp $%.4x", adr)
}

func jmp_ind(e *Emu6502) string {
	vec, str := e.getAddress()
	adr := e.readVector(vec)
	e.address = adr
	return str + fmt.Sprintf("   jmp ($%.4x)", adr)
}

func jsr_abs(e *Emu6502) string {
	adr, str := e.getAddress()
	adr--
	e.push(uint8(adr >> 8))
	e.push(uint8(adr & 0x00ff))
	e.address = adr
	return str + fmt.Sprintf("   jmp $%.4x", adr)
}

func rts(e *Emu6502) string {
	lo := e.pop()
	hi := e.pop()
	e.address = uint16(lo) + uint16(hi)*265
	e.address++
	return "         rts"
}

func rti(e *Emu6502) string {
	st := e.pop()
	e.setStatus(st)
	lo := e.pop()
	hi := e.pop()
	e.address = uint16(lo) + uint16(hi)*265
	return "         rts"
}

// badr set the new calculated address
func badr(e *Emu6502, v uint8) {
	v1 := uint16(v)
	if v > 0x80 {
		v1 = uint16(0xff00) + uint16(v)
	}
	e.address = uint16(e.address + v1)
}

func bcc(e *Emu6502) string {
	v := e.getMnemonic()
	if !e.cf {
		badr(e, v)
	}
	return fmt.Sprintf("%.2x   bcc $%.4x", v, e.address)
}

func bcs(e *Emu6502) string {
	v := e.getMnemonic()
	if e.cf {
		badr(e, v)
	}
	return fmt.Sprintf("%.2x   bcs $%.4x", v, e.address)
}

func beq(e *Emu6502) string {
	v := e.getMnemonic()
	if e.zf {
		badr(e, v)
	}
	return fmt.Sprintf("%.2x   beq $%.4x", v, e.address)
}

func bne(e *Emu6502) string {
	v := e.getMnemonic()
	if !e.zf {
		badr(e, v)
	}
	return fmt.Sprintf("%.2x   bne $%.4x", v, e.address)
}

func bpl(e *Emu6502) string {
	v := e.getMnemonic()
	if !e.nf {
		badr(e, v)
	}
	return fmt.Sprintf("%.2x   bpl $%.4x", v, e.address)
}

func bmi(e *Emu6502) string {
	v := e.getMnemonic()
	if e.nf {
		badr(e, v)
	}
	return fmt.Sprintf("%.2x   bmi $%.4x", v, e.address)
}

func bvc(e *Emu6502) string {
	v := e.getMnemonic()
	if !e.vf {
		badr(e, v)
	}
	return fmt.Sprintf("%.2x   bvc $%.4x", v, e.address)
}

func bvs(e *Emu6502) string {
	v := e.getMnemonic()
	if e.vf {
		badr(e, v)
	}
	return fmt.Sprintf("%.2x   bvs $%.4x", v, e.address)
}

func sec(e *Emu6502) string {
	e.cf = true
	return "         sec"
}

func clc(e *Emu6502) string {
	e.cf = false
	return "         clc"
}

func sei(e *Emu6502) string {
	e.jf = true
	return "         sei"
}

func cli(e *Emu6502) string {
	e.jf = false
	return "         cli"
}

func clv(e *Emu6502) string {
	e.vf = false
	return "         clv"
}

func sed(e *Emu6502) string {
	e.df = true
	return "         sed"
}

func cld(e *Emu6502) string {
	e.df = false
	return "         cld"
}
