package input

import "github.com/veandco/go-sdl2/sdl"

type keyState struct {
	key               sdl.Keycode
	state             byte
	pressedThisFrame  bool
	releasedThisFrame bool
}

var (
	keyboardKeys = make(map[sdl.Keycode]*keyState)
)

func HandleKeyboardEvent(e *sdl.KeyboardEvent) {

	ks := keyboardKeys[e.Keysym.Sym]
	if ks == nil {
		ks = &keyState{key: e.Keysym.Sym}
		keyboardKeys[e.Keysym.Sym] = ks
	}

	ks.state = e.State
	ks.pressedThisFrame = e.Repeat == 0 && e.State == sdl.PRESSED
	ks.releasedThisFrame = e.Repeat == 0 && e.State == sdl.RELEASED
}

func KeyClicked(kc sdl.Keycode) bool {

	key := keyboardKeys[kc]
	if key == nil {
		return false
	}

	return key.pressedThisFrame
}

func KeyReleased(kc sdl.Keycode) bool {

	key := keyboardKeys[kc]
	if key == nil {
		return false
	}

	return key.releasedThisFrame
}

func KeyDown(kc sdl.Keycode) bool {

	key := keyboardKeys[kc]
	if key == nil {
		return false
	}

	return key.state == sdl.PRESSED
}

func KeyUp(kc sdl.Keycode) bool {

	key := keyboardKeys[kc]
	if key == nil {
		return true
	}

	return key.state == sdl.RELEASED
}
