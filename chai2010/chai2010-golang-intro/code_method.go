// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// START // OMIT
type Point struct { // HL
	X, Y float // Upper case means exported
}

func (p *Point) Scale(s float) { // HL
	p.X *= s; p.Y *= s; // p is explicit
}

func (p *Point) Abs() float { // HL
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

x := &Point{ 3, 4 };
x.Scale(5); // HL
