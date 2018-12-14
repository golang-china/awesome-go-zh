// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import "fmt"

type TZ int // HL

const (
	HOUR TZ = 60 * 60
	UTC  TZ = 0 * HOUR
	EST  TZ = -5 * HOUR //...
)

var timeZones = map[string]TZ{"UTC": UTC, "EST": EST}

func (tz TZ) String() string { // Method on TZ (not ptr) // HL
	for name, zone := range timeZones {
		if tz == zone {
			return name
		}
	}
	return fmt.Sprintf("%+d:%02d", tz/3600, (tz%3600)/60)
}

func main() {
	fmt.Println(EST) // Print* know about method String() // HL
	fmt.Println(5 * HOUR / 2)
}

// Output (two lines) EST +2:30
