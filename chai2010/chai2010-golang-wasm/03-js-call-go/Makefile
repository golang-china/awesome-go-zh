default:
	GOOS=js GOARCH=wasm go build -o hello.wasm
	node $(shell go env GOROOT)/misc/wasm/wasm_exec.js hello.wasm 

clean:
