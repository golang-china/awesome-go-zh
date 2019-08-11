// 02-go-call-js/main.go
package main

import (
	"syscall/js"
)

func main() {
	console_log := js.Global().Get("console").Get("log") // HL
	console_log.Invoke("Hello wasm!") // HL

	js.Global().Call("eval", `
		console.log("hello, wasm!");
	`)
}
