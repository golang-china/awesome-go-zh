package nes

import "log"

const (
	// 1kB
	MAPPER4_CHRROM_BANK_SIZE = 0x0400
	// 8kB
	MAPPER4_PRGROM_BANK_SIZE = 0x2000
)

type Mapper4 struct {
	ROM *Rom

	// Written at 0x8000 - 0xa0000 (on even addrs)
	Register uint8

	// R0 - R7 internal registers
	BankRegisters [8]uint8

	// 8 1kB regions
	CHRMappings [8]uint8

	// 4 8kB regions
	PRGMappings [4]uint8
}

func MakeMapper4(rom *Rom) *Mapper4 {
	m := &Mapper4{ROM: rom}

	n := uint8(len(m.ROM.PRGROM) / MAPPER4_PRGROM_BANK_SIZE)

	m.PRGMappings[0] = 0
	m.PRGMappings[1] = 1
	m.PRGMappings[2] = n - 2
	m.PRGMappings[3] = n - 1

	return m
}

func (m *Mapper4) Read8(addr uint16) uint8 {
	switch {
	// PPU CHRROM access (pattern tables)
	case addr >= 0x0000 && addr < 0x2000:
		bank := addr / MAPPER4_CHRROM_BANK_SIZE
		rest := addr % MAPPER4_CHRROM_BANK_SIZE
		mappedBank := m.CHRMappings[bank]
		return m.ROM.CHRROM[uint64(mappedBank)*MAPPER4_CHRROM_BANK_SIZE+uint64(rest)]

		// PRGROM access
	case addr >= 0x8000:
		bank := (addr - 0x8000) / MAPPER4_PRGROM_BANK_SIZE
		rest := (addr - 0x8000) % MAPPER4_PRGROM_BANK_SIZE
		mappedBank := m.PRGMappings[bank]
		return m.ROM.PRGROM[uint64(mappedBank)*MAPPER4_PRGROM_BANK_SIZE+uint64(rest)]

	default:
		log.Fatalf("Invalid read from mapper at %x", addr)
	}
	return 0
}

func (m *Mapper4) Write8(addr uint16, v uint8) {
	switch {
	case addr >= 0x8000 && addr < 0xa000:
		if addr&0x1 == 0 {
			m.Register = v
		} else {
			m.BankRegisters[m.Register&0x07] = v
		}

		m.doBankSwitch()

	case addr >= 0xa000 && addr < 0xbfff && addr&0x1 == 0x0:
		m.ROM.Header.VerticalMirror = v&0x01 == 0x00

	default:
		//log.Fatalf("Invalid write to mapper at %x", addr)
	}
}

// These black magic mappings are documented on http://wiki.nesdev.com/w/index.php/MMC3
func (m *Mapper4) doBankSwitch() {
	chrMode := (m.Register >> 7) & 0x01
	prgMode := (m.Register >> 6) & 0x01

	if chrMode == 0x00 {
		m.CHRMappings[0] = m.BankRegisters[0] & 0xfe
		m.CHRMappings[1] = m.BankRegisters[0] | 0x01
		m.CHRMappings[2] = m.BankRegisters[1] & 0xfe
		m.CHRMappings[3] = m.BankRegisters[1] | 0x01
		m.CHRMappings[4] = m.BankRegisters[2]
		m.CHRMappings[5] = m.BankRegisters[3]
		m.CHRMappings[6] = m.BankRegisters[4]
		m.CHRMappings[7] = m.BankRegisters[5]

	} else {
		m.CHRMappings[0] = m.BankRegisters[2]
		m.CHRMappings[1] = m.BankRegisters[3]
		m.CHRMappings[2] = m.BankRegisters[4]
		m.CHRMappings[3] = m.BankRegisters[5]
		m.CHRMappings[4] = m.BankRegisters[0] & 0xfe
		m.CHRMappings[5] = m.BankRegisters[0] | 0x01
		m.CHRMappings[6] = m.BankRegisters[1] & 0xfe
		m.CHRMappings[7] = m.BankRegisters[1] | 0x01
	}

	n := uint8(len(m.ROM.PRGROM) / MAPPER4_PRGROM_BANK_SIZE)

	if prgMode == 0x00 {
		m.PRGMappings[0] = m.BankRegisters[6]
		m.PRGMappings[1] = m.BankRegisters[7]
		m.PRGMappings[2] = n - 2
		m.PRGMappings[3] = n - 1
	} else {
		m.PRGMappings[0] = n - 2
		m.PRGMappings[1] = m.BankRegisters[7]
		m.PRGMappings[2] = m.BankRegisters[6]
		m.PRGMappings[3] = n - 1
	}
}
