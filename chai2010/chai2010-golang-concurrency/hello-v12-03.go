// +build ignore

package main

import "time"

func main() {
	for i := 0; i < 10; i++ {
		i := i // 定义新的局部变量, 每次迭代都不同 // HL
		go func() { println(i) }()
	}
	time.Sleep(time.Second)
}
