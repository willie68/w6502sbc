package emulator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhaPla(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00,
	}
	e := getEmu(data)
	e.sp = 0xff
	e.a = 0x67
	str := pha(e)

	fmt.Printf("pha output: %s\n\r", str)

	ast.Equal(uint8(0x67), e.pop())
	e.push(0x67)

	str = pla(e)

	fmt.Printf("pla output: %s\n\r", str)

	ast.Equal(uint8(0x67), e.a)
}

func TestPhpPlp(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00,
	}
	e := getEmu(data)
	e.setStatus(0xc7)
	e.sp = 0xff
	str := php(e)

	fmt.Printf("php output: %s\n\r", str)

	ast.Equal(uint8(0xc7), e.pop())
	e.push(0xc7)

	str = plp(e)

	fmt.Printf("pla output: %s\n\r", str)

	ast.Equal(uint8(0xc7), e.getStatus())
}

func TestNop(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00,
	}
	e := getEmu(data)
	st := e.getStatus()
	e.a = 0x67
	e.x = 0x67
	e.y = 0x67
	str := nop(e)
	fmt.Printf("php output: %s\n\r", str)

	ast.Equal(uint16(0xe000), e.address)
	ast.Equal(st, e.getStatus())
	ast.Equal(uint8(0x67), e.a)
	ast.Equal(uint8(0x67), e.x)
	ast.Equal(uint8(0x67), e.y)
}

func TestBrk(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x1ffe: 0x00, 0x1fff: 0xf0,
	}
	e := getEmu(data)
	e.setStatus(0xdf)
	str := brk(e)
	fmt.Printf("brk output: %s\n\r", str)

	ast.Equal(uint16(0xf000), e.address)
	ast.Equal(uint8(0xdf), e.pop())

	adr := uint16(e.pop()) + uint16(e.pop())*256
	ast.Equal(uint16(0xe001), adr)
}
