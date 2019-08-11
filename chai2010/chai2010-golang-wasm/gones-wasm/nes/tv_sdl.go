// +build !js

package nes

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCREEN_WIDTH  = 256
	SCREEN_HEIGHT = 240
)

type TV struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
}

func MakeTV() *TV {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(
		"Awesomenes",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		SCREEN_WIDTH*3,
		SCREEN_HEIGHT*3,
		sdl.WINDOW_SHOWN)

	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)

	texture, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, SCREEN_WIDTH, SCREEN_HEIGHT)

	if err != nil {
		panic(err)
	}

	renderer.SetLogicalSize(SCREEN_WIDTH, SCREEN_HEIGHT)

	// Initialize joystick if present
	sdl.JoystickOpen(0)

	return &TV{
		window:   window,
		renderer: renderer,
		texture:  texture,
	}
}

func (tv *TV) UpdateInputState(ctrlr *Controller) {
	for evt := sdl.PollEvent(); evt != nil; evt = sdl.PollEvent() {
		switch evt.(type) {
		case *sdl.KeyboardEvent:
			tv.handleKBDEvevent(ctrlr, evt.(*sdl.KeyboardEvent))

		case *sdl.JoyHatEvent:
			tv.handleJoyHatEvent(ctrlr, evt.(*sdl.JoyHatEvent))

		case *sdl.JoyButtonEvent:
			tv.handleJoyButtonEvent(ctrlr, evt.(*sdl.JoyButtonEvent))
		}
	}
}

func (tv *TV) handleKBDEvevent(ctrlr *Controller, evt *sdl.KeyboardEvent) {
	if evt.Repeat != 0 {
		return
	}

	fn := ctrlr.PushButton

	if evt.Type == sdl.KEYUP {
		fn = ctrlr.ReleaseButton
	}

	switch evt.Keysym.Sym {
	case sdl.K_RETURN:
		fn(CONTROLLER_BUTTONS_START)
	case sdl.K_RSHIFT:
		fn(CONTROLLER_BUTTONS_SELECT)
	case sdl.K_a:
		fn(CONTROLLER_BUTTONS_A)
	case sdl.K_s:
		fn(CONTROLLER_BUTTONS_B)
	case sdl.K_UP:
		fn(CONTROLLER_BUTTONS_UP)
	case sdl.K_RIGHT:
		fn(CONTROLLER_BUTTONS_RIGHT)
	case sdl.K_DOWN:
		fn(CONTROLLER_BUTTONS_DOWN)
	case sdl.K_LEFT:
		fn(CONTROLLER_BUTTONS_LEFT)
	}
}

func (tv *TV) handleJoyHatEvent(ctrlr *Controller, evt *sdl.JoyHatEvent) {
	v := evt.Value

	pressOrRelease := func(SDL_BTN uint8, CTRL_BTN uint8) {
		if v&SDL_BTN == 0 {
			ctrlr.ReleaseButton(CTRL_BTN)
		} else {
			ctrlr.PushButton(CTRL_BTN)
		}
	}

	pressOrRelease(sdl.HAT_UP, CONTROLLER_BUTTONS_UP)
	pressOrRelease(sdl.HAT_RIGHT, CONTROLLER_BUTTONS_RIGHT)
	pressOrRelease(sdl.HAT_DOWN, CONTROLLER_BUTTONS_DOWN)
	pressOrRelease(sdl.HAT_LEFT, CONTROLLER_BUTTONS_LEFT)
}

func (tv *TV) handleJoyButtonEvent(ctrlr *Controller, evt *sdl.JoyButtonEvent) {
	pressOrRelease := func(SDL_BTN uint8, CTRL_BTN uint8) {
		if evt.Button == SDL_BTN {
			if evt.State == sdl.RELEASED {
				ctrlr.ReleaseButton(CTRL_BTN)
			} else {
				ctrlr.PushButton(CTRL_BTN)
			}
		}
	}

	// Down arrow
	pressOrRelease(0, CONTROLLER_BUTTONS_B)
	// Right arrow
	pressOrRelease(1, CONTROLLER_BUTTONS_A)
	// SL
	pressOrRelease(4, CONTROLLER_BUTTONS_SELECT)
	// SR
	pressOrRelease(5, CONTROLLER_BUTTONS_START)
}

func (tv *TV) SetFrame(pixels []byte) {
	tv.texture.Update(nil, pixels, SCREEN_WIDTH*4)
}

func (tv *TV) ShowPixels() {
	tv.renderer.Clear()
	tv.renderer.Copy(tv.texture, nil, nil)
	tv.renderer.Present()
}

func (tv *TV) Cleanup() {
	tv.window.Destroy()
	sdl.Quit()
}
