package emulator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLDA_direct(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00, 0x80, 0x23,
	}
	e := getEmu(data)

	str := lda_direct(e)
	fmt.Printf("lda output: %s\n\r", str)

	ast.Equal(data[0], e.a)
	ast.Equal(uint16(0xe001), e.address)
	testFlags(ast, nil, True, False, e)

	str = lda_direct(e)
	fmt.Printf("lda output: %s\n\r", str)

	ast.Equal(data[1], e.a)
	testFlags(ast, nil, False, True, e)

	str = lda_direct(e)
	fmt.Printf("lda output: %s\n\r", str)

	ast.Equal(data[2], e.a)
	testFlags(ast, nil, False, False, e)

	e.address = 0xe000
	str = ldy_direct(e)
	fmt.Printf("ldy output: %s\n\r", str)

	ast.Equal(data[0], e.y)
	ast.Equal(uint16(0xe001), e.address)
	testFlags(ast, nil, True, False, e)
}

func TestLDA_abs(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x02, 0xe0, 0x00,
		0x05, 0xe0, 0x80,
		0x08, 0xe0, 0x23,
	}
	e := getEmu(data)

	str := lda_abs(e)
	fmt.Printf("lda output: %s\n\r", str)

	ast.Equal(data[2], e.a)
	ast.Equal(uint16(0xe002), e.address)
	testFlags(ast, nil, True, False, e)

	e.address = 0xe003
	str = lda_abs(e)
	fmt.Printf("lda output: %s\n\r", str)

	ast.Equal(data[5], e.a)
	testFlags(ast, nil, False, True, e)

	e.address = 0xe006
	str = lda_abs(e)
	fmt.Printf("lda output: %s\n\r", str)

	ast.Equal(data[8], e.a)
	testFlags(ast, nil, False, False, e)

	e.address = 0xe000
	str = ldy_abs(e)
	fmt.Printf("ldy output: %s\n\r", str)

	ast.Equal(data[2], e.y)
	ast.Equal(uint16(0xe002), e.address)
	testFlags(ast, nil, True, False, e)
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
	str := lda_abs_x(e)
	fmt.Printf("lda output: %s\n\r", str)

	ast.Equal(data[5], e.a)
	ast.Equal(uint16(0xe002), e.address)
	testFlags(ast, nil, False, True, e)

	e.address = 0xe000
	str = ldy_abs_x(e)
	fmt.Printf("ldy output: %s\n\r", str)

	ast.Equal(data[5], e.y)
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
	str := lda_abs_y(e)
	fmt.Printf("lda output: %s\n\r", str)

	ast.Equal(data[5], e.a)
	ast.Equal(uint16(0xe002), e.address)
	testFlags(ast, nil, False, True, e)

	e.address = 0xe000
	str = ldx_abs_y(e)
	fmt.Printf("ldx output: %s\n\r", str)

	ast.Equal(data[5], e.x)
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
	str := lda_zp(e)
	fmt.Printf("lda output: %s\n\r", str)

	ast.Equal(uint8(3), e.a)
	ast.Equal(uint16(0xe001), e.address)
	testFlags(ast, nil, False, False, e)

	e.address = 0xe000
	str = ldx_zp(e)
	fmt.Printf("ldx output: %s\n\r", str)

	ast.Equal(uint8(3), e.x)
	ast.Equal(uint16(0xe001), e.address)
	testFlags(ast, nil, False, False, e)

	e.address = 0xe000
	str = ldy_zp(e)
	fmt.Printf("ldy output: %s\n\r", str)

	ast.Equal(uint8(3), e.y)
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
	str := lda_zp_x(e)
	fmt.Printf("lda output: %s\n\r", str)

	ast.Equal(uint8(3), e.a)
	ast.Equal(uint16(0xe001), e.address)
	testFlags(ast, nil, False, False, e)

	e.address = 0xe000
	str = ldy_zp_x(e)
	fmt.Printf("ldy output: %s\n\r", str)

	ast.Equal(uint8(3), e.y)
	ast.Equal(uint16(0xe001), e.address)
	testFlags(ast, nil, False, False, e)
}

func TestLDX_zp_y(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x40, 0xe0, 0x00,
	}
	e := getEmu(data)
	e.ram.data[0x43] = 3
	e.y = 3
	str := ldx_zp_y(e)
	fmt.Printf("ldx output: %s\n\r", str)

	ast.Equal(uint8(3), e.x)
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

func TestTransfer(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00,
	}
	e := getEmu(data)
	e.a = 0x73
	e.x = 0x00
	str := tax(e)
	fmt.Printf("tax output: %s\n\r", str)
	ast.Equal(uint8(0x73), e.x)
	testFlags(ast, nil, False, False, e)

	str = tay(e)
	fmt.Printf("tay output: %s\n\r", str)
	ast.Equal(uint8(0x73), e.y)
	testFlags(ast, nil, False, False, e)

	e.a = 0x00
	str = tax(e)
	fmt.Printf("tax output: %s\n\r", str)
	ast.Equal(uint8(0x00), e.x)
	testFlags(ast, nil, True, False, e)

	str = tay(e)
	fmt.Printf("tay output: %s\n\r", str)
	ast.Equal(uint8(0x00), e.y)
	testFlags(ast, nil, True, False, e)

	e.a = 0x83
	str = tax(e)
	fmt.Printf("tax output: %s\n\r", str)
	ast.Equal(uint8(0x83), e.x)
	testFlags(ast, nil, False, True, e)

	str = tay(e)
	fmt.Printf("tay output: %s\n\r", str)
	ast.Equal(uint8(0x83), e.y)
	testFlags(ast, nil, False, True, e)

	e.x = 0x83
	str = txa(e)
	fmt.Printf("txa output: %s\n\r", str)
	ast.Equal(uint8(0x83), e.a)

	e.y = 0x32
	str = tya(e)
	fmt.Printf("tya output: %s\n\r", str)
	ast.Equal(uint8(0x32), e.a)

	v := e.sp
	str = tsx(e)
	fmt.Printf("tsx output: %s\n\r", str)
	ast.Equal(v, e.x)

	e.x = 0xff
	str = txs(e)
	fmt.Printf("txs output: %s\n\r", str)
	ast.Equal(uint8(0xff), e.sp)
}

func TestSTAXY_abs(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00, 0x70, 0x00,
		0x05, 0xe0, 0x80,
		0x08, 0xe0, 0x23,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x7000), 0)
	e.a = 0x67
	str := sta_abs(e)
	fmt.Printf("sta output: %s\n\r", str)

	ast.Equal(uint8(0x67), e.ram.getMem(uint16(0x7000)))
	ast.Equal(uint16(0xe002), e.address)

	e.address = 0xe000
	e.x = 0x63
	str = stx_abs(e)
	fmt.Printf("stx output: %s\n\r", str)

	ast.Equal(uint8(0x63), e.ram.getMem(uint16(0x7000)))
	ast.Equal(uint16(0xe002), e.address)

	e.address = 0xe000
	e.y = 0x56
	str = sty_abs(e)
	fmt.Printf("sty output: %s\n\r", str)

	ast.Equal(uint8(0x56), e.ram.getMem(uint16(0x7000)))
	ast.Equal(uint16(0xe002), e.address)
}

func TestSTA_abs_x(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00, 0x70, 0x00,
		0x05, 0xe0, 0x80,
		0x08, 0xe0, 0x23,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x7000), 0)
	e.a = 0x67
	e.x = 0x03
	str := sta_abs_x(e)
	fmt.Printf("sta output: %s\n\r", str)

	ast.Equal(uint8(0x67), e.ram.getMem(uint16(0x7003)))
	ast.Equal(uint16(0xe002), e.address)
}

func TestSTA_abs_y(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00, 0x70, 0x00,
		0x05, 0xe0, 0x80,
		0x08, 0xe0, 0x23,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x7000), 0)
	e.a = 0x67
	e.y = 0x03
	str := sta_abs_y(e)
	fmt.Printf("sta output: %s\n\r", str)

	ast.Equal(uint8(0x67), e.ram.getMem(uint16(0x7003)))
	ast.Equal(uint16(0xe002), e.address)
}

func TestSTAXY_zp(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x70, 0x70, 0x00,
		0x05, 0xe0, 0x80,
		0x08, 0xe0, 0x23,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x0070), 0)
	e.a = 0x67
	str := sta_zp(e)
	fmt.Printf("sta output: %s\n\r", str)

	ast.Equal(uint8(0x67), e.ram.getMem(uint16(0x0070)))
	ast.Equal(uint16(0xe001), e.address)

	e.address = 0xe000
	e.x = 0x63
	str = stx_zp(e)
	fmt.Printf("stx output: %s\n\r", str)

	ast.Equal(uint8(0x63), e.ram.getMem(uint16(0x0070)))
	ast.Equal(uint16(0xe001), e.address)

	e.address = 0xe000
	e.y = 0x56
	str = sty_zp(e)
	fmt.Printf("sty output: %s\n\r", str)

	ast.Equal(uint8(0x56), e.ram.getMem(uint16(0x0070)))
	ast.Equal(uint16(0xe001), e.address)
}

func TestSTAY_zp_x(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x70, 0x70, 0x00,
		0x05, 0xe0, 0x80,
		0x08, 0xe0, 0x23,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x0070), 0)
	e.a = 0x67
	e.x = 0x03
	str := sta_zp_x(e)
	fmt.Printf("sta output: %s\n\r", str)

	ast.Equal(uint8(0x67), e.ram.getMem(uint16(0x0073)))
	ast.Equal(uint16(0xe001), e.address)

	e.address = 0xe000
	e.y = 0x56
	e.x = 0x03
	str = sty_zp_x(e)
	fmt.Printf("sty output: %s\n\r", str)

	ast.Equal(uint8(0x56), e.ram.getMem(uint16(0x0073)))
	ast.Equal(uint16(0xe001), e.address)
}

func TestSTX_zp_y(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x70, 0x70, 0x00,
		0x05, 0xe0, 0x80,
		0x08, 0xe0, 0x23,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x0070), 0)
	e.x = 0x63
	e.y = 0x03
	str := stx_zp_y(e)
	fmt.Printf("stx output: %s\n\r", str)

	ast.Equal(uint8(0x63), e.ram.getMem(uint16(0x0073)))
	ast.Equal(uint16(0xe001), e.address)
}

func TestSTA_ind_x(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x70, 0x70, 0x00,
		0x05, 0xe0, 0x80,
		0x08, 0xe0, 0x23,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x0074), 0x70)
	e.ram.setMem(uint16(0x0075), 0x70)
	e.a = 0x67
	e.x = 0x04
	str := sta_ind_x(e)
	fmt.Printf("sta output: %s\n\r", str)

	ast.Equal(uint8(0x67), e.ram.getMem(uint16(0x7070)))
	ast.Equal(uint16(0xe001), e.address)
}

func TestSTA_ind_y(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x70, 0x70, 0x00,
		0x05, 0xe0, 0x80,
		0x08, 0xe0, 0x23,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x0070), 0x70)
	e.ram.setMem(uint16(0x0071), 0x70)
	e.a = 0x67
	e.y = 0x04
	str := sta_ind_y(e)
	fmt.Printf("sta output: %s\n\r", str)

	ast.Equal(uint8(0x67), e.ram.getMem(uint16(0x7074)))
	ast.Equal(uint16(0xe001), e.address)
}
