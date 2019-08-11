package nes

import "log"

type Mapper3 struct {
	ROM *Rom

	selectedCHRROMBank uint8
}

func (m *Mapper3) Read8(addr uint16) uint8 {
	switch {
	// PPU pattern tables - to be read by the PPU
	// We'll read from the previously selected CHRROM bank
	case addr >= 0x0000 && addr < 0x2000:
		return m.ROM.CHRROM[0x2000*uint16(m.selectedCHRROMBank)+addr]

	// CHRROM. Rom can either have 1 or 2 PRGROM banks. If it has only one, we need to
	// mirror access to the first one.
	case addr >= 0x8000:
		return m.ROM.PRGROM.Read8((addr - 0x8000) % (0x4000 * uint16(m.ROM.Header.NPRGROMBanks)))
	default:
		log.Fatalf("Invalid read from mapper at %x", addr)
		return 0
	}
}

func (m *Mapper3) Write8(addr uint16, v uint8) {
	switch {
	case addr >= 0x8000:
		// Select which CHRROM bank is in use. Only supports 4 banks => 32 KiB, so we only
		// look to the two least significant bits of v here
		m.selectedCHRROMBank = v & 0x3

	default:
		log.Fatalf("Invalid write to mapper at %x", addr)
	}
}
