// +build ignore

package main

import (
	"fmt"
	"time"
)

func main() {
	go func() { // HL
		for i := 0; ; i++ { // HL
			fmt.Println(i) // HL
		} // HL
	}() // HL

	time.Sleep(time.Second)
}
