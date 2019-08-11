// 05-js-value-to-go/main.go
package main

import (
	"syscall/js"
)

func main() {
	js.Global().Call("eval", `
		a_bool = true;
		a_int = 123;
		a_string = 'abc';
	`)

	println(js.Global().Get("a_bool").Bool())     // HL
	println(js.Global().Get("a_int").Int())       // HL
	println(js.Global().Get("a_string").String()) // HL
}
