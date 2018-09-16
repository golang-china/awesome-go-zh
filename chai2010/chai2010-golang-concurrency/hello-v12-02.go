// +build ignore

package main

import "time"

func main() {
	for i := 0; i < 10; i++ {
		go func(i int) { println(i) }(i) // 函数参数是传值 // HL
	}
	time.Sleep(time.Second)
}
