package nes

import (
	"log"
)

const (
	SCANLINE_TYPE_PRE     = 0x1
	SCANLINE_TYPE_VISIBLE = 0x2
	SCANLINE_TYPE_POST    = 0x3
	SCANLINE_TYPE_VBLANK  = 0x4

	SCANLINE_NMI = 241

	DOT_TYPE_VISIBLE   = 0x1
	DOT_TYPE_PREFETCH  = 0x2
	DOT_TYPE_INVISIBLE = 0x3
)

// From http://www.thealmightyguru.com/Games/Hacking/Wiki/index.php?title=NES_Palette
var Palette = [64]uint32{
	0x7C7C7C, 0x0000FC, 0x0000BC, 0x4428BC, 0x940084, 0xA80020, 0xA81000, 0x881400,
	0x503000, 0x007800, 0x006800, 0x005800, 0x004058, 0x000000, 0x000000, 0x000000,
	0xBCBCBC, 0x0078F8, 0x0058F8, 0x6844FC, 0xD800CC, 0xE40058, 0xF83800, 0xE45C10,
	0xAC7C00, 0x00B800, 0x00A800, 0x00A844, 0x008888, 0x000000, 0x000000, 0x000000,
	0xF8F8F8, 0x3CBCFC, 0x6888FC, 0x9878F8, 0xF878F8, 0xF85898, 0xF87858, 0xFCA044,
	0xF8B800, 0xB8F818, 0x58D854, 0x58F898, 0x00E8D8, 0x787878, 0x000000, 0x000000,
	0xFCFCFC, 0xA4E4FC, 0xB8B8F8, 0xD8B8F8, 0xF8B8F8, 0xF8A4C0, 0xF0D0B0, 0xFCE0A8,
	0xF8D878, 0xD8F878, 0xB8F8B8, 0xB8F8D8, 0x00FCFC, 0xF8D8F8, 0x000000, 0x000000,
}

func scanlineType(scanlineN int) int {
	switch {
	case scanlineN == 261:
		return SCANLINE_TYPE_PRE

	case scanlineN < 240:
		return SCANLINE_TYPE_VISIBLE

	case scanlineN == 240:
		return SCANLINE_TYPE_POST

	case scanlineN >= 241 && scanlineN <= 260:
		return SCANLINE_TYPE_VBLANK

	default:
		log.Fatalf("Invalid scanline number %d\n", scanlineN)
		return 0
	}
}

/*
  Screen resolution: 256 cols * 240 rows pixels
  Scanlines: 262 per frame
  Dots:      341 per scanline

  Timings extracted from http://wiki.nesdev.com/w/images/d/d1/Ntsc_timing.png
*/

func (ppu *PPU) TickScanline() {
	line := ppu.Scanline
	lineType := scanlineType(line)

	// Pre-render scanline
	if lineType == SCANLINE_TYPE_PRE {
		ppu.tickPreScanline()

		// Visible scanline
	} else if lineType == SCANLINE_TYPE_VISIBLE {
		ppu.tickVisibleScanline()

	} else if line == SCANLINE_NMI {
		if ppu.Dot == 1 {
			ppu.setVerticalBlank()
		}
	} else if lineType == SCANLINE_TYPE_POST {
		// Currently setting frame at the same time as vblank is triggered, so no-op for now
	}

	ppu.Dot += 1
	if ppu.Dot == 341 {
		ppu.Scanline += 1
		if ppu.Scanline == 262 {
			// Wrap around
			ppu.Scanline = 0
		}
		ppu.Dot = 0
	}
}

func (ppu *PPU) tickPreScanline() {
	dot := ppu.Dot

	if dot == 1 {
		//Not in VBlank anymore. Prepare for next visible scanlines.
		ppu.STATUS.VBlankStarted = false
		ppu.STATUS.Sprite0Hit = false
		ppu.STATUS.SpriteOverflow = false

	} else if dot >= 280 && dot <= 304 {
		if ppu.MASK.shouldRender() {
			ppu.ADDR.TransferY()
		}
	} else if dot == 257 {
		ppu.spriteCount = 0
	}

	// Now do everything a visible line does
	ppu.tickVisibleScanline()
}

func (ppu *PPU) tickVisibleScanline() {
	dot := ppu.Dot
	isFetchTime := (dot >= 1 && dot <= 256) || (dot >= 321 && dot <= 336)

	if !ppu.MASK.shouldRender() {
		return
	}

	if dot >= 1 && dot <= 256 {
		ppu.RenderSinglePixel()
	}

	// Background evaluation

	if isFetchTime {

		ppu.BgTileShiftLow <<= 1
		ppu.BgTileShiftHigh <<= 1
		ppu.AttrShiftLow <<= 1
		ppu.AttrShiftHigh <<= 1
		ppu.AttrShiftLow |= (ppu.AttrLatchLow << 0)
		ppu.AttrShiftHigh |= (ppu.AttrLatchHigh << 1)

		switch ppu.Dot % 8 {
		case 1:
			ppu.tempTileAddr = ppu.ADDR.NameTableAddr()

			// Feed new data into the background tile shift registers
			ppu.BgTileShiftLow |= uint16(ppu.BgLatchLow)
			ppu.BgTileShiftHigh |= uint16(ppu.BgLatchHigh)

			// Feed new data into the attribute latches
			ppu.AttrLatchLow = (ppu.AttrTableLatch >> 0) & 0x1
			ppu.AttrLatchHigh = (ppu.AttrTableLatch >> 1) & 0x1
		case 2:
			ppu.NameTableLatch = ppu.Read(ppu.tempTileAddr)
		case 3:
			ppu.tempTileAddr = ppu.ADDR.AttrTableAddr()
		case 4:
			shift := ((ppu.ADDR.VAddr >> 4) & 4) | (ppu.ADDR.VAddr & 2)
			ppu.AttrTableLatch = ppu.Read(ppu.tempTileAddr) >> shift
		case 5:
			ppu.tempTileAddr = ppu.LowBGTileAddr()
		case 6:
			ppu.BgLatchLow = ppu.Read(ppu.tempTileAddr)
		case 7:
			ppu.tempTileAddr = ppu.HighBGTileAddr()
		case 0:
			ppu.BgLatchHigh = ppu.Read(ppu.tempTileAddr)
			ppu.ADDR.IncrementCoarseX()
		}
	}

	// Sprite evaluation

	if dot == 257 {
		ppu.EvalSprites()
	}

	if dot == 256 {
		ppu.ADDR.IncrementFineY()
	}

	if dot == 257 {
		ppu.ADDR.TransferX()
	}
}

func (ppu *PPU) RenderSinglePixel() {
	x := ppu.Dot - 2
	y := ppu.Scanline

	backgroundPixel := ppu.GetBgPixel()

	// Any sprite on this dot?
	oamEntry := ppu.GetCurrentOAMEntry()
	spritePixel := ppu.GetSpritePixel(oamEntry)

	colorOffset := ppu.EvalSpritePriority(x, spritePixel, backgroundPixel, oamEntry)

	addr := ppu.Read(0x3f00 + uint16(colorOffset))
	c := Palette[addr]

	r := uint8((c >> 16) & 0xff)
	g := uint8((c >> 8) & 0xff)
	b := uint8((c >> 0) & 0xff)

	if x >= 0 && x <= 256 && y < 240 {
		pos := 4 * (y*256 + x)
		ppu.Pixels[pos+0] = 0xff
		ppu.Pixels[pos+1] = b
		ppu.Pixels[pos+2] = g
		ppu.Pixels[pos+3] = r
	}
}

func (ppu *PPU) EvalSpritePriority(x int, spritePixel uint8, backgroundPixel uint8, oamEntry *OAMEntry) uint8 {
	isTransparent := func(pixel uint8) bool {
		// 2 least significant bits are the actual tile data
		return pixel&0x03 == 0
	}

	isBgTransp := isTransparent(backgroundPixel)
	isSpriteTransp := isTransparent(spritePixel)

	if isBgTransp && isSpriteTransp {
		return 0
	} else if !isBgTransp && isSpriteTransp {
		return backgroundPixel
	} else if isBgTransp && !isSpriteTransp {
		return spritePixel | 0x10

		// Both bg and sprite are opaque
	} else {
		// Collision?
		if oamEntry.id == 0 && x < 255 {
			ppu.STATUS.Sprite0Hit = true
		}

		if oamEntry.priority == 0 {
			return spritePixel | 0x10
		} else {
			return backgroundPixel
		}
	}
}

func unpackOAMEntry(entryIdx int, oamData []uint8) (y, tileN, attrs, x uint8) {
	baseIdx := 4 * entryIdx

	y = oamData[baseIdx+0]
	tileN = oamData[baseIdx+1]
	attrs = oamData[baseIdx+2]
	x = oamData[baseIdx+3]

	return
}

func isVFlipped(spriteAttrs uint8) bool {
	return spriteAttrs>>7 != 0
}

func isHFlipped(spriteAttrs uint8) bool {
	return spriteAttrs>>6 != 0
}

func (ppu *PPU) EvalSprites() {
	nSpritesInScanline := 0

	for oamIdx := 0; oamIdx < len(ppu.oamData)/4; oamIdx++ {
		y, tileN, attrs, x := unpackOAMEntry(oamIdx, ppu.oamData[:])

		if int(y) <= ppu.Scanline && int(y)+int(ppu.CTRL.SpriteSize) > ppu.Scanline {
			if nSpritesInScanline < 8 {

				spriteRow := ppu.Scanline - int(y)

				if isVFlipped(attrs) {
					spriteRow = int(ppu.CTRL.SpriteSize) - spriteRow
				}

				var patternTableBaseAddr uint16

				if ppu.CTRL.SpriteSize == 8 {
					patternTableBaseAddr = ppu.CTRL.SpritePatTableAddr
				} else {
					patternTableBaseAddr = 0x1000 * uint16(tileN&0x01)

					if spriteRow >= 8 {
						spriteRow -= 8
						tileN |= 0x01
					} else {
						tileN &= 0xfe
					}
				}

				addr := patternTableBaseAddr + uint16(tileN)*16 + uint16(spriteRow)

				ppu.oamEntries[nSpritesInScanline].id = oamIdx
				ppu.oamEntries[nSpritesInScanline].x = x
				ppu.oamEntries[nSpritesInScanline].attrs = attrs
				ppu.oamEntries[nSpritesInScanline].tileLow = ppu.Read(addr)
				ppu.oamEntries[nSpritesInScanline].tileHigh = ppu.Read(addr + 8)
				ppu.oamEntries[nSpritesInScanline].priority = (attrs >> 5) & 0x1

				nSpritesInScanline++
			} else {
				ppu.STATUS.SpriteOverflow = true
			}
		}
	}
	ppu.spriteCount = nSpritesInScanline
}

func (ppu *PPU) GetBgPixel() uint8 {
	if ppu.MASK.showBg {
		// Pull bg attribute and tile from the 4 shift registers into a 4-bit word
		fx := ppu.ADDR.FineXScroll
		return uint8(
			uint16(((ppu.AttrShiftHigh>>(7-fx))&0x1)<<3)|
				uint16(((ppu.AttrShiftLow>>(7-fx))&0x1)<<2)|
				(((ppu.BgTileShiftHigh>>(15-fx))&0x1)<<1)|
				(((ppu.BgTileShiftLow>>(15-fx))&0x1)<<0)) & 0x0f
	} else {
		return 0
	}
}

func (ppu *PPU) GetCurrentOAMEntry() *OAMEntry {
	for oamIdx := 0; oamIdx < ppu.spriteCount; oamIdx++ {
		entry := ppu.oamEntries[oamIdx]

		// We know this sprite is visible in this scanline, but is it
		// visible in this column?
		posX := ppu.Dot - 1

		if int(entry.x) <= posX && int(entry.x)+8 >= posX {
			return &entry
		}
	}

	return nil
}

func (ppu *PPU) GetSpritePixel(entry *OAMEntry) uint8 {
	if !ppu.MASK.showSprites || entry == nil {
		return 0
	}
	posX := ppu.Dot - 1

	// Intra-sprite, fine x scroll
	fx := uint8(posX - 1 - int(entry.x))

	if isHFlipped(entry.attrs) {
		fx = 7 - fx
	}

	return ((entry.attrs&0x3)<<2 |
		((entry.tileHigh>>(7-fx))&0x1)<<1 |
		((entry.tileLow>>(7-fx))&0x1)<<0) & 0x0f
}
