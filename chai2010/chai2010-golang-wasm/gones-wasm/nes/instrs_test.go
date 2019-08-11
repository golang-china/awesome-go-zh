package nes

import (
	"log"
	"math"
	"testing"
)

func TestBRK(t *testing.T) {
	cpu := makeTestCPU()
	cpu.mem.Write8(0xfffe, 0xad)
	cpu.mem.Write8(0xffff, 0xde)

	brk(cpu, AddrModeImplied)

	if cpu.regs.P&(0x1<<StatusFlagB) == 0 {
		t.Fatalf("Wrong value 0 for status bit B")
	}
	if cpu.regs.PC != 0xdead {
		t.Fatalf("Wrong value for PC register")
	}
}

func TestORA(t *testing.T) {
	cpu := makeTestCPU()
	cpu.regs.PC = 0x0004
	cpu.mem.Write8(cpu.regs.PC+1, 0xad)
	cpu.regs.A = 0x4a
	cpu.regs.X = 0x02

	// AddrModeXIndirect will read from (mem[PC + 1] + X) & 0xff
	cpu.mem.Write8(0xaf, 0x8d)
	cpu.mem.Write8(0x8d, 0xc9)

	ora(cpu, AddrModeXIndirect)

	if cpu.regs.A != uint8(0x4a|0xc9) {
		t.Fatalf("Wrong value for reg A: %x", cpu.regs.A)
	}
	if cpu.regs.P&(0x1<<StatusFlagN) == 0 {
		t.Fatalf("Flag N shouldve been set")
	}
	if cpu.regs.P&(0x1<<StatusFlagZ) != 0 {
		t.Fatalf("Flag Z should not have been set")
	}
}

func TestEOR(t *testing.T) {
	cpu := makeTestCPU()
	cpu.regs.PC = 0x0004
	cpu.mem.Write8(cpu.regs.PC+1, 0xad)
	cpu.regs.A = 0x4a
	cpu.regs.X = 0x02

	// AddrModeXIndirect will read from (mem[PC + 1] + X) & 0xff
	cpu.mem.Write8(0xaf, 0x8d)
	cpu.mem.Write8(0x8d, 0xc9)

	eor(cpu, AddrModeXIndirect)

	if cpu.regs.A != uint8(0x4a^0xc9) {
		t.Fatalf("Wrong value for reg A: %x", cpu.regs.A)
	}
	if cpu.regs.P&(0x1<<StatusFlagN) == 0 {
		t.Fatalf("Flag N shouldve been set")
	}
	if cpu.regs.P&(0x1<<StatusFlagZ) != 0 {
		t.Fatalf("Flag Z should not have been set")
	}
}

func TestASL(t *testing.T) {
	cpu := makeTestCPU()
	cpu.regs.A = 0x8a

	asl(cpu, AddrModeAccumulator)

	if cpu.regs.A != uint8((0x8a<<1)&0xff) {
		t.Fatalf("Wrong value for reg A: %x", cpu.regs.A)
	}
	if cpu.regs.P&(0x1<<StatusFlagC) == 0 {
		t.Fatalf("Flag C shouldve been set")
	}
	if cpu.regs.P&(0x1<<StatusFlagZ) != 0 {
		t.Fatalf("Flag Z should not have been set")
	}
}

func TestBPL(t *testing.T) {
	cpu := makeTestCPU()
	cpu.regs.PC = 0x0004

	// 0xfa = signed -6
	cpu.mem.Write8(cpu.regs.PC+1, 0xfa)

	bpl(cpu, AddrModeRelative)

	// PC should be at 0x0000 (relative jump of -6 plus 2 for the
	// BPL instruct itself)
	if cpu.regs.PC != 0x0000 {
		t.Fatalf("Wrong value for reg PC: %x", cpu.regs.PC)
	}

	cpu.regs.PC = 0x0004

	// 0xfa = signed +6
	cpu.mem.Write8(cpu.regs.PC+1, 0x06)

	bpl(cpu, AddrModeRelative)

	// PC should be at 0x000c (relative jump of +6 plus 2 for the
	// BPL instruct itself)
	if cpu.regs.PC != 0x000c {
		t.Fatalf("Wrong value for reg PC: %x", cpu.regs.PC)
	}
}

func TestJSR(t *testing.T) {
	cpu := makeTestCPU()
	cpu.regs.PC = 0x0004
	cpu.mem.Write16(0x0005, 0xfffb)

	jsr(cpu, AddrModeAbs)

	if cpu.regs.PC != 0xfffb {
		t.Fatalf("Unexpected jump to %x", cpu.regs.PC)
	}

	if p := cpu.Pop16(); p != 0x0003 {
		t.Fatalf("Unexpected pop of %x", p)
	}
}

func TestAND(t *testing.T) {
	cpu := makeTestCPU()

	// AddrModeXIndirect
	cpu.regs.PC = 0x0004
	cpu.regs.A = 0xff
	cpu.regs.X = 0xfb

	// Operand
	cpu.mem.Write8(0x0005, 0x80)

	// Actual address
	cpu.mem.Write8((0x80+0xfb)&0xff, 0xcc)

	// Valur at address
	cpu.mem.Write8(0xcc, 0xab)

	and(cpu, AddrModeXIndirect)

	if cpu.regs.A != 0xab {
		t.Fatalf("Unexpected value of A %x", cpu.regs.A)
	}

	// AddrModeZeroX
	cpu.regs.PC = 0x0004
	cpu.regs.A = 0xf4
	cpu.regs.X = 0xfb
	cpu.mem.Write8(0x0005, 0x80)
	cpu.mem.Write8((0x80+0xfb)&0xff, 0xae)

	and(cpu, AddrModeZeroX)

	if cpu.regs.A != (0xae & 0xf4) {
		t.Fatalf("Unexpected value of A %x", cpu.regs.A)
	}
}

func TestBIT(t *testing.T) {
	cpu := makeTestCPU()
	cpu.regs.PC = 0x0004
	cpu.regs.A = 0xff
	cpu.mem.Write8(0x0005, 0xab)
	cpu.mem.Write8(0x00ab, 0xf0)

	bit(cpu, AddrModeZeroPage)

	if cpu.regs.A != 0xff {
		t.Fatalf("Unexpected value of A %x", cpu.regs.A)
	}

	if cpu.getFlag(StatusFlagZ) {
		t.Fatalf("And result should not be zero")
	}

	if !cpu.getFlag(StatusFlagV) {
		t.Fatalf("Flag V shouldve been set")
	}

	if !cpu.getFlag(StatusFlagN) {
		t.Fatalf("Flag N shouldve been set")
	}
}

func TestROL(t *testing.T) {
	cpu := makeTestCPU()

	cpu.regs.PC = 0x0004
	cpu.mem.Write16(0x0005, 0xdead)
	cpu.regs.X = 0x02

	cpu.mem.Write8(0xdeaf, 0xfe)

	rol(cpu, AddrModeAbsX)

	if res := cpu.mem.Read8(0xdeaf); res != 0xfd {
		t.Fatalf("Value is not rotated left: %x", res)
	}

	if cpu.getFlag(StatusFlagZ) {
		t.Fatalf("Result should not have been zero")
	}

	if !cpu.getFlag(StatusFlagN) {
		t.Fatalf("Flag N shouldve been set")
	}
}

func TestROR(t *testing.T) {
	cpu := makeTestCPU()

	cpu.regs.A = 0xf0

	ror(cpu, AddrModeAccumulator)

	if res := cpu.regs.A; res != 0x78 {
		t.Fatalf("Value is not rotated left: %x", res)
	}

	if cpu.getFlag(StatusFlagZ) {
		t.Fatalf("Result should not have been zero")
	}

	if cpu.getFlag(StatusFlagN) {
		t.Fatalf("Flag N should have not been set")
	}
}

func TestPHP(t *testing.T) {
	cpu := makeTestCPU()
	cpu.regs.P = 0xab

	php(cpu, AddrModeImplied)

	if res := cpu.Pop8(); res != 0xab {
		t.Fatalf("Pushed wrong value for P: %x", res)
	}
}

func TestPLP(t *testing.T) {
	cpu := makeTestCPU()
	cpu.Push8(0xab)

	plp(cpu, AddrModeImplied)

	if res := cpu.regs.P; res != 0xab {
		t.Fatalf("Popped wrong value for P: %x", res)
	}
}

func TestBMI(t *testing.T) {
	cpu := makeTestCPU()
	cpu.regs.PC = 0xaaaa
	cpu.setFlag(StatusFlagN)

	// -128 in two's complement
	cpu.mem.Write8(0xaaab, 0x80)

	expected := uint16(0xaaaa - 128 + 2)

	bmi(cpu, AddrModeRelative)

	if res := cpu.regs.PC; res != expected {
		t.Fatalf("Popped wrong value for P: %x", res)
	}
}

func TestBNE(t *testing.T) {
	cpu := makeTestCPU()
	cpu.regs.PC = 0xaaaa
	cpu.resetFlag(StatusFlagN)

	// -128 in two's complement
	cpu.mem.Write8(0xaaab, 0x80)

	expected := uint16(0xaaaa - 128 + 2)

	bne(cpu, AddrModeRelative)

	if res := cpu.regs.PC; res != expected {
		t.Fatalf("Popped wrong value for P: %x", res)
	}
}

func TestCMP(t *testing.T) {
	cpu := makeTestCPU()
	cpu.regs.A = 0xf0
	cpu.regs.PC = 0x0004
	cpu.mem.Write8(cpu.regs.PC+1, 0x43)
	cpu.mem.Write8(0x43, 0xa0)
	cpu.regs.Y = 0x0a

	cpu.mem.Write8(0xaa, 0xbb)

	// Will compare 0xf0 with 0xbb
	cmp(cpu, AddrModeIndirectY)

	if cpu.getFlag(StatusFlagZ) {
		t.Fatalf("Flag Z should not have been set")
	}

	// C = A >= m
	if !cpu.getFlag(StatusFlagC) {
		t.Fatalf("Flag C should have been set")
	}

	if cpu.getFlag(StatusFlagN) {
		t.Fatalf("Flag N should not have been set")
	}

	cpu.regs.A = 0x00
	cpu.regs.PC = 0x0004
	cpu.mem.Write8(0x0005, 0x0a)

	cmp(cpu, AddrModeImmediate)

	if cpu.getFlag(StatusFlagZ) {
		t.Fatalf("Flag Z should not have been set")
	}

	if !cpu.getFlag(StatusFlagN) {
		t.Fatalf("Flag N should have been set")
	}

	if cpu.getFlag(StatusFlagC) {
		t.Fatalf("Flag N should not have been set")
	}
}

func TestDEC(t *testing.T) {
	cpu := makeTestCPU()
	cpu.regs.PC = 0x0004
	cpu.mem.Write8(cpu.regs.PC+1, 0x43)
	cpu.mem.Write8(0x0043, 0xf0)

	dec(cpu, AddrModeZeroPage)

	if v := cpu.mem.Read8(0x0043); v != 0xef {
		t.Fatalf("Value was not decremented: %x", v)
	}

	if !cpu.getFlag(StatusFlagN) {
		t.Fatalf("Flag N should have been set")
	}

	if cpu.getFlag(StatusFlagZ) {
		t.Fatalf("Flag Z should not have been set")
	}
}

func to_uint8(v int) uint8 {
	if v < -128 {
		log.Fatalf("Invalid value %x", v)
	} else if v > 127 {
		log.Fatalf("Invalid value %x", v)
	} else if v < 0 {
		return uint8(int(math.Pow(2, 8)) + v)
	} else {
		return uint8(v)
	}
	return 0
}

func TestADC(t *testing.T) {
	cpu := makeTestCPU()
	cpu.regs.PC = 0x0004

	type testCase struct {
		A        uint8
		m        uint8
		result   uint8
		carry    bool
		overflow bool
	}

	cases := []testCase{
		{A: to_uint8(-39), m: to_uint8(92), result: to_uint8(53), carry: true, overflow: false},
		{A: to_uint8(-19), m: to_uint8(-7), result: to_uint8(-26), carry: true, overflow: false},
		{A: to_uint8(44), m: to_uint8(45), result: to_uint8(89), carry: false, overflow: false},
		{A: to_uint8(104), m: to_uint8(45), result: to_uint8(149 - 256), carry: false, overflow: true},
		{A: to_uint8(-103), m: to_uint8(-69), result: to_uint8(256 - 172), carry: true, overflow: true},
	}

	for _, test := range cases {
		cpu.setFlag(StatusFlagC)
		cpu.regs.A = test.A
		cpu.mem.Write8(cpu.regs.PC+1, test.m)

		adc(cpu, AddrModeImmediate)

		if v := cpu.regs.A; v != test.result {
			t.Fatalf("Wrong result: %x", v)
		}

		if cpu.getFlag(StatusFlagC) != test.carry {
			t.Fatalf("Wrong carry")
		}

		if cpu.getFlag(StatusFlagV) != test.overflow {
			t.Fatalf("Wrong overflow")
		}
	}
}

func TestSBC(t *testing.T) {
	cpu := makeTestCPU()
	cpu.regs.PC = 0x0004

	type testCase struct {
		A        uint8
		m        uint8
		result   uint8
		carry    bool
		overflow bool
	}

	cases := []testCase{
		{A: to_uint8(92), m: to_uint8(39), result: to_uint8(53), carry: true, overflow: false},
		{A: to_uint8(-19), m: to_uint8(7), result: to_uint8(-26), carry: true, overflow: false},
		{A: to_uint8(44), m: to_uint8(-45), result: to_uint8(89), carry: false, overflow: false},
		{A: to_uint8(104), m: to_uint8(-45), result: to_uint8(149 - 256), carry: false, overflow: true},
		{A: to_uint8(-103), m: to_uint8(69), result: to_uint8(256 - 172), carry: true, overflow: true},
	}

	for _, test := range cases {
		cpu.setFlag(StatusFlagC)
		cpu.regs.A = test.A
		cpu.mem.Write8(cpu.regs.PC+1, test.m)

		sbc(cpu, AddrModeImmediate)

		if v := cpu.regs.A; v != test.result {
			t.Fatalf("Wrong result: %x", v)
		}

		if cpu.getFlag(StatusFlagC) != test.carry {
			t.Fatalf("Wrong carry")
		}

		if cpu.getFlag(StatusFlagV) != test.overflow {
			t.Fatalf("Wrong overflow")
		}
	}
}
