// +build ignore

package main

import (
    "sync"
)

var msg string
var done sync.Mutex // HL

func main() {
    done.Lock() // HL
    go func() {
        msg = "hello, world"
        done.Unlock() // HL
    }()

    done.Lock() // HL
    println(msg)
}
