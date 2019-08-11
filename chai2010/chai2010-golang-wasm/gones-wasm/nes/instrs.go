package nes

import (
	"log"
)

type instr struct {
	name     string
	size     uint8
	cycles   uint8
	addrMode addressingMode
	fn       func(cpu *CPU, addrMode addressingMode)
}

type addressingMode uint8

const (
	AddrModeAbs addressingMode = iota
	AddrModeAbsX
	AddrModeAbsY
	AddrModeAccumulator
	AddrModeImmediate
	AddrModeImplied
	AddrModeIndirect
	AddrModeXIndirect
	AddrModeIndirectY
	AddrModeRelative
	AddrModeZeroPage
	AddrModeZeroX
	AddrModeZeroY
)

func calculateAddr(cpu *CPU, addrMode addressingMode) uint16 {
	switch addrMode {
	case AddrModeAbs:
		return cpu.mem.Read16(cpu.regs.PC + 1)

	case AddrModeAbsX:
		return cpu.mem.Read16(cpu.regs.PC+1) + uint16(cpu.regs.X)

	case AddrModeAbsY:
		return cpu.mem.Read16(cpu.regs.PC+1) + uint16(cpu.regs.Y)

	case AddrModeAccumulator:
		log.Fatalf("It makes no sense to calculate addresses in accumulator addressing mode")
		return 0xffff

	case AddrModeImmediate:
		return cpu.regs.PC + 1

	case AddrModeRelative:
		offset := uint16(cpu.mem.Read8(cpu.regs.PC + 1))
		if offset < 0x80 {
			return cpu.regs.PC + 2 + offset
		} else {
			return cpu.regs.PC + 2 + offset - 0x100
		}

	case AddrModeIndirect:
		m := cpu.mem.Read16(cpu.regs.PC + 1)
		return cpu.read16AndMaybeWrap(m)

	case AddrModeXIndirect:
		m := cpu.mem.Read8(cpu.regs.PC + 1)
		x := cpu.regs.X
		return cpu.read16AndMaybeWrap(uint16(m + x))

	case AddrModeIndirectY:
		m := cpu.mem.Read8(cpu.regs.PC + 1)
		y := cpu.regs.Y
		return cpu.read16AndMaybeWrap(uint16(m)) + uint16(y)

	case AddrModeZeroPage:
		return uint16(cpu.mem.Read8(cpu.regs.PC + 1))

	case AddrModeZeroX:
		return uint16(cpu.mem.Read8(cpu.regs.PC+1) + cpu.regs.X)

	default:
		log.Fatalf("Invalid addressing mode")
		return 0xffff
	}
}

// Break
func brk(cpu *CPU, addrMode addressingMode) {
	cpu.Push16(cpu.regs.PC)
	cpu.Push8(cpu.regs.P)
	cpu.setFlag(StatusFlagB)
	cpu.regs.PC = cpu.mem.Read16(0xfffe)
	cpu.jumped = true
}

// Logical inclusive OR
func ora(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	cpu.regs.A = cpu.regs.A | cpu.mem.Read8(addr)
	cpu.setOrReset(StatusFlagN, cpu.regs.A&0x80 != 0)
	cpu.setOrReset(StatusFlagZ, cpu.regs.A == 0)
}

func eor(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	cpu.regs.A = cpu.regs.A ^ cpu.mem.Read8(addr)
	cpu.setOrReset(StatusFlagN, cpu.regs.A&0x80 != 0)
	cpu.setOrReset(StatusFlagZ, cpu.regs.A == 0)
}

// Arithmetic shift left
func asl(cpu *CPU, addrMode addressingMode) {
	shiftL := func(v uint8) uint8 {
		cpu.setOrReset(StatusFlagC, v&0x80 != 0)
		v = v << 1
		setOrResetNZ(v, cpu)
		return v
	}

	if addrMode == AddrModeAccumulator {
		cpu.regs.A = shiftL(cpu.regs.A)
	} else {
		addr := calculateAddr(cpu, addrMode)
		cpu.mem.Write8(addr, shiftL(cpu.mem.Read8(addr)))
	}
}

// Logical shift right
func lsr(cpu *CPU, addrMode addressingMode) {
	shiftR := func(v uint8) uint8 {
		cpu.setOrReset(StatusFlagC, v&0x1 == 0x1)
		v = v >> 1
		setOrResetNZ(v, cpu)
		return v
	}

	if addrMode == AddrModeAccumulator {
		cpu.regs.A = shiftR(cpu.regs.A)
	} else {
		addr := calculateAddr(cpu, addrMode)
		cpu.mem.Write8(addr, shiftR(cpu.mem.Read8(addr)))
	}
}

// Add with carry - borrowed from github.com/fogleman/nes
func adc(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	a := cpu.regs.A
	b := cpu.mem.Read8(addr)
	var c uint8 = 0
	if cpu.getFlag(StatusFlagC) {
		c = 1
	}
	cpu.regs.A = a + b + c
	setOrResetNZ(cpu.regs.A, cpu)

	cpu.setOrReset(StatusFlagC, int(a)+int(b)+int(c) > 0xFF)
	cpu.setOrReset(StatusFlagV, (a^b)&0x80 == 0 && (a^cpu.regs.A)&0x80 != 0)
}

// Subtract with carry - borrowed from github.com/fogleman/nes
func sbc(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	a := cpu.regs.A
	b := cpu.mem.Read8(addr)
	var c uint8 = 0
	if cpu.getFlag(StatusFlagC) {
		c = 1
	}
	cpu.regs.A = a - b - (1 - c)
	setOrResetNZ(cpu.regs.A, cpu)

	cpu.setOrReset(StatusFlagC, int(a)-int(b)-int(1-c) >= 0x00)
	cpu.setOrReset(StatusFlagV, (a^b)&0x80 != 0 && (a^cpu.regs.A)&0x80 != 0)
}

// Push processor state
// See http://wiki.nesdev.com/w/index.php/CPU_status_flag_behavior
// for the meaning of | 0x10. Hint: it has none
func php(cpu *CPU, addrMode addressingMode) {
	cpu.Push8(cpu.regs.P | 0x10)
}

func pha(cpu *CPU, addrMode addressingMode) {
	cpu.Push8(cpu.regs.A)
}

// Pull processor state
func plp(cpu *CPU, addrMode addressingMode) {
	cpu.regs.P = cpu.Pop8()&0xEF | 0x20
}

func pla(cpu *CPU, addrMode addressingMode) {
	cpu.regs.A = cpu.Pop8()
	setOrResetNZ(cpu.regs.A, cpu)
}

// Clear carry
func clc(cpu *CPU, addrMode addressingMode) {
	cpu.resetFlag(StatusFlagC)
}

func cld(cpu *CPU, addrMode addressingMode) {
	cpu.resetFlag(StatusFlagD)
}

func cli(cpu *CPU, addrMode addressingMode) {
	cpu.resetFlag(StatusFlagI)
}

func clv(cpu *CPU, addrMode addressingMode) {
	cpu.resetFlag(StatusFlagV)
}

// Jump
func jmp(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	cpu.regs.PC = addr
	cpu.jumped = true
}

// Jump to subroutine
func jsr(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	cpu.Push16(cpu.regs.PC - 1 + 3)
	cpu.regs.PC = addr
	cpu.jumped = true
}

// Return from interrupt
func rti(cpu *CPU, addrMode addressingMode) {
	cpu.regs.P = cpu.Pop8()&0xEF | 0x20
	cpu.regs.PC = cpu.Pop16()
	cpu.jumped = true
}

// Return from subroutine
func rts(cpu *CPU, addrMode addressingMode) {
	cpu.regs.PC = cpu.Pop16() + 1
	cpu.jumped = true
}

func and(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	cpu.regs.A = cpu.regs.A & cpu.mem.Read8(addr)
	cpu.setOrReset(StatusFlagN, cpu.regs.A&0x80 != 0)
	cpu.setOrReset(StatusFlagZ, cpu.regs.A == 0)
}

func bit(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	v := cpu.mem.Read8(addr)

	cpu.setOrReset(StatusFlagZ, cpu.regs.A&v == 0)
	cpu.setOrReset(StatusFlagN, (v>>7)&0x1 == 0x1)
	cpu.setOrReset(StatusFlagV, (v>>6)&0x1 == 0x1)
}

// Rotate left
func rol(cpu *CPU, addrMode addressingMode) {
	inner := func(v uint8) uint8 {
		var c uint8 = 0
		if cpu.getFlag(StatusFlagC) {
			c = 1
		}
		cpu.setOrReset(StatusFlagC, (v>>7)&0x1 == 0x1)
		v = (v << 1) | c
		setOrResetNZ(v, cpu)
		return v
	}

	if addrMode == AddrModeAccumulator {
		cpu.regs.A = inner(cpu.regs.A)
	} else {
		addr := calculateAddr(cpu, addrMode)
		cpu.mem.Write8(addr, inner(cpu.mem.Read8(addr)))
	}
}

// Rotate right
func ror(cpu *CPU, addrMode addressingMode) {
	inner := func(v uint8) uint8 {
		var c uint8 = 0
		if cpu.getFlag(StatusFlagC) {
			c = 1
		}
		cpu.setOrReset(StatusFlagC, v&0x1 == 0x1)
		v = (v >> 1) | (c << 7)
		setOrResetNZ(v, cpu)
		return v
	}

	if addrMode == AddrModeAccumulator {
		cpu.regs.A = inner(cpu.regs.A)
	} else {
		addr := calculateAddr(cpu, addrMode)
		cpu.mem.Write8(addr, inner(cpu.mem.Read8(addr)))
	}
}

func bcc(cpu *CPU, addrMode addressingMode) {
	if !cpu.getFlag(StatusFlagC) {
		cpu.regs.PC = calculateAddr(cpu, addrMode)
		cpu.jumped = true
	}
}

func bcs(cpu *CPU, addrMode addressingMode) {
	if cpu.getFlag(StatusFlagC) {
		cpu.regs.PC = calculateAddr(cpu, addrMode)
		cpu.jumped = true
	}
}

// Branch if negative
func bmi(cpu *CPU, addrMode addressingMode) {
	if cpu.getFlag(StatusFlagN) {
		cpu.regs.PC = calculateAddr(cpu, addrMode)
		cpu.jumped = true
	}
}

// Branch if equal
func beq(cpu *CPU, addrMode addressingMode) {
	if cpu.getFlag(StatusFlagZ) {
		cpu.regs.PC = calculateAddr(cpu, addrMode)
		cpu.jumped = true
	}
}

// Branch if not equal
func bne(cpu *CPU, addrMode addressingMode) {
	if !cpu.getFlag(StatusFlagZ) {
		cpu.regs.PC = calculateAddr(cpu, addrMode)
		cpu.jumped = true
	}
}

// Branch if positive
func bpl(cpu *CPU, addrMode addressingMode) {
	if !cpu.getFlag(StatusFlagN) {
		cpu.regs.PC = calculateAddr(cpu, addrMode)
		cpu.jumped = true
	}
}

// Branch if overflow clear
func bvc(cpu *CPU, addrMode addressingMode) {
	if !cpu.getFlag(StatusFlagV) {
		cpu.regs.PC = calculateAddr(cpu, addrMode)
		cpu.jumped = true
	}
}

// Branch if overflow set
func bvs(cpu *CPU, addrMode addressingMode) {
	if cpu.getFlag(StatusFlagV) {
		cpu.regs.PC = calculateAddr(cpu, addrMode)
		cpu.jumped = true
	}
}

func _comp(v1 uint8, v2 uint8, cpu *CPU) {
	v := v1 - v2
	cpu.setOrReset(StatusFlagZ, v == 0)
	cpu.setOrReset(StatusFlagN, v>>7 == 0x1)
	cpu.setOrReset(StatusFlagC, v1 >= v2)
}

// Compare with A
func cmp(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	m := cpu.mem.Read8(addr)
	_comp(cpu.regs.A, m, cpu)
}

// Compare with X
func cpx(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	m := cpu.mem.Read8(addr)
	_comp(cpu.regs.X, m, cpu)
}

// Compare with Y
func cpy(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	m := cpu.mem.Read8(addr)
	_comp(cpu.regs.Y, m, cpu)
}

func dec(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	v := cpu.mem.Read8(addr) - 1
	cpu.mem.Write8(addr, v)
	cpu.setOrReset(StatusFlagZ, v == 0)
	cpu.setOrReset(StatusFlagN, v>>7 == 0x1)
}

func dex(cpu *CPU, addrMode addressingMode) {
	cpu.regs.X -= 1
	cpu.setOrReset(StatusFlagZ, cpu.regs.X == 0)
	cpu.setOrReset(StatusFlagN, cpu.regs.X>>7 == 0x1)
}

func dey(cpu *CPU, addrMode addressingMode) {
	cpu.regs.Y -= 1
	cpu.setOrReset(StatusFlagZ, cpu.regs.Y == 0)
	cpu.setOrReset(StatusFlagN, cpu.regs.Y>>7 == 0x1)
}

func inc(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	v := cpu.mem.Read8(addr) + 1
	cpu.mem.Write8(addr, v)
	cpu.setOrReset(StatusFlagZ, v == 0)
	cpu.setOrReset(StatusFlagN, v>>7 == 0x1)
}

func inx(cpu *CPU, addrMode addressingMode) {
	cpu.regs.X += 1
	cpu.setOrReset(StatusFlagZ, cpu.regs.X == 0)
	cpu.setOrReset(StatusFlagN, cpu.regs.X>>7 == 0x1)
}

func iny(cpu *CPU, addrMode addressingMode) {
	cpu.regs.Y += 1
	cpu.setOrReset(StatusFlagZ, cpu.regs.Y == 0)
	cpu.setOrReset(StatusFlagN, cpu.regs.Y>>7 == 0x1)
}

func lda(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	cpu.regs.A = cpu.mem.Read8(addr)
	cpu.setOrReset(StatusFlagZ, cpu.regs.A == 0)
	cpu.setOrReset(StatusFlagN, cpu.regs.A>>7 == 0x1)
}

func sta(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	cpu.mem.Write8(addr, cpu.regs.A)
}

func ldx(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	cpu.regs.X = cpu.mem.Read8(addr)
	cpu.setOrReset(StatusFlagZ, cpu.regs.X == 0)
	cpu.setOrReset(StatusFlagN, cpu.regs.X>>7 == 0x1)
}

func stx(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	cpu.mem.Write8(addr, cpu.regs.X)
}

func ldy(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	cpu.regs.Y = cpu.mem.Read8(addr)
	cpu.setOrReset(StatusFlagZ, cpu.regs.Y == 0)
	cpu.setOrReset(StatusFlagN, cpu.regs.Y>>7 == 0x1)
}

func sty(cpu *CPU, addrMode addressingMode) {
	addr := calculateAddr(cpu, addrMode)
	cpu.mem.Write8(addr, cpu.regs.Y)
}

func nop(cpu *CPU, addrMode addressingMode) {
}

func sec(cpu *CPU, addrMode addressingMode) {
	cpu.setFlag(StatusFlagC)
}

func sed(cpu *CPU, addrMode addressingMode) {
	cpu.setFlag(StatusFlagD)
}

func sei(cpu *CPU, addrMode addressingMode) {
	cpu.setFlag(StatusFlagI)
}

func setOrResetNZ(v uint8, cpu *CPU) {
	cpu.setOrReset(StatusFlagZ, v == 0)
	cpu.setOrReset(StatusFlagN, v>>7 == 0x1)
}

func tax(cpu *CPU, addrMode addressingMode) {
	cpu.regs.X = cpu.regs.A
	setOrResetNZ(cpu.regs.X, cpu)
}

func tay(cpu *CPU, addrMode addressingMode) {
	cpu.regs.Y = cpu.regs.A
	setOrResetNZ(cpu.regs.Y, cpu)
}

func tsx(cpu *CPU, addrMode addressingMode) {
	cpu.regs.X = cpu.regs.SP
	setOrResetNZ(cpu.regs.X, cpu)
}

func txa(cpu *CPU, addrMode addressingMode) {
	cpu.regs.A = cpu.regs.X
	setOrResetNZ(cpu.regs.A, cpu)
}

func tya(cpu *CPU, addrMode addressingMode) {
	cpu.regs.A = cpu.regs.Y
	setOrResetNZ(cpu.regs.A, cpu)
}
func txs(cpu *CPU, addrMode addressingMode) {
	cpu.regs.SP = cpu.regs.X
}
