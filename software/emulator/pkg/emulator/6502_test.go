package emulator

import (
	"github.com/stretchr/testify/assert"
)

var True *bool
var False *bool

func init() {
	b := true
	True = &b
	c := false
	False = &c
}

func testFlags(ast *assert.Assertions, cf *bool, zf *bool, nf *bool, e *Emu6502) {
	if cf != nil {
		ast.Equal(*cf, e.cf, "carry flag")
	}
	if zf != nil {
		ast.Equal(*zf, e.zf, "zero flag")
	}
	if nf != nil {
		ast.Equal(*nf, e.nf, "negative flag")
	}
}

func getEmu(data []uint8) *Emu6502 {
	e := &Emu6502{}
	e.init()
	e.highrom.data = data
	e.highrom.start = 0xe000
	e.highrom.end = e.highrom.start + uint16(len(data)-1)
	ram := memory{
		readonly: false,
		start:    uint16(0),
		data:     make([]byte, 0x8000),
		end:      uint16(0x7fff),
	}
	e.ram = ram
	e.address = 0xe000
	e.sp = 0xff
	e.cf = false
	e.df = false
	return e
}
