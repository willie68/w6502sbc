package emulator

import (
	"fmt"

	log "github.com/willie68/w6502sbc/tree/main/software/emulator/internal/logging"
)

type build6502 struct {
	highrom memory
	ram     memory
}

type emu6502 struct {
	highrom memory
	ram     memory

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
}

type memory struct {
	readonly bool
	start    uint16
	end      uint16
	data     []byte
}

func (m *memory) getMem(adr uint16) uint8 {
	return m.data[adr]
}

func (m *memory) setMem(adr uint16, dt uint8) {
	if !m.readonly {
		m.data[adr] = dt
	}
}

func NewEmu6502() build6502 {
	return build6502{}
}

func (b build6502) WithRAM(start, end uint16) build6502 {
	ram := make([]byte, end-start)
	b.ram = memory{readonly: false, start: start, data: ram, end: end}
	return b
}

func (b build6502) WithROM(start uint16, data []byte) build6502 {
	b.highrom = memory{readonly: true, start: start, end: start + uint16(len(data)) - 1, data: data}
	return b
}

func (b build6502) With6522(adr uint16) build6502 {
	return b
}

func (b build6502) Build() emu6502 {
	emu := emu6502{
		highrom: b.highrom,
		ram:     b.ram,
	}
	emu.init()
	return emu
}

func (e *emu6502) init() {
	e.a = 0
	e.x = 0
	e.y = 0
}

func (e *emu6502) Start() {
	log.Logger.Info("starting emulation")
	e.address = e.readVector(uint16(0xFFFC))
}

func (e *emu6502) Reset() {
}

func (e *emu6502) NMI() {
}

func (e *emu6502) IRQ() {
	adr, _ := e.getAddress()
	e.push(uint8(adr >> 8))
	e.push(uint8(adr & 0x00ff))
	st := e.getStatus()
	e.push(st)
}

func (e *emu6502) Step() string {
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

func (e *emu6502) Adr() uint16 {
	return e.address
}

func (e *emu6502) SP() uint8 {
	return e.sp
}

func (e *emu6502) A() uint8 {
	return e.a
}

func (e *emu6502) X() uint8 {
	return e.x
}

func (e *emu6502) Y() uint8 {
	return e.y
}

func (e *emu6502) readVector(adr uint16) uint16 {
	lo := uint16(e.getMemory(adr))
	hi := uint16(e.getMemory(adr + 1))
	return hi*256 + lo
}

func (e *emu6502) getMemory(adr uint16) uint8 {
	if adr >= e.ram.start && adr <= e.ram.end {
		ramadr := adr - e.ram.start
		return e.ram.getMem(ramadr)
	}
	if adr >= e.highrom.start && adr <= e.highrom.end {
		romadr := adr - e.highrom.start
		return e.highrom.getMem(romadr)
	}
	return 0
}

func (e *emu6502) getMnemonic() uint8 {
	b := e.getMemory(e.address)
	e.address++
	return b
}

func (e *emu6502) setMemory(adr uint16, dt uint8) {
	if adr >= e.ram.start && adr <= e.ram.end {
		ramadr := adr - e.ram.start
		e.ram.setMem(ramadr, dt)
	}
}

func (e *emu6502) setFlags(v uint8, cf *bool, vf *bool) {
	e.zf = v == 0
	e.nf = (v & 0x80) > 0
	if cf != nil {
		e.cf = *cf
	}
	if vf != nil {
		e.vf = *vf
	}
}

func (e *emu6502) getAddress() (uint16, string) {
	lo := e.getMnemonic()
	hi := e.getMnemonic()
	return uint16(hi)*256 + uint16(lo), fmt.Sprintf("%.2x %.2x", lo, hi)
}

func (e *emu6502) getZPAddress() (uint16, string) {
	lo := e.getMnemonic()
	return uint16(lo), fmt.Sprintf("%.2x   ", lo)
}

func (e *emu6502) push(v uint8) {
	adr := uint16(0x0100) + uint16(e.sp)
	e.setMemory(adr, v)
	e.sp--
}

func (e *emu6502) pop() uint8 {
	e.sp++
	adr := uint16(0x0100) + uint16(e.sp)
	return e.getMemory(adr)
}

func (e *emu6502) setStatus(st uint8) {
	e.nf = (st & 0x80) > 0
	e.vf = (st & 0x40) > 0
	e.nf = true
	e.bf = (st & 0x10) > 0
	e.df = (st & 0x08) > 0
	e.jf = (st & 0x04) > 0
	e.zf = (st & 0x02) > 0
	e.cf = (st & 0x01) > 0
}

func (e *emu6502) getStatus() uint8 {
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
