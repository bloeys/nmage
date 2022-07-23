package engine

import (
	"runtime"

	"github.com/bloeys/nmage/assert"
	"github.com/bloeys/nmage/input"
	"github.com/bloeys/nmage/renderer"
	"github.com/bloeys/nmage/timing"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/inkyblackness/imgui-go/v4"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	isInited = false
)

type Window struct {
	SDLWin         *sdl.Window
	GlCtx          sdl.GLContext
	EventCallbacks []func(sdl.Event)
	Rend           renderer.Render
}

func (w *Window) handleInputs() {

	input.EventLoopStart()
	imIO := imgui.CurrentIO()

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

		//Fire callbacks
		for i := 0; i < len(w.EventCallbacks); i++ {
			w.EventCallbacks[i](event)
		}

		//Internal processing
		switch e := event.(type) {

		case *sdl.MouseWheelEvent:

			input.HandleMouseWheelEvent(e)

			xDelta, yDelta := input.GetMouseWheelMotion()
			imIO.AddMouseWheelDelta(float32(xDelta), float32(yDelta))

		case *sdl.KeyboardEvent:
			input.HandleKeyboardEvent(e)

			if e.Type == sdl.KEYDOWN {
				imIO.KeyPress(int(e.Keysym.Scancode))
			} else if e.Type == sdl.KEYUP {
				imIO.KeyRelease(int(e.Keysym.Scancode))
			}

		case *sdl.TextInputEvent:
			imIO.AddInputCharacters(string(e.Text[:]))

		case *sdl.MouseButtonEvent:
			input.HandleMouseBtnEvent(e)

		case *sdl.MouseMotionEvent:
			input.HandleMouseMotionEvent(e)

		case *sdl.WindowEvent:
			if e.Event == sdl.WINDOWEVENT_SIZE_CHANGED {
				w.handleWindowResize()
			}

		case *sdl.QuitEvent:
			input.HandleQuitEvent(e)
		}
	}

	// If a mouse press event came, always pass it as "mouse held this frame", so we don't miss click-release events that are shorter than 1 frame.
	x, y, _ := sdl.GetMouseState()
	imIO.SetMousePosition(imgui.Vec2{X: float32(x), Y: float32(y)})

	imIO.SetMouseButtonDown(0, input.MouseDown(sdl.BUTTON_LEFT))
	imIO.SetMouseButtonDown(1, input.MouseDown(sdl.BUTTON_RIGHT))
	imIO.SetMouseButtonDown(2, input.MouseDown(sdl.BUTTON_MIDDLE))

	imIO.KeyShift(sdl.SCANCODE_LSHIFT, sdl.SCANCODE_RSHIFT)
	imIO.KeyCtrl(sdl.SCANCODE_LCTRL, sdl.SCANCODE_RCTRL)
	imIO.KeyAlt(sdl.SCANCODE_LALT, sdl.SCANCODE_RALT)
}

func (w *Window) handleWindowResize() {

	fbWidth, fbHeight := w.SDLWin.GLGetDrawableSize()
	if fbWidth <= 0 || fbHeight <= 0 {
		return
	}
	gl.Viewport(0, 0, fbWidth, fbHeight)
}

func (w *Window) Destroy() error {
	return w.SDLWin.Destroy()
}

func Init() error {

	isInited = true

	runtime.LockOSThread()
	timing.Init()
	err := initSDL()

	return err
}

func initSDL() error {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return err
	}

	sdl.ShowCursor(1)

	sdl.GLSetAttribute(sdl.MAJOR_VERSION, 4)
	sdl.GLSetAttribute(sdl.MINOR_VERSION, 1)

	// R(0-255) G(0-255) B(0-255)
	sdl.GLSetAttribute(sdl.GL_RED_SIZE, 8)
	sdl.GLSetAttribute(sdl.GL_GREEN_SIZE, 8)
	sdl.GLSetAttribute(sdl.GL_BLUE_SIZE, 8)

	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	sdl.GLSetAttribute(sdl.GL_DEPTH_SIZE, 24)
	sdl.GLSetAttribute(sdl.GL_STENCIL_SIZE, 8)

	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)

	return nil
}

func CreateOpenGLWindow(title string, x, y, width, height int32, flags WindowFlags, rend renderer.Render) (*Window, error) {
	return createWindow(title, x, y, width, height, WindowFlags_OPENGL|flags, rend)
}

func CreateOpenGLWindowCentered(title string, width, height int32, flags WindowFlags, rend renderer.Render) (*Window, error) {
	return createWindow(title, -1, -1, width, height, WindowFlags_OPENGL|flags, rend)
}

func createWindow(title string, x, y, width, height int32, flags WindowFlags, rend renderer.Render) (*Window, error) {

	assert.T(isInited, "engine.Init was not called!")
	if x == -1 && y == -1 {
		x = sdl.WINDOWPOS_CENTERED
		y = sdl.WINDOWPOS_CENTERED
	}

	sdlWin, err := sdl.CreateWindow(title, x, y, width, height, uint32(flags))
	if err != nil {
		return nil, err
	}
	win := &Window{
		SDLWin:         sdlWin,
		EventCallbacks: make([]func(sdl.Event), 0),
		Rend:           rend,
	}

	win.GlCtx, err = sdlWin.GLCreateContext()
	if err != nil {
		return nil, err
	}

	err = initOpenGL()
	if err != nil {
		return nil, err
	}

	return win, err
}

func initOpenGL() error {

	if err := gl.Init(); err != nil {
		return err
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.ClearColor(0, 0, 0, 1)
	return nil
}

func SetVSync(enabled bool) {
	assert.T(isInited, "engine.Init was not called!")

	if enabled {
		sdl.GLSetSwapInterval(1)
	} else {
		sdl.GLSetSwapInterval(0)
	}
}
