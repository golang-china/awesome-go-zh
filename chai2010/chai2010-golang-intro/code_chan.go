// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

var c chan string;
c = make(chan string); // HL

c <- "Hello"; // infix send // HL
// in a different goroutine
greeting := <-c; // prefix receive // HL

cc := new(chan chan string);
cc <- c; // handing off a capability
