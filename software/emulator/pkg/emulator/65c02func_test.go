package emulator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTsb(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x20, 0x00, 0x21, 0x20, 0x00, 0x21,
	}
	e := getEmu(data)

	e.a = 0x40
	v := i_tsb(e, 0x80)

	ast.Equal(uint8(0xC0), v)

	for i := uint16(0); i < 0xff; i++ {
		e.ram.setMem(i, 0x0f)
	}
	ast.Equal(uint8(0x0f), e.ram.getMem(uint16(0x0020)))

	str := tsb_abs(e)

	fmt.Printf("trb output: %s\n\r", str)
	ast.Equal(uint8(0x4f), e.ram.getMem(uint16(0x0020)))

	str = tsb_zp(e)

	fmt.Printf("trb output: %s\n\r", str)
	ast.Equal(uint8(0x4f), e.ram.getMem(uint16(0x0021)))
}

func TestTrb(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x20, 0x00, 0x21, 0x20, 0x00, 0x21,
	}
	e := getEmu(data)

	e.a = 0x40
	v := i_trb(e, 0xc0)

	ast.Equal(uint8(0x80), v)

	for i := uint16(0); i < 0xff; i++ {
		e.ram.setMem(i, 0xff)
	}
	ast.Equal(uint8(0xff), e.ram.getMem(uint16(0x0020)))

	str := trb_abs(e)

	fmt.Printf("trb output: %s\n\r", str)
	ast.Equal(uint8(0xbf), e.ram.getMem(uint16(0x0020)))

	str = trb_zp(e)

	fmt.Printf("trb output: %s\n\r", str)
	ast.Equal(uint8(0xbf), e.ram.getMem(uint16(0x0021)))
}

func TestStz(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x20, 0x00, 0x21, 0x20, 0x00, 0x21,
	}
	e := getEmu(data)
	for i := uint16(0); i < 0xff; i++ {
		e.ram.setMem(i, 0xff)
	}
	ast.Equal(uint8(0xff), e.ram.getMem(uint16(0x0020)))

	str := stz_abs(e)

	fmt.Printf("stz output: %s\n\r", str)
	ast.Equal(uint8(0x00), e.ram.getMem(uint16(0x0020)))

	str = stz_zp(e)

	fmt.Printf("stz output: %s\n\r", str)
	ast.Equal(uint8(0x00), e.ram.getMem(uint16(0x0021)))

	e.x = 0x20

	str = stz_abs_x(e)

	fmt.Printf("stz output: %s\n\r", str)
	ast.Equal(uint8(0x00), e.ram.getMem(uint16(0x0040)))

	str = stz_zp_x(e)

	fmt.Printf("stz output: %s\n\r", str)
	ast.Equal(uint8(0x00), e.ram.getMem(uint16(0x0041)))
}

func TestPhxPlx(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00,
	}
	e := getEmu(data)
	e.sp = 0xff
	e.x = 0x67
	str := phx(e)

	fmt.Printf("phx output: %s\n\r", str)

	ast.Equal(uint8(0x67), e.pop())
	e.push(0x67)

	str = plx(e)

	fmt.Printf("plx output: %s\n\r", str)

	ast.Equal(uint8(0x67), e.x)
}

func TestPhyPly(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00,
	}
	e := getEmu(data)
	e.sp = 0xff
	e.y = 0x67
	str := phy(e)

	fmt.Printf("phy output: %s\n\r", str)

	ast.Equal(uint8(0x67), e.pop())
	e.push(0x67)

	str = ply(e)

	fmt.Printf("ply output: %s\n\r", str)

	ast.Equal(uint8(0x67), e.y)
}

func TestBra(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00, 0x02, 0xfe,
	}
	e := getEmu(data)
	e.address = 0xe001
	str := bra(e)

	fmt.Printf("bra output: %s\n\r", str)

	ast.Equal(uint16(0xe004), e.address)

	e.address = 0xe002
	str = bra(e)

	fmt.Printf("bra output: %s\n\r", str)

	ast.Equal(uint16(0xe001), e.address)
}

func TestInc(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00, 0x80, 0x23,
	}
	e := getEmu(data)

	e.a = 0x04
	str := inc(e)
	fmt.Printf("inc output: %s\n\r", str)

	ast.Equal(uint8(0x05), e.a)
	testFlags(ast, nil, False, False, e)

	e.a = 0x7f
	str = inc(e)
	fmt.Printf("inc output: %s\n\r", str)

	ast.Equal(uint8(0x80), e.a)
	testFlags(ast, nil, False, True, e)

	e.a = 0xff
	str = inc(e)
	fmt.Printf("inc output: %s\n\r", str)

	ast.Equal(uint8(0x00), e.a)
	testFlags(ast, nil, True, False, e)
}

func TestDec(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00, 0x80, 0x23,
	}
	e := getEmu(data)

	e.a = 0x01
	str := dec(e)
	fmt.Printf("dec output: %s\n\r", str)

	ast.Equal(uint8(0x00), e.a)
	testFlags(ast, nil, True, False, e)

	e.a = 0x80
	str = dec(e)
	fmt.Printf("dec output: %s\n\r", str)

	ast.Equal(uint8(0x7f), e.a)
	testFlags(ast, nil, False, False, e)

	e.a = 0x00
	str = dec(e)
	fmt.Printf("dec output: %s\n\r", str)

	ast.Equal(uint8(0xff), e.a)
	testFlags(ast, nil, False, True, e)
}

func TestJmpAbsX(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00, 0x70, 0x23,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x7070), uint8(0x70))
	e.ram.setMem(uint16(0x7071), uint8(0x80))
	e.x = 0x70
	str := jmp_abs_x(e)
	fmt.Printf("jmp output: %s\n\r", str)

	ast.Equal(uint16(0x8070), e.address)
}
