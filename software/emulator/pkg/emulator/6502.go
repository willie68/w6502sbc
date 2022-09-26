package emulator

import (
	"fmt"

	log "github.com/willie68/w6502sbc/tree/main/software/emulator/internal/logging"
)

type Build6502 struct {
	memmap []memory
	c02    bool
}

type Emu6502 struct {
	functions []func(*Emu6502) string
	memmap    []memory

	a, x, y uint8
	sp      uint8
	address uint16
	cf      bool // carry flag
	zf      bool // zero flag
	jf      bool // iterrupt flag
	df      bool // decimal flag
	bf      bool // break flag
	vf      bool // overflow flag
	nf      bool // negative flag
	c02     bool // is cmos cpu
	wait    bool // 65C02 wait until iterrupt
	stop    bool // 65C02 stop until reset
}

func NewEmu6502() *Build6502 {
	return &Build6502{
		memmap: make([]memory, 0),
		c02:    false,
	}
}

func (b *Build6502) WithRAM(start, end uint16) *Build6502 {
	ram := make([]byte, end-start)
	b.memmap = append(b.memmap, memory{readonly: false, start: start, data: ram, end: end})
	return b
}

func (b *Build6502) WithROM(start uint16, data []byte) *Build6502 {
	b.memmap = append(b.memmap, memory{readonly: true, start: start, end: start + uint16(len(data)) - 1, data: data})
	return b
}

func (b *Build6502) With6522(adr uint16) *Build6502 {
	return b
}

func (b *Build6502) With65C02() *Build6502 {
	b.c02 = true
	return b
}

func (b *Build6502) Build() Emu6502 {
	f := functions
	if b.c02 {
		f = append(f, cfunc...)
	}
	emu := Emu6502{
		c02:       true,
		functions: f,
		memmap:    b.memmap,
		wait:      false,
		stop:      false,
	}
	emu.init()
	return emu
}

func (e *Emu6502) init() {
	e.a = 0
	e.x = 0
	e.y = 0
}

func (e *Emu6502) Start() string {
	log.Logger.Info("starting emulation")
	return e.Reset()
}

func (e *Emu6502) Reset() string {
	e.jf = true
	if e.c02 {
		e.df = false
		e.wait = false
		e.stop = false
	}
	e.address = e.readVector(uint16(0xFFFC))
	return fmt.Sprintf("read vector $FFFC with value of $%.4x", e.address)
}

func (e *Emu6502) NMI() {
	e.wait = false
}

func (e *Emu6502) IRQ() {
	e.wait = false
	if e.jf {
		adr, _ := e.getAddress()
		e.push(uint8(adr >> 8))
		e.push(uint8(adr & 0x00ff))
		st := e.getStatus()
		e.push(st)
	}
}

func (e *Emu6502) Step() string {
	if e.wait {
		return "no step possible, cpu wait, need reset or iterrupt"
	}
	if e.stop {
		return "no step possible, cpu stoped, need reset"
	}
	output := ""
	output = fmt.Sprintf("$%.4x ", e.address)
	mne := e.getMnemonic()
	output = output + fmt.Sprintf("%.2x ", mne)

	fct := functions[mne]
	if fct != nil {
		res := fct(e)
		output = output + res
	} else {
		switch mne {
		default:
			output = output + "           illegal opcode"
		}
	}
	return output
}

func (e *Emu6502) Adr() uint16 {
	return e.address
}

func (e *Emu6502) SP() uint8 {
	return e.sp
}

func (e *Emu6502) A() uint8 {
	return e.a
}

func (e *Emu6502) X() uint8 {
	return e.x
}

func (e *Emu6502) Y() uint8 {
	return e.y
}

func (e *Emu6502) ST() uint8 {
	return e.getStatus()
}

func (e *Emu6502) readVector(adr uint16) uint16 {
	lo := uint16(e.getMemory(adr))
	hi := uint16(e.getMemory(adr + 1))
	return hi*256 + lo
}

func (e *Emu6502) getMemory(adr uint16) uint8 {
	for _, m := range e.memmap {
		if m.IsMapped(adr) {
			return m.GetMem(adr)
		}
	}
	return 0
}

func (e *Emu6502) setMemory(adr uint16, dt uint8) {
	for _, m := range e.memmap {
		if m.IsMapped(adr) {
			m.SetMem(adr, dt)
		}
	}
}

func (e *Emu6502) getMnemonic() uint8 {
	b := e.getMemory(e.address)
	e.address++
	return b
}

func (e *Emu6502) setFlags(v uint8, cf *bool, vf *bool) {
	e.zf = v == 0
	e.nf = (v & 0x80) > 0
	if cf != nil {
		e.cf = *cf
	}
	if vf != nil {
		e.vf = *vf
	}
}

func (e *Emu6502) getAddress() (uint16, string) {
	lo := e.getMnemonic()
	hi := e.getMnemonic()
	return uint16(hi)*256 + uint16(lo), fmt.Sprintf("%.2x %.2x", lo, hi)
}

func (e *Emu6502) getZPAddress() (uint16, string) {
	lo := e.getMnemonic()
	return uint16(lo), fmt.Sprintf("%.2x   ", lo)
}

func (e *Emu6502) push(v uint8) {
	adr := uint16(0x0100) + uint16(e.sp)
	e.setMemory(adr, v)
	e.sp--
}

func (e *Emu6502) pop() uint8 {
	e.sp++
	adr := uint16(0x0100) + uint16(e.sp)
	return e.getMemory(adr)
}

func (e *Emu6502) setStatus(st uint8) {
	e.nf = (st & 0x80) > 0
	e.vf = (st & 0x40) > 0
	e.nf = true
	e.bf = (st & 0x10) > 0
	e.df = (st & 0x08) > 0
	e.jf = (st & 0x04) > 0
	e.zf = (st & 0x02) > 0
	e.cf = (st & 0x01) > 0
}

func (e *Emu6502) getStatus() uint8 {
	st := uint8(0)
	if e.nf {
		st = st + 0x80
	}
	if e.vf {
		st = st + 0x40
	}
	if e.bf {
		st = st + 0x10
	}
	if e.df {
		st = st + 0x08
	}
	if e.jf {
		st = st + 0x04
	}
	if e.zf {
		st = st + 0x02
	}
	if e.cf {
		st = st + 0x01
	}
	return st
}
