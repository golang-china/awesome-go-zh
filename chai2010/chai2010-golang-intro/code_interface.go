// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// START // OMIT
type Abser interface { // 接口名一般以 "er" 为后缀 // HL
	Abs() float;
}

var m Abser;
m = x; // x是 *Point 类型, 已经实现了 Abs() 方法 (隐式转换) // HL
mag := m.Abs(); // HL

type Point3 struct { X, Y, Z float }
func (p *Point3) Abs() float { // HL
	return math.Sqrt(p.X*p.X + p.Y*p.Y + p.Z*p.Z)
}

m = &Point3{ 3, 4, 5 }; // 自动转换为接口类型 // HL
mag += m.Abs();

type Polar struct { R, 囧 float }
func (p Polar) Abs() float { return p.R }

m = Polar{ 2.0, PI/2 }; // 自动转换为接口类型 // HL
mag += m.Abs();
