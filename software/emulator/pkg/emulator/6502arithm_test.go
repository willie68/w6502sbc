package emulator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSbc(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00,
	}
	e := getEmu(data)
	e.a = 0x30
	e.cf = true
	sbc(e, uint8(0x20))
	ast.Equal(uint8(0x10), e.a)
	testFlags(ast, True, False, False, e)

	e.a = 0x21
	e.cf = false
	sbc(e, uint8(0x20))
	ast.Equal(uint8(0x00), e.a)
	testFlags(ast, True, True, False, e)

	e.a = 0x10
	e.cf = true
	sbc(e, uint8(0x20))
	ast.Equal(uint8(0xF0), e.a)
	testFlags(ast, False, False, True, e)
}
