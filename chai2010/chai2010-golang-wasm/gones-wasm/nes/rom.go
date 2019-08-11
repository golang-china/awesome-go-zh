package nes

import (
	"io/ioutil"
	"log"
)

//http://fms.komkon.org/EMUL8/NES.html
type Rom struct {
	Header *RomHeader
	PRGROM Memory
	CHRROM Memory
	PRGRAM Memory
}

type RomHeader struct {
	MapperN uint8

	// 16kB each
	NPRGROMBanks uint8

	// 8kB each
	NCHRROMBanks uint8

	HasTrainer bool

	VerticalMirror bool
}

func ReadROMData(path string, data []byte) *Rom {
	if len(data) == 0 {
		data, _ = ioutil.ReadFile(path)
	}

	if string(data[:3]) != "NES" {
		log.Fatalf("Invalid ROM file" + string(data[:3]))
	}

	header := &RomHeader{
		MapperN:        (data[6] >> 4) | (data[7] & 0xf0),
		NPRGROMBanks:   data[4],
		NCHRROMBanks:   data[5],
		HasTrainer:     (data[6] & (0x1 << 2)) > 0,
		VerticalMirror: data[6]&0x1 == 0x1,
	}

	var (
		prgBeginning uint64 = 16
		prgEnd       uint64 = 16 + uint64(header.NPRGROMBanks)*0x4000
	)

	if header.HasTrainer {
		prgBeginning += 512
		prgEnd += 512
	}

	var (
		chrBeginning uint64 = prgEnd
		chrEnd       uint64 = prgEnd + uint64(header.NCHRROMBanks)*0x2000
	)

	rom := &Rom{
		Header: header,
		PRGROM: data[prgBeginning:prgEnd],
		CHRROM: data[chrBeginning:chrEnd],
		PRGRAM: make(Memory, 0x2000),
	}

	return rom
}

func ReadROM(path string) *Rom {
	data, _ := ioutil.ReadFile(path)

	if string(data[:3]) != "NES" {
		log.Fatalf("Invalid ROM file" + string(data[:3]))
	}

	header := &RomHeader{
		MapperN:        (data[6] >> 4) | (data[7] & 0xf0),
		NPRGROMBanks:   data[4],
		NCHRROMBanks:   data[5],
		HasTrainer:     (data[6] & (0x1 << 2)) > 0,
		VerticalMirror: data[6]&0x1 == 0x1,
	}

	var (
		prgBeginning uint64 = 16
		prgEnd       uint64 = 16 + uint64(header.NPRGROMBanks)*0x4000
	)

	if header.HasTrainer {
		prgBeginning += 512
		prgEnd += 512
	}

	var (
		chrBeginning uint64 = prgEnd
		chrEnd       uint64 = prgEnd + uint64(header.NCHRROMBanks)*0x2000
	)

	rom := &Rom{
		Header: header,
		PRGROM: data[prgBeginning:prgEnd],
		CHRROM: data[chrBeginning:chrEnd],
		PRGRAM: make(Memory, 0x2000),
	}

	return rom
}
