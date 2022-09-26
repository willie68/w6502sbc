package emulator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStp(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0xea, 0x02, 0x21, 0x02, 0x00, 0x21,
		0x1ffa: 0x00, 0x1ffb: 0x00,
		0x1ffc: 0x00, 0x1ffd: 0xe0,
		0x1ffe: 0x00, 0x1fff: 0x00,
	}
	e := getEmu(data)
	e.address = 0xe000
	ast.False(e.stop)
	stp(e)
	ast.True(e.stop)
	adr := e.address
	e.Step()
	ast.Equal(adr, e.address)
	e.IRQ()
	e.Step()
	ast.Equal(adr, e.address)
	e.NMI()
	e.Step()
	ast.Equal(adr, e.address)

	e.Reset()
	ast.Equal(uint16(0xe000), e.address)
}

func TestWai(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0xea, 0x02, 0x21, 0x02, 0x00, 0x21,
	}
	e := getEmu(data)
	e.address = 0xe000
	ast.False(e.wait)
	wai(e)
	ast.True(e.wait)
	adr := e.address
	e.Step()
	ast.Equal(adr, e.address)
	e.IRQ()
	e.Step()
	ast.Equal(adr+1, e.address)

	e.address = 0xe000
	ast.False(e.wait)
	wai(e)
	ast.True(e.wait)
	adr = e.address
	e.Step()
	ast.Equal(adr, e.address)
	e.NMI()
	e.Step()
	ast.Equal(adr+1, e.address)
}

func Test_mb_(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x20, 0x02, 0x21, 0x02, 0x00, 0x21,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x20), 0x00)
	e.ram.setMem(uint16(0x21), 0xff)
	smbs := []func(*Emu6502) string{
		smb0, smb1, smb2, smb3, smb4, smb5, smb6, smb7,
	}
	v := uint8(0x00)
	for i, f := range smbs {

		e.address = 0xe000
		str := f(e)
		fmt.Printf("smb%d output: %s\n\r", i, str)

		ast.Equal(uint16(0xe001), e.address)
		v = v + (0x01 << i)
		ast.Equal(v, e.ram.getMem(uint16(0x0020)))
	}

	rmbs := []func(*Emu6502) string{
		rmb0, rmb1, rmb2, rmb3, rmb4, rmb5, rmb6, rmb7,
	}
	v = uint8(0xff)
	for i, f := range rmbs {
		e.address = 0xe002
		str := f(e)
		fmt.Printf("rmb%d output: %s\n\r", i, str)

		ast.Equal(uint16(0xe003), e.address)
		v = v - (0x01 << i)
		ast.Equal(v, e.ram.getMem(uint16(0x0021)))
	}
}

func TestBb_(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x20, 0x02, 0x21, 0x02, 0x00, 0x21,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x20), 0xff)
	e.ram.setMem(uint16(0x21), 0x00)
	bbrs := []func(*Emu6502) string{
		bbr0, bbr1, bbr2, bbr3, bbr4, bbr5, bbr6, bbr7,
	}
	for i, f := range bbrs {

		e.address = 0xe000
		str := f(e)
		fmt.Printf("bbr%d output: %s\n\r", i, str)

		ast.Equal(uint16(0xe002), e.address)

		e.address = 0xe002
		str = f(e)
		fmt.Printf("bbr%d output: %s\n\r", i, str)

		ast.Equal(uint16(0xe006), e.address)
	}

	bbss := []func(*Emu6502) string{
		bbs0, bbs1, bbs2, bbs3, bbs4, bbs5, bbs6, bbs7,
	}
	for i, f := range bbss {
		e.address = 0xe000
		str := f(e)
		fmt.Printf("bbs%d output: %s\n\r", i, str)

		ast.Equal(uint16(0xe004), e.address)

		e.address = 0xe002
		str = f(e)
		fmt.Printf("bbs%d output: %s\n\r", i, str)

		ast.Equal(uint16(0xe004), e.address)
	}
}

func TestIBbx(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x20, 0x02, 0x21, 0x02, 0x00, 0x21,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x20), 0xff)
	e.ram.setMem(uint16(0x21), 0x00)

	for i := 0; i < 8; i++ {

		e.address = 0xe000
		v, _ := i_bbx(e, 0x01<<i, false)

		ast.Equal(uint16(0xe004), e.address)
		ast.Equal(uint16(0x20), v)

		e.address = 0xe002
		v, _ = i_bbx(e, 0x01<<i, false)

		ast.Equal(uint16(0xe004), e.address)
		ast.Equal(uint16(0x21), v)

		e.address = 0xe000
		v, _ = i_bbx(e, 0x01<<i, true)

		ast.Equal(uint16(0xe002), e.address)
		ast.Equal(uint16(0x20), v)

		e.address = 0xe002
		v, _ = i_bbx(e, 0x01<<i, true)

		ast.Equal(uint16(0xe006), e.address)
		ast.Equal(uint16(0x21), v)
	}
}

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

func TestBitDirect(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x73, 0x70, 0x23,
	}
	e := getEmu(data)
	e.a = 0x02
	e.cf = false
	e.zf = false
	e.nf = false
	str := bit_direct(e)
	fmt.Printf("bit output: %s\n\r", str)

	testFlags(ast, False, False, False, e)
}

func TestBitZpX(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x40, 0xe0, 0x00,
	}
	e := getEmu(data)
	e.ram.data[0x43] = 0x03
	e.x = 3

	e.a = 0x04
	e.cf = false
	e.zf = false
	e.nf = false

	str := bit_zp_x(e)
	fmt.Printf("bit output: %s\n\r", str)

	testFlags(ast, False, True, False, e)
}

func TestBitAbsX(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x00, 0xe0, 0x03, 0x04,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x7070), uint8(0x00))
	e.ram.setMem(uint16(0x7071), uint8(0xe0))
	e.x = 3

	e.a = 0x04
	e.cf = false
	e.zf = false
	e.nf = false

	str := bit_abs_x(e)
	fmt.Printf("bit output: %s\n\r", str)

	testFlags(ast, False, False, False, e)
}

func TestStaZpInd(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x70, 0x70, 0x23,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x0070), uint8(0x02))
	e.ram.setMem(uint16(0x0071), uint8(0x70))
	e.a = 0x73
	str := sta_zp_ind(e)
	fmt.Printf("sta output: %s\n\r", str)

	ast.Equal(uint8(0x73), e.ram.getMem(uint16(0x7002)))
	testFlags(ast, False, False, False, e)
}

func TestSbcZpInd(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x70, 0x70, 0x20,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x0070), uint8(0x02))
	e.ram.setMem(uint16(0x0071), uint8(0xe0))
	e.a = 0x70
	e.cf = true
	str := sbc_zp_ind(e)
	fmt.Printf("sbc output: %s\n\r", str)

	ast.Equal(uint8(0x50), e.a)
}

func TestOraZpInd(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x70, 0x70, 0x23,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x0070), uint8(0x02))
	e.ram.setMem(uint16(0x0071), uint8(0xe0))
	e.a = 0x51
	e.cf = false
	str := ora_zp_ind(e)
	fmt.Printf("ora output: %s\n\r", str)

	ast.Equal(uint8(0x73), e.a)
	testFlags(ast, False, False, False, e)
}

func TestLdaZpInd(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x70, 0x70, 0x23,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x0070), uint8(0x02))
	e.ram.setMem(uint16(0x0071), uint8(0xe0))
	str := lda_zp_ind(e)
	fmt.Printf("lda output: %s\n\r", str)

	ast.Equal(uint8(0x23), e.a)
	testFlags(ast, False, False, False, e)
}

func TestEorZpInd(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x70, 0x70, 0x23,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x0070), uint8(0x02))
	e.ram.setMem(uint16(0x0071), uint8(0xe0))
	e.a = 0x71
	e.cf = false
	str := eor_zp_ind(e)
	fmt.Printf("eor output: %s\n\r", str)

	ast.Equal(uint8(0x52), e.a)
	testFlags(ast, False, False, False, e)
}

func TestCmpZpInd(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x70, 0x70, 0x20,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x0070), uint8(0x02))
	e.ram.setMem(uint16(0x0071), uint8(0xe0))
	e.a = 0x70
	e.cf = false
	str := cmp_zp_ind(e)
	fmt.Printf("cmp output: %s\n\r", str)

	ast.Equal(uint8(0x70), e.a)
	testFlags(ast, True, False, False, e)
}

func TestAndZpInd(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x70, 0x70, 0x20,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x0070), uint8(0x02))
	e.ram.setMem(uint16(0x0071), uint8(0xe0))
	e.a = 0x70
	e.cf = false
	str := and_zp_ind(e)
	fmt.Printf("and output: %s\n\r", str)

	ast.Equal(uint8(0x20), e.a)
	testFlags(ast, False, False, False, e)
}

func TestAdcZpInd(t *testing.T) {
	ast := assert.New(t)
	data := []uint8{
		0x70, 0x70, 0x20,
	}
	e := getEmu(data)
	e.ram.setMem(uint16(0x0070), uint8(0x02))
	e.ram.setMem(uint16(0x0071), uint8(0xe0))
	e.a = 0x70
	e.cf = false
	str := adc_zp_ind(e)
	fmt.Printf("adc output: %s\n\r", str)

	ast.Equal(uint8(0x90), e.a)
	testFlags(ast, False, False, True, e)
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
