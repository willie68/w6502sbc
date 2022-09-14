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
	zf      bool // zero flag
	cf      bool // carry flag
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
}

func (e *emu6502) Step() string {
	output := ""
	output = fmt.Sprintf("$%.4x ", e.address)
	cmd := e.getMnemonic()
	output = output + fmt.Sprintf("%.2x ", cmd)

	switch cmd {
	case 0x8d: // STA $0000

		lo := e.getMnemonic()
		hi := e.getMnemonic()
		adr := uint16(hi)*256 + uint16(lo)
		e.setMemory(adr, e.a)
		output = output + fmt.Sprintf("%.2x %.2x     sta $%.4x", lo, hi, adr)
	case 0xa9: // LDA #
		e.a = e.getMnemonic()
		e.zf = e.a == 0
		e.nf = (e.a & 0x80) > 0
		output = output + fmt.Sprintf("%.2x        lda #$%.2x", e.a, e.a)
	case 0xea: // NOP
		break
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
