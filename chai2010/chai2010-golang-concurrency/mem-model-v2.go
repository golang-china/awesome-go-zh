// +build ignore

package main

var msg string
var done bool = false

func main() {
	msg = "hello, world"
	done = true

	for {
		if done { // HL
			println(msg)
			break
		}
        println("retry...")
	}
}
