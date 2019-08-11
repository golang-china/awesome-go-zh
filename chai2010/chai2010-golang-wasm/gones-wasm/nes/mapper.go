package nes

import "log"

type Mapper interface {
	Read8(addr uint16) uint8
	Write8(addr uint16, v uint8)
}

func MakeMapper(rom *Rom) Mapper {
	switch rom.Header.MapperN {
	case 0:
		return &Mapper0{ROM: rom}
	case 3:
		return &Mapper3{ROM: rom}
	case 4:
		return MakeMapper4(rom)
	default:
		log.Fatalf("Unsupported mapper: %x", rom.Header.MapperN)
		return nil
	}
}
