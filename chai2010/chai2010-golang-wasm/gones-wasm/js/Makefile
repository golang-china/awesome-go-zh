go_js_wasm_exec:=$(shell go env GOROOT)/misc/wasm/go_js_wasm_exec

default:
	GOOS=js GOARCH=wasm go build -o index.out.wasm
	#@$(go_js_wasm_exec) index.out.wasm
	#node wasm_exec.js index.out.wasm

dev:
	go build -o a.out
	./a.out

run:
	node wasm_exec.js index.out.wasm

server:
	go run server.go
