// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// START // OMIT
x := longCalculation(17); // runs too long
c := make(chan int); // HL

func wrapper(a int, c chan int) { // HL
	result := longCalculation(a);
	c <- result; // HL
}
go wrapper(17, c); // HL

// do something for a while; then...
x := <-c; // HL
