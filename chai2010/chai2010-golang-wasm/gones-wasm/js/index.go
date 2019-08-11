package main

import (
	"fmt"
	"runtime"
	"syscall/js"
	"time"

	"github.com/chai2010/gones-wasm/nes"
	"github.com/chai2010/gones-wasm/static"
)

const (
	FRAME_RATE   = 60
	FRAME_CYCLES = nes.CPU_FREQ / FRAME_RATE // 29829

	// All timings in nanoseconds
	FRAME_TIME = 1e9 / float64(FRAME_RATE) // 16.666666ms
)

func init() {
	fmt.Println("FRAME_RATE:", FRAME_RATE)
	fmt.Println("FRAME_CYCLES:", FRAME_CYCLES)

	FRAME_TIME := FRAME_TIME
	fmt.Printf("FRAME_TIME: %v\n", time.Duration(int64(FRAME_TIME))*time.Nanosecond)
}

// 2个goroutine分别定期发出cpu和tv的刷新事件
// 1个goroutine阻塞读取事件，并执行

var (
	tv         *nes.TV
	ppu        *nes.PPU
	cpu        *nes.CPU
	controller *nes.Controller

	cpuCycles int
	newCycles int
	ppuCycles int
	t0        int64
)

func SetupNes() {
	tv = nes.MakeTV()

	rompath := "SuperMarioBros.nes"
	rom := nes.ReadROMData(rompath, []byte(static.Files[rompath]))
	mapper := nes.MakeMapper(rom)

	ppu = nes.MakePPU(nil, rom, mapper)
	ppu.Reset()

	controller = nes.MakeController()
	cpuAddrSpace := nes.MakeCPUAddrSpace(rom, ppu, controller, mapper)
	cpu = nes.MakeCPU(cpuAddrSpace)

	ppu.CPU = cpu
	ppu.TV = tv

	t0 = time.Now().UnixNano()
}

// shift 选择
// enter开始
// 小键盘1=A
// 小键盘2=B
// 方向键是游戏用的

func main() {
	SetupNes()

	if runtime.GOARCH == "wasm" && runtime.GOOS == "js" {
		var cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			StepNes()
			return nil
		})
		js.Global().Set("goSetupNes", cb)

		var keydowncb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			KeyDown(fmt.Sprint(args[0]))
			return nil
		})
		js.Global().Set("goKeyDown", keydowncb)
		var keyupcb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			KeyUp(fmt.Sprint(args[0]))
			return nil
		})
		js.Global().Set("goKeyUp", keyupcb)

		frameLoop := js.Global().Get("doFrameLoop") //doFrameLoop
		frameLoop.Invoke()

		time.Sleep(time.Second)

		select {}
		<-make(chan bool)
		return
	}
}

func KeyDown(key string) {
	tv.KeyDown(controller, key)
}

func KeyUp(key string) {
	tv.KeyUp(controller, key)
}

func StepNes() {
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

			break
		}
	}
}
