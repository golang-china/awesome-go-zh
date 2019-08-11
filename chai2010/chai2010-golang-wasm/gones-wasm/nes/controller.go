package nes

type Controller struct {
	// Latch containing the buttons state in order from MSB to LSB:
	// A, B, Select, Start, Up, Down, Left, Right
	state  uint8
	strobe bool
}

const (
	CONTROLLER_BUTTONS_A = iota
	CONTROLLER_BUTTONS_B
	CONTROLLER_BUTTONS_SELECT
	CONTROLLER_BUTTONS_START
	CONTROLLER_BUTTONS_UP
	CONTROLLER_BUTTONS_DOWN
	CONTROLLER_BUTTONS_LEFT
	CONTROLLER_BUTTONS_RIGHT
)

func MakeController() *Controller {
	return &Controller{
		state:  0,
		strobe: false,
	}
}

func (ctrl *Controller) ReadState() uint8 {
	if ctrl.strobe {
		// Always return A state during strobe
		return ctrl.state & 0x80
	} else {
		buttonState := (ctrl.state >> 7)
		ctrl.state <<= 1
		ctrl.state |= buttonState

		// http://wiki.nesdev.com/w/index.php/Standard_controller
		return 0x40 | buttonState
	}
}

func (ctrl *Controller) Write(v uint8) {
	ctrl.strobe = v&0x1 == 0x1
}

func (ctrl *Controller) PushButton(btn uint8) {
	pos := 7 - btn
	ctrl.state |= (0x1 << pos)
}

func (ctrl *Controller) ReleaseButton(btn uint8) {
	pos := 7 - btn
	ctrl.state &= ^(0x1 << pos)
}
