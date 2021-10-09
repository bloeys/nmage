package main

import (
	"fmt"
	"math"
	"time"

	"github.com/bloeys/go-sdl-engine/input"
	"github.com/bloeys/go-sdl-engine/logging"
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
		logging.ErrLog.Panicln("Failed to init SDL. Err:", err.Error())
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
		logging.ErrLog.Panicln("Failed to create window. Err: " + err.Error())
	}
	defer window.Destroy()

	//Create GL context
	glContext, err = window.GLCreateContext()
	if err != nil {
		logging.ErrLog.Panicln("Creating OpenGL context failed. Err: " + err.Error())
	}
	defer sdl.GLDeleteContext(glContext)

	if err := gl.Init(); err != nil {
		logging.ErrLog.Panicln("Initing OpenGL Context failed. Err: " + err.Error())
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

var simpleProg shaders.Program

func loadShaders() {

	simpleVert, err := shaders.NewShaderFromFile("simpleVert", "./res/shaders/simple.vert.glsl", shaders.Vertex)
	panicIfErr(err, "Parsing vert shader failed")

	simpleFrag, err := shaders.NewShaderFromFile("simpleFrag", "./res/shaders/simple.frag.glsl", shaders.Fragment)
	panicIfErr(err, "Parsing frag shader failed")

	simpleProg = shaders.NewProgram("simple")
	simpleProg.AttachShader(simpleVert)
	simpleProg.AttachShader(simpleFrag)
	simpleProg.Link()
}

var vao uint32

func gameLoop() {

	//vertex positions in opengl coords
	verts := []float32{
		-0.5, 0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
		-0.5, -0.5, 0,
	}

	//Trianlge indices used for drawing
	indices := []uint32{
		0, 1, 2,
		0, 2, 3,
	}

	//Create a VAO to store the different VBOs of a given object/set of vertices
	gl.GenVertexArrays(1, &vao)

	//Bind the VAO first so later buffer binds/VBOs are put within this VAO
	gl.BindVertexArray(vao)

	//Gen buffer to hold EBOs and fill it with data. Note that an EBO must NOT be unbound before the VAO, otherwise
	//its settings are lost as the VAO records its actions
	var ebo uint32
	gl.GenBuffers(1, &ebo)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	//Gen vertPos VBO and fill it with data
	var vertPosVBO uint32
	gl.GenBuffers(1, &vertPosVBO)

	gl.BindBuffer(gl.ARRAY_BUFFER, vertPosVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)

	//Assign vertPos VBO to vertPos shader attribute by specifying that each vertPos variable
	//takes 3 values from the VBO, where each value is a float.

	//We also specify the total size (in bytes) of the values used for a single vertPos.
	//The offset defines the bytes to skip between each set of vertPos values
	vertPosLoc := uint32(gl.GetAttribLocation(simpleProg.ID, gl.Str("vertPos\x00")))
	gl.VertexAttribPointer(vertPosLoc, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	//Vertex attributes are disabled by default, so we need to finally enable it
	gl.EnableVertexAttribArray(vertPosLoc)

	//We are done working with VBOs so can unbind the VAO to avoid corrupting it later.
	//Note: Actions (binding/setting buffers, enabling/disabling attribs) done between
	//bind and unbind of a VAO are recorded by it, and when its rebinded before a draw the
	//settings are retrieved, therefore keep in mind work after a VAO unbind will be lost.
	gl.BindVertexArray(0)

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

	simpleProg.Use()

	deg2rad := math.Pi / 180.0
	t := float64(time.Now().UnixMilli()) / 10

	x := float32(math.Sin(t*deg2rad*0.3)+1) * 0.5
	y := float32(math.Sin(t*deg2rad*0.5)+1) * 0.5
	z := float32(math.Sin(t*deg2rad*0.7)+1) * 0.5
	simpleProg.SetUniformF32("c", x, y, z)

	gl.BindVertexArray(vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}

func panicIfErr(err error, msg string) {

	if err == nil {
		return
	}

	logging.ErrLog.Panicln(msg+". Err:", err.Error())
}
