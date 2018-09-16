// +build ignore

package main

var msg string
var done = make(chan struct{}) // HL

func main() {
    go func() {
        msg = "hello, world"
        done <- struct{}{} // HL
    }()

    <-done // HL
    println(msg)
}
