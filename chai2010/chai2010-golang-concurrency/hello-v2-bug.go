// +build ignore

package main

import "time"

func main() {
	go println("你好, 并发!")
	time.Sleep(time.Second) // HL
}

func println(s string) { // HL
	time.Sleep(time.Second*2) // HL
	print(s+"\n")
}
