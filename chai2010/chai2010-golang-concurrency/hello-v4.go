// +build ignore

package main

import (
	"fmt"
)

func main() {
	go func() {
		for i := 0; ; i++ {
			fmt.Println(i)
		}
	}()

	for {} // 死循环是大杀器 // HL
}
