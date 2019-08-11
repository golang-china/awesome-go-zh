package nes

import "log"

type Mapper0 struct {
	ROM *Rom
}

func (m *Mapper0) Read8(addr uint16) uint8 {
	switch {
	// PPU pattern tables - to be read by the PPU
	case addr >= 0x0000 && addr < 0x2000:
		return m.ROM.CHRROM[addr]

	// CHRROM. Rom can either have 1 or 2 PRGROM banks. If it has only one, we need to
	// mirror access to the first one.
	case addr >= 0x8000:
		return m.ROM.PRGROM.Read8((addr - 0x8000) % (0x4000 * uint16(m.ROM.Header.NPRGROMBanks)))
	default:
		log.Fatalf("Invalid read from mapper at %x", addr)
		return 0
	}
}

func (m *Mapper0) Write8(addr uint16, v uint8) {
	switch {
	case addr >= 0x8000:
		// No-op
	default:
		log.Fatalf("Invalid write to mapper at %x", addr)
	}
}
