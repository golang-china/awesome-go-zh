package nes

import "fmt"

type Memory []uint8

func (m Memory) Read8(p uint16) uint8 {
	return m[p]
}

func (m Memory) Write8(p uint16, v uint8) {
	m[p] = v
}

func (m Memory) Dump(addr, amount uint16) {
	var row, col uint16

	for ; row < amount/8; row++ {
		fmt.Printf("%04x  ", addr+row*8)
		for col = 0; col < 8; col++ {
			fmt.Printf("%02x ", m.Read8(row*16+col))
		}
		fmt.Println("")
	}
}
