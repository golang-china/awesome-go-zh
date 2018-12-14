// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main // HL

import "os" // HL
import "flag"

var nFlag = flag.Bool("n", false, `no \n`) // HL

func main() { // HL
	flag.Parse()
	s := ""
	for i := 0; i < flag.NArg(); i++ {
		if i > 0 {
			s += " "
		}
		s += flag.Arg(i)
	}
	if !*nFlag {
		s += "\n"
	}
	os.Stdout.WriteString(s)
}
