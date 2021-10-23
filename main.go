package main

import (
	"github.com/bloeys/go-sdl-engine/input"
	"github.com/bloeys/go-sdl-engine/logging"
	"github.com/bloeys/go-sdl-engine/shaders"
	"github.com/go-gl/gl/v4.6-compatibility/gl"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth  = 1280
	winHeight = 720
)

var (
	isRunning bool = true

	window *sdl.Window
)

func main() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		logging.ErrLog.Fatalln("Failed to init SDL. Err:", err)
	}

	window, err = sdl.CreateWindow("Go SDL Engine", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, winWidth, winHeight, sdl.WINDOW_OPENGL)
	if err != nil {
		logging.ErrLog.Fatalln("Failed to create window. Err: ", err)
	}
	defer window.Destroy()

	glCtx, err := window.GLCreateContext()
	if err != nil {
		logging.ErrLog.Fatalln("Failed to create OpenGL context. Err: ", err)
	}
	defer sdl.GLDeleteContext(glCtx)

	initOpenGL()
	loadShaders()

	//Game loop
	for isRunning {

		handleInputs()
		runGameLogic()
		draw()

		sdl.Delay(17)
	}
}

func initOpenGL() {

	if err := gl.Init(); err != nil {
		logging.ErrLog.Fatalln(err)
	}

	sdl.GLSetAttribute(sdl.MAJOR_VERSION, 4)
	sdl.GLSetAttribute(sdl.MINOR_VERSION, 6)

	// R(0-255) G(0-255) B(0-255)
	sdl.GLSetAttribute(sdl.GL_RED_SIZE, 8)
	sdl.GLSetAttribute(sdl.GL_GREEN_SIZE, 8)
	sdl.GLSetAttribute(sdl.GL_BLUE_SIZE, 8)

	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	gl.ClearColor(0, 0, 0, 1)

	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_COMPATIBILITY)
	// sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
}

func loadShaders() {

	simpleShader, err := shaders.NewShaderProgram()
	if err != nil {
		logging.ErrLog.Fatalln("Failed to create new shader program. Err: ", err)
	}

	vertShader, err := shaders.LoadAndCompilerShader("./res/shaders/simple.vert.glsl", shaders.VertexShaderType)
	if err != nil {
		logging.ErrLog.Fatalln("Failed to create new shader. Err: ", err)
	}

	fragShader, err := shaders.LoadAndCompilerShader("./res/shaders/simple.frag.glsl", shaders.FragmentShaderType)
	if err != nil {
		logging.ErrLog.Fatalln("Failed to create new shader. Err: ", err)
	}

	simpleShader.AttachShader(vertShader)
	simpleShader.AttachShader(fragShader)
	simpleShader.Link()
}

func handleInputs() {

	input.EventLoopStart()

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

		switch e := event.(type) {

		case *sdl.KeyboardEvent:
			input.HandleKeyboardEvent(e)
		case *sdl.MouseButtonEvent:
			input.HandleMouseEvent(e)
		case *sdl.QuitEvent:
			isRunning = false
		}
	}
}

func runGameLogic() {

}

func draw() {

	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.Begin(gl.TRIANGLES)

	gl.Vertex3f(-0.5, 0.5, 0)
	gl.Vertex3f(0.5, 0.5, 0)
	gl.Vertex3f(-0.5, -0.5, 0)

	gl.Vertex3f(0.5, 0.5, 0)
	gl.Vertex3f(0.5, -0.5, 0)
	gl.Vertex3f(-0.5, -0.5, 0)

	gl.End()

	window.GLSwap()
}
