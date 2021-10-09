package input

import "github.com/veandco/go-sdl2/sdl"

type InputKey int

var (
	anyKeyDown      bool
	anyMouseBtnDown bool
)

//EventLoopStarted should be called just before processing SDL events
func EventLoopStarted() {

	anyKeyDown = false
	anyMouseBtnDown = false

	//Clear XThisFrame which is needed because a repeat event needs multiple frames to happen
	for _, v := range mouseBtns {

		v.isDoubleClick = false
		v.pressedThisFrame = false
		v.releasedThisFrame = false

		if v.state == sdl.PRESSED {
			anyMouseBtnDown = true
		}
	}

	for _, v := range keyboardKeys {

		v.pressedThisFrame = false
		v.releasedThisFrame = false

		if v.state == sdl.PRESSED {
			anyKeyDown = true
		}
	}
}

func AnyKeyDown() bool {
	return anyKeyDown
}

func AnyMouseBtnDown() bool {
	return anyMouseBtnDown
}
