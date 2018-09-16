// +build ignore

package main

import (
	"time"
)

func main() {
	go println("你好, 并发!") // 干活的

	go func() { <-make(chan int) } () // 滥竽充数的, Goroutine 泄露 // HL
	go func() { for{} } () // 浪费资源的, 但不是 Goroutine 泄露 // HL
	go func() {} () // 滥竽充数的, 但不是 Goroutine 泄露 // HL

	time.Sleep(time.Second)
	println("Done")
}
