// +build ignore

package main

import "time"

func main() {
	for i := 0; i < 10; i++ {
		go func() { println(i) }() // 对比多次执行结果 // HL
	}
	time.Sleep(time.Second)
}
