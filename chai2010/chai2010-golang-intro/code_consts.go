// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// START // OMIT
type TZ int // HL

const (
	UTC TZ = 0*60*60 // HL
	EST TZ = -5*60*60
)

// iota枚举:
const (
	bit0, mask0 uint32 = 1<<iota, 1<<iota - 1 // HL
	bit1, mask1 uint32 = 1<<iota, 1<<iota - 1
	bit2, mask2 // 缺省时, 和上一行相同 // HL
)

// 高精度:
const Ln2= 0.693147180559945309417232121458176568075500134360255254120680009
const Log2E= 1/Ln2 // 高精度 // HL
