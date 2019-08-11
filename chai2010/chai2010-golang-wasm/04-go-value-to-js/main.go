// 04-go-value-to-js/main.go
package main

import (
	"syscall/js"
)

func main() {
	js.Global().Set("a_bool", js.ValueOf(true))
	js.Global().Set("a_int", js.ValueOf(123))
	js.Global().Set("a_string", js.ValueOf("abc"))

	js.Global().Call("eval", `
		console.log(typeof a_bool, a_bool);
		console.log(typeof a_int, a_int);
		console.log(typeof a_string, a_string);
	`)
}
