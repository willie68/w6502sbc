package emulator

type Mm interface {
	IsMapped(adr uint16) bool
	GetMem(adr uint16) uint8
	SetMem(adr uint16, dt uint8)
}

type memory struct {
	readonly bool
	start    uint16
	end      uint16
	data     []byte
}

func (m *memory) GetMem(adr uint16) uint8 {
	ia := adr - m.start
	return m.data[ia]
}

func (m *memory) SetMem(adr uint16, dt uint8) {
	if !m.readonly {
		ia := adr - m.start
		m.data[ia] = dt
	}
}

func (m *memory) IsMapped(adr uint16) bool {
	return (adr >= m.start) && (adr < m.end)
}
