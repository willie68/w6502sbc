package emulator

import (
	"testing"

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

func TestTransfer(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00,
	}
	e := getEmu(data)
	e.a = 0x73
	e.x = 0x00
	tax(e)
	ast.Equal(uint8(0x73), e.x)
	testFlags(ast, nil, False, False, e)

	tay(e)
	ast.Equal(uint8(0x73), e.y)
	testFlags(ast, nil, False, False, e)

	e.a = 0x00
	tax(e)
	ast.Equal(uint8(0x00), e.x)
	testFlags(ast, nil, True, False, e)

	tay(e)
	ast.Equal(uint8(0x00), e.y)
	testFlags(ast, nil, True, False, e)

	e.a = 0x83
	tax(e)
	ast.Equal(uint8(0x83), e.x)
	testFlags(ast, nil, False, True, e)

	tay(e)
	ast.Equal(uint8(0x83), e.y)
	testFlags(ast, nil, False, True, e)
}

func testFlags(ast *assert.Assertions, cf *bool, zf *bool, nf *bool, e *emu6502) {
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

func getEmu(data []uint8) *emu6502 {
	e := &emu6502{}
	e.init()
	e.highrom.data = data
	e.highrom.start = 0xe000
	e.highrom.end = e.highrom.start + uint16(len(data))
	ram := memory{
		readonly: false,
		start:    uint16(0),
		data:     make([]byte, 1024),
		end:      uint16(1023),
	}
	e.ram = ram
	e.address = 0xe000
	return e
}
