package main

import (
	"fmt"

	"github.com/bloeys/go-sdl-engine/input"
	"github.com/bloeys/go-sdl-engine/shaders"
	"github.com/bloeys/go-sdl-engine/timing"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth  int32 = 1280
	winHeight int32 = 720
)

var (
	isRunning = true
	window    *sdl.Window
	glContext sdl.GLContext
)

func main() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic("Failed to init SDL. Err: " + err.Error())
	}
	defer sdl.Quit()

	//Size of each pixel field
	err = sdl.GLSetAttribute(sdl.GL_RED_SIZE, 8)
	panicIfErr(err, "")

	err = sdl.GLSetAttribute(sdl.GL_GREEN_SIZE, 8)
	panicIfErr(err, "")

	err = sdl.GLSetAttribute(sdl.GL_BLUE_SIZE, 8)
	panicIfErr(err, "")

	err = sdl.GLSetAttribute(sdl.GL_ALPHA_SIZE, 8)
	panicIfErr(err, "")

	//Min frame buffer size
	err = sdl.GLSetAttribute(sdl.GL_BUFFER_SIZE, 4*8)
	panicIfErr(err, "")

	//Whether to enable a double buffer
	err = sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	panicIfErr(err, "")

	//Run in compatiability (old and modern opengl) or modern (core) opengl only
	err = sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	panicIfErr(err, "")

	//Set wanted opengl version
	err = sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4)
	panicIfErr(err, "")

	err = sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 6)
	panicIfErr(err, "")

	//Create window
	window, err = sdl.CreateWindow(
		"Go Game Engine",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		winWidth,
		winHeight,
		sdl.WINDOW_OPENGL)
	if err != nil {
		panic("Failed to create window. Err: " + err.Error())
	}
	defer window.Destroy()

	//Create GL context
	glContext, err = window.GLCreateContext()
	if err != nil {
		panic("Creating OpenGL context failed. Err: " + err.Error())
	}
	defer sdl.GLDeleteContext(glContext)

	if err := gl.Init(); err != nil {
		panic("Initing OpenGL Context failed. Err: " + err.Error())
	}

	initGL()
	loadShaders()
	gameLoop()
}

func initGL() {

	gl.ClearColor(0, 0, 0, 0)

	gl.Enable(gl.DEPTH_TEST)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)
	gl.Viewport(0, 0, winWidth, winHeight)
}

func loadShaders() {

	simpleProg := shaders.NewProgram("simple")
	simpleVert, err := shaders.NewShaderFromFile("./res/shaders/simple.vert.glsl", shaders.Vertex)
	panicIfErr(err, "Parsing vert shader failed")
	simpleFrag, err := shaders.NewShaderFromFile("./res/shaders/simple.frag.glsl", shaders.Fragment)
	panicIfErr(err, "Parsing frag shader failed")

	simpleProg.AttachShader(simpleVert)
	simpleProg.AttachShader(simpleFrag)
}

func gameLoop() {

	for isRunning {

		timing.FrameStarted()

		handleEvents()
		update()
		draw()

		window.GLSwap()

		timing.FrameEnded()
		window.SetTitle(fmt.Sprintf("FPS: %.2f; dt: %.3f", timing.FPS(), timing.DT()))
	}
}

func handleEvents() {

	input.EventLoopStarted()

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

		switch e := event.(type) {

		case *sdl.KeyboardEvent:
			input.HandleKeyboardEvent(e)

		case *sdl.MouseButtonEvent:
			input.HandleMouseBtnEvent(e)

		case *sdl.QuitEvent:
			println("Quit at ", e.Timestamp)
			isRunning = false
		}
	}
}

func update() {
}

func draw() {
	//Clear screen and depth buffers
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func panicIfErr(err error, msg string) {

	if err == nil {
		return
	}

	panic(msg + "; Err: " + err.Error())
}
