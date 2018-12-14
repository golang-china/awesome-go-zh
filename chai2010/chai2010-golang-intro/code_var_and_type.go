// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// START // OMIT
weekend := []string{ "Saturday", "Sunday" } // HL

timeZones := map[string]TZ { // HL
	"UTC":UTC, "EST":EST, "CST":CST, //...
}

func add(a, b int) int { return a+b } // HL

type Op func (int, int) int // HL

type RPC struct { // HL
	a, b int;
	op Op;
	result *int;
}

rpc := RPC{ 1, 2, add, new(int) }; // HL
