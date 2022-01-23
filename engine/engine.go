package engine

import (
	"github.com/bloeys/nmage/timing"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

type Window struct {
	SDLWin *sdl.Window
	GlCtx  sdl.GLContext
}

func (w *Window) Destroy() error {
	return w.SDLWin.Destroy()
}

func Init() error {

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

func CreateOpenGLWindow(title string, x, y, width, height int32, flags WindowFlags) (*Window, error) {
	return createWindow(title, x, y, width, height, WindowFlags_OPENGL|flags)
}

func CreateOpenGLWindowCentered(title string, width, height int32, flags WindowFlags) (*Window, error) {
	return createWindow(title, -1, -1, width, height, WindowFlags_OPENGL|flags)
}

func createWindow(title string, x, y, width, height int32, flags WindowFlags) (*Window, error) {

	if x == -1 && y == -1 {
		x = sdl.WINDOWPOS_CENTERED
		y = sdl.WINDOWPOS_CENTERED
	}

	sdlWin, err := sdl.CreateWindow(title, x, y, width, height, uint32(flags))
	if err != nil {
		return nil, err
	}
	win := &Window{SDLWin: sdlWin}

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

	gl.ClearColor(0, 0, 0, 1)
	return nil
}

func SetVSync(enabled bool) {
	if enabled {
		sdl.GLSetSwapInterval(1)
	} else {
		sdl.GLSetSwapInterval(0)
	}
}
