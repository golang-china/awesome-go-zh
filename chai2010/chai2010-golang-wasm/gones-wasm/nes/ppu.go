package nes

import (
	"log"
)

func byteSetter(v uint8, bitN uint8, ifNotSet uint8, ifSet uint8) uint8 {
	if v&(0x1<<bitN) == 0 {
		return ifNotSet
	} else {
		return ifSet
	}
}

func addrSetter(v uint8, bitN uint8, ifNotSet uint16, ifSet uint16) uint16 {
	if v&(0x1<<bitN) == 0 {
		return ifNotSet
	} else {
		return ifSet
	}
}

func boolSetter(v uint8, bitN uint8, ifNotSet bool, ifSet bool) bool {
	if v&(0x1<<bitN) == 0 {
		return ifNotSet
	} else {
		return ifSet
	}
}

const (
	MS_READ_EXT  = false
	MS_WRITE_EXT = true
)

type PPUCTRL struct {
	NameTableAddr     uint16
	VRAMReadIncrement uint16
	// Addr for 8x8 sprites only (ignored for 16x16)
	SpritePatTableAddr uint16
	BgTableAddr        uint16
	SpriteSize         uint8
	MasterSlave        bool
	NMIonVBlank        bool
}

func (ctrl *PPUCTRL) Set(v uint8) {
	switch v & 0x3 {
	case 0x0:
		ctrl.NameTableAddr = 0x2000
	case 0x1:
		ctrl.NameTableAddr = 0x2400
	case 0x2:
		ctrl.NameTableAddr = 0x2800
	case 0x3:
		ctrl.NameTableAddr = 0x2c00
	}

	ctrl.VRAMReadIncrement = addrSetter(v, 2, 0x0001, 0x0020)
	ctrl.SpritePatTableAddr = addrSetter(v, 3, 0x0000, 0x1000)
	ctrl.BgTableAddr = addrSetter(v, 4, 0x0000, 0x1000)
	ctrl.SpriteSize = byteSetter(v, 5, 8, 16)
	ctrl.MasterSlave = boolSetter(v, 6, MS_READ_EXT, MS_WRITE_EXT)
	ctrl.NMIonVBlank = boolSetter(v, 7, false, true)
}

type PPUADDR struct {
	// Internal PPU registers v, t, w, x
	VAddr       uint16
	TAddr       uint16
	WriteHi     bool
	FineXScroll uint8
}

func (ppu *PPU) LowBGTileAddr() uint16 {
	return ppu.CTRL.BgTableAddr + uint16(ppu.NameTableLatch)*16 + ppu.ADDR.FineY()
}

func (ppu *PPU) HighBGTileAddr() uint16 {
	return ppu.LowBGTileAddr() + 8
}

func (addr *PPUADDR) NameTableAddr() uint16 {
	return 0x2000 | (addr.VAddr & 0x0fff)
}

// http://wiki.nesdev.com/w/index.php/PPU_scrolling
func (addr *PPUADDR) AttrTableAddr() uint16 {
	v := addr.VAddr
	return 0x23c0 | (v & 0x0c00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)
}

func (addr *PPUADDR) FineY() uint16 {
	return (addr.VAddr >> 12) & 0x07
}

func (addr *PPUADDR) Write(v uint8) {
	if addr.WriteHi == false {
		addr.TAddr = (addr.TAddr & 0x80FF) | (uint16(v)&0x3F)<<8
		addr.WriteHi = true
	} else {
		addr.TAddr = (addr.TAddr & 0xFF00) | uint16(v)
		addr.VAddr = addr.TAddr
		addr.WriteHi = false
	}
}

// http://wiki.nesdev.com/w/index.php/PPU_scrolling
func (addr *PPUADDR) TransferX() {
	addr.VAddr = (addr.VAddr & 0xFBE0) | (addr.TAddr & 0x041F)
}

func (addr *PPUADDR) TransferY() {
	addr.VAddr = (addr.VAddr & 0x841F) | (addr.TAddr & 0x7BE0)
}

func (addr *PPUADDR) SetOnCTRLWrite(v uint8) {
	addr.TAddr = (addr.TAddr & 0xf3ff) | uint16(v&0x03)<<10
}

func (addr *PPUADDR) SetOnSTATUSRead() {
	addr.WriteHi = false
}

// http://wiki.nesdev.com/w/index.php/PPU_scrolling
func (addr *PPUADDR) SetOnSCROLLWrite(v uint8) {
	if addr.WriteHi == false {
		addr.TAddr |= uint16(v >> 3)
		addr.FineXScroll = v & 0x7
		addr.WriteHi = true
	} else {
		addr.TAddr |= (uint16(v)&0x07)<<12 | (uint16(v)&0xf8)<<2
		addr.WriteHi = false
	}
}

// http://wiki.nesdev.com/w/index.php/PPU_scrolling#Y_increment
func (addr *PPUADDR) IncrementFineY() {
	v := addr.VAddr
	var y uint16

	if (v & 0x7000) != 0x7000 {
		v += 0x1000
		addr.VAddr = v
	} else {
		v &= 0x8fff
		y = (v & 0x03e0) >> 5
		if y == 29 {
			y = 0
			v ^= 0x0800
		} else {
			if y == 31 {
				y = 0
			} else {
				y += 1
			}
		}
		addr.VAddr = (v & 0xfc1f) | (y << 5)
	}
}

// http://wiki.nesdev.com/w/index.php/PPU_scrolling#X_increment
func (addr *PPUADDR) IncrementCoarseX() {
	v := addr.VAddr

	if (v & 0x001F) == 31 {
		v &= 0xFFE0
		v ^= 0x0400
	} else {
		v += 1
	}

	addr.VAddr = v
}

type PPUMASK struct {
	greyscale       bool
	showBgLeft      bool
	showSpritesLeft bool
	showBg          bool
	showSprites     bool
	emphasisRed     bool
	emphasisGreen   bool
	emphasisBlue    bool
}

func (mask *PPUMASK) Set(v uint8) {
	mask.greyscale = boolSetter(v, 0, false, true)
	mask.showBgLeft = boolSetter(v, 1, false, true)
	mask.showSpritesLeft = boolSetter(v, 2, false, true)
	mask.showBg = boolSetter(v, 3, false, true)
	mask.showSprites = boolSetter(v, 4, false, true)
	mask.emphasisRed = boolSetter(v, 5, false, true)
	mask.emphasisGreen = boolSetter(v, 6, false, true)
	mask.emphasisBlue = boolSetter(v, 7, false, true)
}

func (mask *PPUMASK) shouldRender() bool {
	return mask.showBg || mask.showSprites
}

type PPUSTATUS struct {
	SpriteOverflow bool
	Sprite0Hit     bool
	VBlankStarted  bool

	// So we can simulate a dirty bus when reading CTRL
	LastWrite uint8
}

func (status *PPUSTATUS) Get() (result uint8) {
	if status.SpriteOverflow {
		result |= (0x1 << 5)
	}
	if status.Sprite0Hit {
		result |= (0x1 << 6)
	}
	if status.VBlankStarted {
		result |= (0x1 << 7)
	}

	result |= (status.LastWrite & 0x1f)

	status.VBlankStarted = false

	return
}

type OAMEntry struct {
	id       int
	x        uint8
	attrs    uint8
	priority uint8
	tileLow  uint8
	tileHigh uint8

	patt uint32
}

type PPU struct {
	CTRL   *PPUCTRL
	ADDR   *PPUADDR
	MASK   *PPUMASK
	STATUS *PPUSTATUS

	CPU    *CPU
	rom    *Rom
	TV     *TV
	Mapper Mapper

	Pixels []byte

	ReadDataBuffer uint8

	Scanline int
	Dot      int

	NameTableLatch uint8
	AttrTableLatch uint8

	// Background latches & shift registers
	BgLatchLow      uint8
	BgLatchHigh     uint8
	BgTileShiftLow  uint16
	BgTileShiftHigh uint16

	// Attribute latches & shift registers
	AttrLatchLow  uint8
	AttrLatchHigh uint8
	AttrShiftLow  uint8
	AttrShiftHigh uint8

	tempTileAddr uint16

	// Holds 8 OAM entries for a single scanline
	spriteCount int
	oamEntries  [8]OAMEntry

	Frame uint64 // frame counter

	// storage variables
	paletteData   [32]byte
	nameTableData [2048]byte
	oamData       [256]byte

	oamAddress byte
}

func MakePPU(cpu *CPU, rom *Rom, mapper Mapper) *PPU {
	ppu := PPU{
		CPU:    cpu,
		rom:    rom,
		Mapper: mapper,
		CTRL:   &PPUCTRL{},
		ADDR:   &PPUADDR{},
		MASK:   &PPUMASK{},
		STATUS: &PPUSTATUS{},
	}
	ppu.Pixels = make([]byte, 4*256*240)
	ppu.Reset()
	return &ppu
}

func (ppu *PPU) Reset() {
	ppu.Dot = 340
	ppu.Scanline = 0
	ppu.Frame = 0
	ppu.CTRL.Set(0)
	ppu.ADDR.SetOnSCROLLWrite(0)
	ppu.MASK.Set(0)
	ppu.writeOAMAddress(0)
}

func (ppu *PPU) readPalette(address uint16) byte {
	if address >= 16 && address%4 == 0 {
		address -= 16
	}
	return ppu.paletteData[address]
}

func (ppu *PPU) writePalette(address uint16, value byte) {
	if address >= 16 && address%4 == 0 {
		address -= 16
	}
	ppu.paletteData[address] = value
}

// VRAM 0x0000 - 0x3eff reads are buffered!
// https://wiki.nesdev.com/w/index.php/PPU_registers#The_PPUDATA_read_buffer_.28post-fetch.29
func (ppu *PPU) ReadData() uint8 {
	current := ppu.Read(ppu.ADDR.VAddr)
	ppu.ADDR.VAddr += ppu.CTRL.VRAMReadIncrement

	if ppu.ADDR.VAddr <= 0x3eff {
		ppu.ReadDataBuffer, current = current, ppu.ReadDataBuffer
	}

	// TODO: fogleman subtracts 0x1000 form this addr. Why?
	return current
}

func (ppu *PPU) WriteData(v uint8) {
	ppu.Write(ppu.ADDR.VAddr, v)
	ppu.ADDR.VAddr += ppu.CTRL.VRAMReadIncrement
}

func (ppu *PPU) writeOAMAddress(value byte) {
	ppu.oamAddress = value
}

func (ppu *PPU) readOAMData() byte {
	return ppu.oamData[ppu.oamAddress]
}

func (ppu *PPU) writeOAMData(value byte) {
	ppu.oamData[ppu.oamAddress] = value
	ppu.oamAddress++
}

func (ppu *PPU) writeDMA(value byte) {
	address := uint16(value) << 8
	for i := 0; i < 256; i++ {
		ppu.oamData[ppu.oamAddress] = ppu.CPU.mem.Read8(address)
		ppu.oamAddress++
		address++
	}
}

func (ppu *PPU) setVerticalBlank() {

	ppu.TV.SetFrame(ppu.Pixels)

	ppu.STATUS.VBlankStarted = true
	if ppu.CTRL.NMIonVBlank {
		ppu.CPU.nmiRequested = true
	}
}

func (ppu *PPU) getMirroredNametableAddr(addr uint16) uint16 {
	if ppu.rom.Header.VerticalMirror {
		return addr % 0x0800
	} else {
		return (addr>>1)&0x400 + addr%0x400
	}
}

func (ppu *PPU) Read(address uint16) byte {
	address = address % 0x4000
	switch {
	case address < 0x2000:
		return ppu.Mapper.Read8(address)
	case address < 0x3F00:
		return ppu.nameTableData[ppu.getMirroredNametableAddr(address)]
	case address < 0x4000:
		return ppu.readPalette(address % 32)
	default:
		log.Fatalf("Invalid read from ppu at address %x", address)
	}
	return 0
}

func (ppu *PPU) Write(address uint16, value byte) {
	address = address % 0x4000
	switch {
	case address < 0x2000:
		ppu.rom.CHRROM[address] = value
	case address < 0x3F00:
		ppu.nameTableData[ppu.getMirroredNametableAddr(address)] = value
	case address < 0x4000:
		ppu.writePalette(address%32, value)
	default:
		log.Fatalf("Invalid write to ppu at address %x", address)
	}
}
