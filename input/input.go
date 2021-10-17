package input

import "github.com/veandco/go-sdl2/sdl"

type keyState struct {
	Key                 sdl.Keycode
	State               int
	IsPressedThisFrame  bool
	IsReleasedThisFrame bool
}

type mouseBtnState struct {
	Btn   int
	State int

	IsPressedThisFrame  bool
	IsReleasedThisFrame bool
	IsDoubleClicked     bool
}

var (
	keyMap      = make(map[sdl.Keycode]*keyState)
	mouseBtnMap = make(map[int]*mouseBtnState)
)

func EventLoopStart() {

	for _, v := range keyMap {
		v.IsPressedThisFrame = false
		v.IsReleasedThisFrame = false
	}

	for _, v := range mouseBtnMap {
		v.IsPressedThisFrame = false
		v.IsReleasedThisFrame = false
		v.IsDoubleClicked = false
	}
}

func HandleKeyboardEvent(e *sdl.KeyboardEvent) {

	ks := keyMap[e.Keysym.Sym]
	if ks == nil {
		ks = &keyState{Key: e.Keysym.Sym}
		keyMap[ks.Key] = ks
	}

	ks.State = int(e.State)
	ks.IsPressedThisFrame = e.State == sdl.PRESSED && e.Repeat == 0
	ks.IsReleasedThisFrame = e.State == sdl.RELEASED && e.Repeat == 0
}

func HandleMouseEvent(e *sdl.MouseButtonEvent) {

	mb := mouseBtnMap[int(e.Button)]
	if mb == nil {
		mb = &mouseBtnState{Btn: int(e.Button)}
		mouseBtnMap[int(e.Button)] = mb
	}

	mb.State = int(e.State)

	mb.IsDoubleClicked = e.Clicks == 2 && e.State == sdl.PRESSED
	mb.IsPressedThisFrame = e.State == sdl.PRESSED
	mb.IsReleasedThisFrame = e.State == sdl.RELEASED
}

func KeyClicked(kc sdl.Keycode) bool {

	ks := keyMap[kc]
	if ks == nil {
		return false
	}

	return ks.IsPressedThisFrame
}

func KeyReleased(kc sdl.Keycode) bool {

	ks := keyMap[kc]
	if ks == nil {
		return false
	}

	return ks.IsReleasedThisFrame
}

func KeyDown(kc sdl.Keycode) bool {

	ks := keyMap[kc]
	if ks == nil {
		return false
	}

	return ks.State == sdl.PRESSED
}

func KeyUp(kc sdl.Keycode) bool {

	ks := keyMap[kc]
	if ks == nil {
		return true
	}

	return ks.State == sdl.RELEASED
}

func MouseClicked(mb int) bool {

	btn := mouseBtnMap[mb]
	if btn == nil {
		return false
	}

	return btn.IsPressedThisFrame
}

func MouseDoubleClicked(mb int) bool {

	btn := mouseBtnMap[mb]
	if btn == nil {
		return false
	}

	return btn.IsDoubleClicked
}

func MouseReleased(mb int) bool {
	btn := mouseBtnMap[mb]
	if btn == nil {
		return false
	}

	return btn.IsReleasedThisFrame
}

func MouseDown(mb int) bool {

	btn := mouseBtnMap[mb]
	if btn == nil {
		return false
	}

	return btn.State == sdl.PRESSED
}

func MouseUp(mb int) bool {

	btn := mouseBtnMap[mb]
	if btn == nil {
		return true
	}

	return btn.State == sdl.RELEASED
}
