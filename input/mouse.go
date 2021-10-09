package input

import "github.com/veandco/go-sdl2/sdl"

type mouseBtnState struct {
	button            byte
	state             byte
	isDoubleClick     bool
	pressedThisFrame  bool
	releasedThisFrame bool
}

var (
	mouseBtns = make(map[byte]*mouseBtnState)
)

func HandleMouseBtnEvent(e *sdl.MouseButtonEvent) {

	mb := mouseBtns[e.Button]
	if mb == nil {
		mb = &mouseBtnState{button: e.Button}
		mouseBtns[e.Button] = mb
	}

	mb.state = e.State
	mb.isDoubleClick = e.Clicks > 1 && e.State == sdl.PRESSED
	mb.pressedThisFrame = e.State == sdl.PRESSED
	mb.releasedThisFrame = e.State == sdl.RELEASED
}

func MouseClicked(mouseBtn byte) bool {

	mb := mouseBtns[mouseBtn]
	if mb == nil {
		return false
	}

	return mb.pressedThisFrame
}

func MouseDoubleClicked(mouseBtn byte) bool {

	mb := mouseBtns[mouseBtn]
	if mb == nil {
		return false
	}

	return mb.isDoubleClick
}

func MouseReleased(mouseBtn byte) bool {

	mb := mouseBtns[mouseBtn]
	if mb == nil {
		return false
	}

	return mb.releasedThisFrame
}

func MouseDown(mouseBtn byte) bool {

	mb := mouseBtns[mouseBtn]
	if mb == nil {
		return false
	}

	return mb.state == sdl.PRESSED
}

func MouseUp(mouseBtn byte) bool {

	mb := mouseBtns[mouseBtn]
	if mb == nil {
		return true
	}

	return mb.state == sdl.RELEASED
}
