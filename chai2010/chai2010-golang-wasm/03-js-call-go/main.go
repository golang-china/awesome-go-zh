// 03-js-call-go/main.go
package main

import (
	"syscall/js"
)

func main() {
	js.Global().Set("println",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} { // HL
			println("hello callback")
			return nil
		}),
	)

	println := js.Global().Get("println")
	println.Invoke()
}
