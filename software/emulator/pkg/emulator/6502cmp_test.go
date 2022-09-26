package emulator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestICmp(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00,
	}
	e := getEmu(data)
	i_cmp(e, 0x00, 0x00)
	testFlags(ast, True, True, False, e)

	i_cmp(e, 0x10, 0x00)
	testFlags(ast, True, False, False, e)

	i_cmp(e, 0x00, 0x10)
	testFlags(ast, False, False, True, e)

	i_cmp(e, 0x70, 0x70)
	testFlags(ast, True, True, False, e)

	i_cmp(e, 0x10, 0x10)
	testFlags(ast, True, True, False, e)

	i_cmp(e, 0xF0, 0xF0)
	testFlags(ast, True, True, False, e)

	i_cmp(e, 0xF0, 0xE0)
	testFlags(ast, True, False, False, e)
}

func TestCmpDirect(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x70, 0x70, 0x20,
	}
	e := getEmu(data)
	e.a = 0x70
	str := cmp_direct(e)
	fmt.Printf("cmp output: %s\n\r", str)

	testFlags(ast, True, True, False, e)

	e.address = 0xe000
	e.a = 0x00
	e.x = 0x70
	str = cpx_direct(e)
	fmt.Printf("cpx output: %s\n\r", str)

	testFlags(ast, True, True, False, e)

	e.address = 0xe000
	e.a = 0x00
	e.x = 0x00
	e.y = 0x70
	str = cpy_direct(e)
	fmt.Printf("cpy output: %s\n\r", str)

	testFlags(ast, True, True, False, e)
}

func TestCmpAbs_(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x02, 0xe0, 0x70, 0x80, 0x75, 0x60,
	}
	e := getEmu(data)
	e.a = 0x70
	str := cmp_abs(e)
	fmt.Printf("cmp output: %s\n\r", str)

	testFlags(ast, True, True, False, e)

	e.address = 0xe000
	e.a = 0x75
	e.x = 0x02
	str = cmp_abs_x(e)
	fmt.Printf("cmp output: %s\n\r", str)

	testFlags(ast, True, True, False, e)

	e.address = 0xe000
	e.a = 0x60
	e.y = 0x03
	str = cmp_abs_y(e)
	fmt.Printf("cmp output: %s\n\r", str)

	testFlags(ast, True, True, False, e)

	e.address = 0xe000
	e.a = 0x00
	e.x = 0x70
	str = cpx_abs(e)
	fmt.Printf("cpx output: %s\n\r", str)

	testFlags(ast, True, True, False, e)

	e.address = 0xe000
	e.a = 0x00
	e.x = 0x00
	e.y = 0x70
	str = cpy_abs(e)
	fmt.Printf("cpy output: %s\n\r", str)

	testFlags(ast, True, True, False, e)
}

func TestCmpZp_(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x20, 0xe0, 0x70, 0x80, 0x75, 0x60,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x0020), uint8(0x70))
	e.ram.setMem(uint16(0x0022), uint8(0x75))
	e.a = 0x70
	str := cmp_zp(e)
	fmt.Printf("cmp output: %s\n\r", str)

	testFlags(ast, True, True, False, e)

	e.address = 0xe000
	e.a = 0x75
	e.x = 0x02
	str = cmp_zp_x(e)
	fmt.Printf("cmp output: %s\n\r", str)

	testFlags(ast, True, True, False, e)

	e.address = 0xe000
	e.a = 0x00
	e.x = 0x70
	str = cpx_zp(e)
	fmt.Printf("cpx output: %s\n\r", str)

	testFlags(ast, True, True, False, e)

	e.address = 0xe000
	e.a = 0x00
	e.x = 0x00
	e.y = 0x70
	str = cpy_zp(e)
	fmt.Printf("cpy output: %s\n\r", str)

	testFlags(ast, True, True, False, e)
}

func TestCmpInd_(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x20, 0xe0, 0x70, 0x80, 0x75, 0x60,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x0023), uint8(0x10))
	e.ram.setMem(uint16(0x0024), uint8(0x01))
	e.ram.setMem(uint16(0x0110), uint8(0x75))
	e.a = 0x75
	e.x = 3

	str := cmp_ind_x(e)
	fmt.Printf("cmp output: %s\n\r", str)

	testFlags(ast, True, True, False, e)

	e.address = 0xe000
	e.ram.setMem(uint16(0x0023), uint8(0x00))
	e.ram.setMem(uint16(0x0024), uint8(0x00))
	e.ram.setMem(uint16(0x0110), uint8(0x00))

	e.ram.setMem(uint16(0x0020), uint8(0x10))
	e.ram.setMem(uint16(0x0021), uint8(0x01))
	e.ram.setMem(uint16(0x0113), uint8(0x75))
	e.a = 0x75
	e.x = 0
	e.y = 3

	str = cmp_ind_y(e)
	fmt.Printf("cmp output: %s\n\r", str)

	testFlags(ast, True, True, False, e)
}
