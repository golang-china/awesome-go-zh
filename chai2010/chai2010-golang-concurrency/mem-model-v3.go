// +build ignore

package main

import (
	"runtime"
)

var msg string
var done bool = false

func main() {
    runtime.GOMAXPROCS(2) // HL

    go func() {
        msg = "hello, world"
        done = true // HL
    }()

    for {
        if done { // HL
            println(msg); break
        }
		println("retry...")
    }
}
