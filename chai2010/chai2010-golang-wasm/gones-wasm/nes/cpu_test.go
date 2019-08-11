package nes

import (
	"testing"
)

type testCPUAddrSpace struct {
	RAM Memory
}

func (as *testCPUAddrSpace) Read8(addr uint16) uint8 {
	return as.RAM.Read8(addr)
}

func (as *testCPUAddrSpace) Write8(addr uint16, v uint8) {
	as.RAM.Write8(addr, v)
}

func (as *testCPUAddrSpace) Read16(addr uint16) uint16 {
	lo := uint16(as.Read8(addr))
	hi := uint16(as.Read8(addr + 1))
	return (hi << 8) + lo
}

func (as *testCPUAddrSpace) Write16(addr uint16, v uint16) {
	as.Write8(addr, uint8(v&0xff))
	as.Write8(addr+1, uint8(v>>8))
}

func makeTestCPU() *CPU {
	return MakeCPU(
		&testCPUAddrSpace{
			RAM: make(Memory, 0x10000),
		},
	)
}

func TestPushPop(t *testing.T) {
	cpu := makeTestCPU()

	cpu.Push8(0xde)
	cpu.Push8(0xad)
	cpu.Push16(0xbeaf)

	if v := cpu.Pop16(); v != 0xbeaf {
		t.Fatalf("Wrong value for Pop16: %x", v)
	}

	if v := cpu.Pop8(); v != 0xad {
		t.Fatalf("Wrong value for Pop8: %x", v)
	}

	if v := cpu.Pop8(); v != 0xde {
		t.Fatalf("Wrong value for Pop8: %x", v)
	}

}
