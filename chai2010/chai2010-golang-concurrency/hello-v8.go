// +build ignore

package main

func main() {
	go println("你好, 并发!") // HL
	<-make(chan int) // 阻塞, 不占用CPU // HL
}
