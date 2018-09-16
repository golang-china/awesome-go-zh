// +build ignore

package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(1) // HL

	go func() { for {} }() // HL

	time.Sleep(time.Second)

	fmt.Println("the answer to life:", 42) // HL
}
