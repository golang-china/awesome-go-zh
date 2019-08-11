package main

import (
	"os"
	"time"

	"github.com/chai2010/gones-wasm/nes"
	"github.com/chai2010/gones-wasm/static"
)

const (
	FRAME_RATE   = 60
	FRAME_CYCLES = nes.CPU_FREQ / FRAME_RATE

	// All timings in nanoseconds
	FRAME_TIME = 1e9 / float64(FRAME_RATE)
)

func main() {
	tv := nes.MakeTV()

	rompath := "BattleCity.nes"
	if len(os.Args) == 2 {
		rompath = os.Args[1]
	}

	rom := nes.ReadROMData(rompath, []byte(static.Files[rompath]))
	mapper := nes.MakeMapper(rom)

	ppu := nes.MakePPU(nil, rom, mapper)
	ppu.Reset()

	controller := nes.MakeController()
	cpuAddrSpace := nes.MakeCPUAddrSpace(rom, ppu, controller, mapper)
	cpu := nes.MakeCPU(cpuAddrSpace)

	ppu.CPU = cpu
	ppu.TV = tv

	var cpuCycles, newCycles, ppuCycles int

	t0 := time.Now().UnixNano()

	for {

		newCycles = cpu.Run()
		cpuCycles += newCycles

		for ppuCycles = 0; ppuCycles < 3*newCycles; ppuCycles++ {
			ppu.TickScanline()
		}

		if cpuCycles > FRAME_CYCLES {
			cpuCycles -= FRAME_CYCLES

			tv.ShowPixels()
			tv.UpdateInputState(controller)

			if delta := FRAME_TIME - float64(time.Now().UnixNano()-t0); delta > 0 {
				time.Sleep(time.Duration(delta))
			}

			t0 = time.Now().UnixNano()
		}
	}
}
