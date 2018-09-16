// +build ignore

package main

func main() {
	const N = 10
	done := make(chan bool, N) // HL

	for i := 0; i < N; i++ {
		go func(i int) {
			println(i, "你好, 并发!")
			done <- true // HL
		}(i)
	}

	for i := 0; i < N; i++ {
		<-done // HL
	}
}
