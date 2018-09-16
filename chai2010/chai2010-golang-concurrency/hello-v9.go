// +build ignore

package main

func main() {
	done := make(chan bool) // HL
	go func() {
		println("你好, 并发!")
		done <- true // HL
	}()

	<-done // HL
}
