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
}
