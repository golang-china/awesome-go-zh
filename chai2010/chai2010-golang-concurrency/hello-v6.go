// +build ignore

package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(1) // HL

	go func() {
		for i := 0; ; i++ {
			fmt.Println(i)
		}
	}()

	select{} // 阻塞, 不占用CPU // HL
}
