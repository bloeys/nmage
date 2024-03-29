package engine

import (
	"runtime"

	imgui "github.com/AllenDang/cimgui-go"
	"github.com/bloeys/nmage/assert"
	"github.com/bloeys/nmage/input"
	"github.com/bloeys/nmage/renderer"
	"github.com/bloeys/nmage/timing"
	nmageimgui "github.com/bloeys/nmage/ui/imgui"
	"github.com/go-gl/gl/v4.1-core/gl"
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
	imIo := imgui.CurrentIO()

	// @TODO: Would be nice to have imgui package process its own events via a callback instead of it being part of engine code
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
			imIo.AddMouseWheelDelta(float32(xDelta), float32(yDelta))

		case *sdl.KeyboardEvent:

			input.HandleKeyboardEvent(e)
			imIo.AddKeyEvent(nmageimgui.SdlScancodeToImGuiKey(e.Keysym.Scancode), e.Type == sdl.KEYDOWN)

			// Send modifier key updates to imgui
			if e.Keysym.Sym == sdl.K_LCTRL || e.Keysym.Sym == sdl.K_RCTRL {
				imIo.SetKeyCtrl(e.Type == sdl.KEYDOWN)
			}

			if e.Keysym.Sym == sdl.K_LSHIFT || e.Keysym.Sym == sdl.K_RSHIFT {
				imIo.SetKeyShift(e.Type == sdl.KEYDOWN)
			}

			if e.Keysym.Sym == sdl.K_LALT || e.Keysym.Sym == sdl.K_RALT {
				imIo.SetKeyAlt(e.Type == sdl.KEYDOWN)
			}

			if e.Keysym.Sym == sdl.K_LGUI || e.Keysym.Sym == sdl.K_RGUI {
				imIo.SetKeySuper(e.Type == sdl.KEYDOWN)
			}

		case *sdl.TextInputEvent:
			imIo.AddInputCharactersUTF8(e.GetText())

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
	imIo.SetMousePos(imgui.Vec2{X: float32(x), Y: float32(y)})

	imIo.SetMouseButtonDown(0, input.MouseDown(sdl.BUTTON_LEFT))
	imIo.SetMouseButtonDown(1, input.MouseDown(sdl.BUTTON_RIGHT))
	imIo.SetMouseButtonDown(2, input.MouseDown(sdl.BUTTON_MIDDLE))
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

	err := sdl.Init(sdl.INIT_TIMER | sdl.INIT_VIDEO)
	if err != nil {
		return err
	}

	sdl.ShowCursor(1)

	sdl.GLSetAttribute(sdl.MAJOR_VERSION, 4)
	sdl.GLSetAttribute(sdl.MINOR_VERSION, 1)

	sdl.GLSetAttribute(sdl.GL_RED_SIZE, 8)
	sdl.GLSetAttribute(sdl.GL_GREEN_SIZE, 8)
	sdl.GLSetAttribute(sdl.GL_BLUE_SIZE, 8)
	sdl.GLSetAttribute(sdl.GL_ALPHA_SIZE, 8)

	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	sdl.GLSetAttribute(sdl.GL_DEPTH_SIZE, 24)
	sdl.GLSetAttribute(sdl.GL_STENCIL_SIZE, 8)

	sdl.GLSetAttribute(sdl.GL_FRAMEBUFFER_SRGB_CAPABLE, 1)

	// Allows us to do MSAA
	sdl.GLSetAttribute(sdl.GL_MULTISAMPLEBUFFERS, 1)
	sdl.GLSetAttribute(sdl.GL_MULTISAMPLESAMPLES, 4)

	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)

	return nil
}

func CreateOpenGLWindow(title string, x, y, width, height int32, flags WindowFlags, rend renderer.Render) (*Window, error) {
	return createWindow(title, x, y, width, height, WindowFlags_OPENGL|flags, rend)
}

func CreateOpenGLWindowCentered(title string, width, height int32, flags WindowFlags, rend renderer.Render) (*Window, error) {
	return createWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, width, height, WindowFlags_OPENGL|flags, rend)
}

func createWindow(title string, x, y, width, height int32, flags WindowFlags, rend renderer.Render) (*Window, error) {

	assert.T(isInited, "engine.Init() was not called!")

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
	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.FRAMEBUFFER_SRGB)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.ClearColor(0, 0, 0, 1)
	return nil
}

func SetSrgbFramebuffer(isEnabled bool) {

	if isEnabled {
		gl.Enable(gl.FRAMEBUFFER_SRGB)
	} else {
		gl.Disable(gl.FRAMEBUFFER_SRGB)
	}
}

func SetVSync(enabled bool) {
	assert.T(isInited, "engine.Init was not called!")

	if enabled {
		sdl.GLSetSwapInterval(1)
	} else {
		sdl.GLSetSwapInterval(0)
	}
}

func SetMSAA(isEnabled bool) {

	if isEnabled {
		gl.Enable(gl.MULTISAMPLE)
	} else {
		gl.Disable(gl.MULTISAMPLE)
	}
}
