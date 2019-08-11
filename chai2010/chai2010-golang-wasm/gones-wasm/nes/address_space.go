package nes

import (
	"log"
)

type AddrSpace interface {
	Read8(addr uint16) uint8
	Write8(addr uint16, v uint8)

	Read16(addr uint16) uint16
	Write16(addr uint16, v uint16)
}

type CPUAddrSpace struct {
	ctrlr *Controller

	RAM Memory
	ROM *Rom

	PPU *PPU

	Mapper Mapper
}

func MakeCPUAddrSpace(rom *Rom, ppu *PPU, ctrlr *Controller, mapper Mapper) *CPUAddrSpace {
	return &CPUAddrSpace{
		ctrlr:  ctrlr,
		RAM:    make(Memory, 0x800),
		ROM:    rom,
		PPU:    ppu,
		Mapper: mapper,
	}
}

//http://wiki.nesdev.com/w/index.php/CPU_memory_map
//https://wiki.nesdev.com/w/index.php/NROM (Hard coded mapper 0 for now)
func (as *CPUAddrSpace) Read8(addr uint16) uint8 {
	switch {
	case addr >= 0 && addr < 0x2000:
		// 0x0800 - 0x1fff mirrors 0x0000 - 0x07ff three times
		return as.RAM.Read8(addr % 0x800)

	// PPU registers
	case addr >= 0x2000 && addr < 0x4000:
		switch 0x2000 + addr%8 {
		case 0x2002:
			as.PPU.ADDR.SetOnSTATUSRead()
			return as.PPU.STATUS.Get()
		case 0x2004:
			return as.PPU.readOAMData()
		case 0x2007:
			return as.PPU.ReadData()
		default:
			return 0
		}

	case addr == 0x4015:
		//log.Printf("Not yet handled read to APU at %x", addr)
		return 0

	case addr == 0x4016:
		//log.Printf("Not yet handled read to controller #1 at %x", addr)
		return as.ctrlr.ReadState()

	case addr == 0x4017:
		//log.Printf("Not yet handled read to controller #2 at %x", addr)
		return 0

	case addr >= 0x4000 && addr < 0x6000:
		//log.Printf("Not yet handled read to %x", addr)
		return 0

	// PRGRAM
	case addr >= 0x6000 && addr < 0x8000:
		return as.ROM.PRGRAM.Read8((addr - 0x6000))

	// ROM PRG banks
	case addr >= 0x8000:
		return as.Mapper.Read8(addr)

	default:
		log.Fatalf("Invalid read from CPU mem space at %x", addr)
		return 0
	}
}

func (as *CPUAddrSpace) Write8(addr uint16, v uint8) {

	switch {
	case addr >= 0 && addr < 0x2000:
		as.RAM.Write8(addr%0x800, v)

	// PPU registers
	case addr >= 0x2000 && addr < 0x4000:
		as.PPU.STATUS.LastWrite = v

		switch 0x2000 + addr%8 {
		case 0x2000:
			as.PPU.CTRL.Set(v)
			as.PPU.ADDR.SetOnCTRLWrite(v)
		case 0x2001:
			as.PPU.MASK.Set(v)
		case 0x2003:
			as.PPU.writeOAMAddress(v)
		case 0x2004:
			as.PPU.writeOAMData(v)
		case 0x2005:
			as.PPU.ADDR.SetOnSCROLLWrite(v)
		case 0x2006:
			as.PPU.ADDR.Write(v)
		case 0x2007:
			as.PPU.WriteData(v)
		}

	case addr >= 0x4000 && addr <= 0x4013:
		//log.Printf("Not yet handled write to APU at %x", addr)

	case addr == 0x4014:
		as.PPU.STATUS.LastWrite = v
		as.PPU.writeDMA(v)

	case addr == 0x4015:
		//log.Printf("Not yet handled write to APU at %x", addr)

	case addr == 0x4016:
		//log.Printf("Not yet handled write to controllers at %x", addr)

	case addr == 0x4017:
	//log.Printf("Not yet handled write to APU at %x", addr)

	case addr >= 0x4000 && addr < 0x6000:
		//log.Printf("Not yet handled write to CPU at %x", addr)

	// PRGRAM
	case addr >= 0x6000 && addr < 0x8000:
		as.ROM.PRGRAM.Write8((addr - 0x6000), v)

	case addr >= 0x8000:
		as.Mapper.Write8(addr, v)

	default:
		log.Printf("Invalid write to CPU mem space at %x", addr)
	}
}

// Little-endian mem layout
func (as *CPUAddrSpace) Read16(addr uint16) uint16 {
	lo := uint16(as.Read8(addr))
	hi := uint16(as.Read8(addr + 1))
	return (hi << 8) + lo
}

func (as *CPUAddrSpace) Write16(addr uint16, v uint16) {
	as.Write8(addr, uint8(v&0xff))
	as.Write8(addr+1, uint8(v>>8))
}
