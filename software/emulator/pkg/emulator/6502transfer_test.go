package emulator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLDA_direct(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00, 0x80, 0x23,
	}
	e := getEmu(data)

	lda_direct(e)

	ast.Equal(data[0], e.a)
	ast.Equal(uint16(0xe001), e.address)
	testFlags(ast, nil, True, False, e)

	lda_direct(e)

	ast.Equal(data[1], e.a)
	testFlags(ast, nil, False, True, e)

	lda_direct(e)

	ast.Equal(data[2], e.a)
	testFlags(ast, nil, False, False, e)
}

func TestLDA_abs(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x02, 0xe0, 0x00,
		0x05, 0xe0, 0x80,
		0x08, 0xe0, 0x23,
	}
	e := getEmu(data)

	lda_abs(e)

	ast.Equal(data[2], e.a)
	ast.Equal(uint16(0xe002), e.address)
	testFlags(ast, nil, True, False, e)

	e.address = 0xe003
	lda_abs(e)

	ast.Equal(data[5], e.a)
	testFlags(ast, nil, False, True, e)

	e.address = 0xe006
	lda_abs(e)

	ast.Equal(data[8], e.a)
	testFlags(ast, nil, False, False, e)
}

func TestLDA_abs_x(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x02, 0xe0, 0x00,
		0x05, 0xe0, 0x80,
		0x08, 0xe0, 0x23,
	}
	e := getEmu(data)
	e.x = 3
	lda_abs_x(e)

	ast.Equal(data[5], e.a)
	ast.Equal(uint16(0xe002), e.address)
	testFlags(ast, nil, False, True, e)
}

func TestLDA_abs_y(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x02, 0xe0, 0x00,
		0x05, 0xe0, 0x80,
		0x08, 0xe0, 0x23,
	}
	e := getEmu(data)
	e.y = 3
	lda_abs_y(e)

	ast.Equal(data[5], e.a)
	ast.Equal(uint16(0xe002), e.address)
	testFlags(ast, nil, False, True, e)
}

func TestLDA_zp(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x43, 0xe0, 0x00,
	}
	e := getEmu(data)
	e.ram.data[0x43] = 3
	lda_zp(e)

	ast.Equal(uint8(3), e.a)
	ast.Equal(uint16(0xe001), e.address)
	testFlags(ast, nil, False, False, e)
}

func TestLDA_zp_x(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x40, 0xe0, 0x00,
	}
	e := getEmu(data)
	e.ram.data[0x43] = 3
	e.x = 3
	lda_zp_x(e)

	ast.Equal(uint8(3), e.a)
	ast.Equal(uint16(0xe001), e.address)
	testFlags(ast, nil, False, False, e)
}

func TestLDA_ind_x(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x40, 0xe0, 0x00,
	}
	e := getEmu(data)
	e.ram.data[0x43] = 0x00
	e.ram.data[0x44] = 0xe0
	e.x = 3
	lda_ind_x(e)

	ast.Equal(data[0], e.a)
	ast.Equal(uint16(0xe001), e.address)
	testFlags(ast, nil, False, False, e)
}

func TestLDA_ind_y(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x40, 0xe0, 0x88,
	}
	e := getEmu(data)
	e.ram.data[0x40] = 0x00
	e.ram.data[0x41] = 0xe0
	e.y = 2
	lda_ind_y(e)

	ast.Equal(data[2], e.a)
	ast.Equal(uint16(0xe001), e.address)
	testFlags(ast, nil, False, True, e)
}

func TestLDXY_direct(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00, 0x80, 0x23,
	}
	e := getEmu(data)

	ldx_direct(e)

	ast.Equal(data[0], e.x)
	ast.Equal(uint16(0xe001), e.address)
	testFlags(ast, nil, True, False, e)

	ldy_direct(e)

	ast.Equal(data[1], e.y)
	testFlags(ast, nil, False, True, e)

	lda_direct(e)

	ast.Equal(data[2], e.a)
	testFlags(ast, nil, False, False, e)
}

func TestLDXY_abs(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x02, 0xe0, 0x00,
		0x05, 0xe0, 0x80,
		0x08, 0xe0, 0x23,
	}
	e := getEmu(data)

	ldx_abs(e)

	ast.Equal(data[2], e.x)
	ast.Equal(uint16(0xe002), e.address)
	testFlags(ast, nil, True, False, e)

	e.address = 0xe003
	ldy_abs(e)

	ast.Equal(data[5], e.y)
	testFlags(ast, nil, False, True, e)

	e.address = 0xe006
	lda_abs(e)

	ast.Equal(data[8], e.a)
	testFlags(ast, nil, False, False, e)
}
