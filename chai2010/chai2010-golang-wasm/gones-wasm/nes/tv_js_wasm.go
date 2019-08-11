// Copyright 2018 chaishushan@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.11

package nes

import (
	"syscall/js"
)

const (
	SCREEN_WIDTH  = 256
	SCREEN_HEIGHT = 240
)

type TV struct {
	window   interface{}
	renderer interface{}
	texture  interface{}
	pixels   []byte
}

func MakeTV() *TV {
	return &TV{}
}

func (tv *TV) SetFrame(pixels []byte) {
	tv.pixels = pixels
	jsSetPixels("nes", tv.pixels, SCREEN_WIDTH, SCREEN_HEIGHT)
}

func (tv *TV) ShowPixels() {
	jsSetPixels("nes", tv.pixels, SCREEN_WIDTH, SCREEN_HEIGHT)
}

func (tv *TV) UpdateInputState(ctrlr *Controller) {
	// TODO
}

func (tv *TV) Cleanup() {
	// TODO
}

func (tv *TV) handleKBDEvevent(ctrlr *Controller, evt interface{}) {
	// TODO
}
func (tv *TV) handleJoyHatEvent(ctrlr *Controller, evt interface{}) {
	// TODO
}
func (tv *TV) handleJoyButtonEvent(ctrlr *Controller, evt interface{}) {
	// TODO
}

func (tv *TV) KeyDown(ctrlr *Controller, key string) {
	fn := ctrlr.PushButton

	switch key {
	case "Enter":
		fn(CONTROLLER_BUTTONS_START)
	case "Shift":
		fn(CONTROLLER_BUTTONS_SELECT)
	case "1":
		fn(CONTROLLER_BUTTONS_A)
	case "2":
		fn(CONTROLLER_BUTTONS_B)
	case "ArrowUp":
		fn(CONTROLLER_BUTTONS_UP)
	case "ArrowRight":
		fn(CONTROLLER_BUTTONS_RIGHT)
	case "ArrowDown":
		fn(CONTROLLER_BUTTONS_DOWN)
	case "ArrowLeft":
		fn(CONTROLLER_BUTTONS_LEFT)
	}
}

func (tv *TV) KeyUp(ctrlr *Controller, key string) {
	fn := ctrlr.ReleaseButton

	switch key {
	case "Enter":
		fn(CONTROLLER_BUTTONS_START)
	case "Shift":
		fn(CONTROLLER_BUTTONS_SELECT)
	case "1":
		fn(CONTROLLER_BUTTONS_A)
	case "2":
		fn(CONTROLLER_BUTTONS_B)
	case "ArrowUp":
		fn(CONTROLLER_BUTTONS_UP)
	case "ArrowRight":
		fn(CONTROLLER_BUTTONS_RIGHT)
	case "ArrowDown":
		fn(CONTROLLER_BUTTONS_DOWN)
	case "ArrowLeft":
		fn(CONTROLLER_BUTTONS_LEFT)
	}
}

func Alert() {
	alert := js.Global().Get("alert")
	alert.Invoke("Hello wasm!")
}

var pixbuf = make([]byte, SCREEN_WIDTH*SCREEN_HEIGHT*4)

func jsSetPixels(canvas_id string, pixel []byte, width, height int) {
	if len(pixel) == 0 {
		return
	}

	if true {
		for i := 0; i < width*height; i++ {
			a := pixel[i*4+0]
			b := pixel[i*4+1]
			g := pixel[i*4+2]
			r := pixel[i*4+3]

			pixbuf[i*4+0] = r
			pixbuf[i*4+1] = g
			pixbuf[i*4+2] = b
			pixbuf[i*4+3] = a
		}
	}

	// 更新js空间的pix缓存信息
	js.Global().Set("nes_width", js.ValueOf(width))
	js.Global().Set("nes_height", js.ValueOf(height))
	js.Global().Call("eval", `nes_pix = new Uint8Array(4*nes_width*nes_height);`)

	nes_pix := js.Global().Get("nes_pix")
	js.CopyBytesToJS(nes_pix, pixbuf)

	jsSetPixels := js.Global().Get("jsSetPixels")
	//	jsSetPixels.Invoke(canvas_id, js.TypedArrayOf(pixbuf), width, height)
	jsSetPixels.Invoke(canvas_id, nes_pix, width, height)
}
