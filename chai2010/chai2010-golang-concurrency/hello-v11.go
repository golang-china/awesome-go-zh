// +build ignore

package main

import "sync"

func main() {
	const N = 10
	var wg sync.WaitGroup // HL

	for i := 0; i < N; i++ {
		wg.Add(1) // 必须在 go 语句前调用! // HL
		go func(i int) {
			defer wg.Done() // HL
			println(i, "你好, 并发!")
		}(i)
	}

	wg.Wait() // HL
}
