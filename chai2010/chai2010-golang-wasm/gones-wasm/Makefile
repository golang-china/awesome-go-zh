# Copyright 2018 <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

nes.out:
	go build -o nes.out

wasm:
	GO111MODULE=on GOOS=js GOARCH=wasm go build -o nes.out.wasm

run: nes.out
	./nes.out #./testdata/SuperMarioBros.nes

clean:
	-rm nes.out
