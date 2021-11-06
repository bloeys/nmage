package main

import (
	"github.com/bloeys/gglm/gglm"
	"github.com/bloeys/go-sdl-engine/buffers"
	"github.com/bloeys/go-sdl-engine/input"
	"github.com/bloeys/go-sdl-engine/logging"
	"github.com/bloeys/go-sdl-engine/res/models"
	"github.com/bloeys/go-sdl-engine/shaders"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth  = 1280
	winHeight = 720
)

var (
	isRunning bool = true
	window    *sdl.Window

	simpleShader shaders.ShaderProgram
	bo           *buffers.BufferObject
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
	loadBuffers()

	simpleShader.SetAttribute("vertPos", bo, bo.VertPosBuf)
	simpleShader.EnableAttribute("vertPos")

	// simpleShader.SetAttribute("vertColor", bo, bo.ColorBuf)
	// simpleShader.EnableAttribute("vertColor")

	modelMat := gglm.NewTrMatId()
	translationMat := gglm.NewTranslationMat(gglm.NewVec3(-0.5, 0, 0))
	scaleMat := gglm.NewScaleMat(gglm.NewVec3(0.25, 0.25, 0.25))
	rotMat := gglm.NewRotMat(gglm.NewQuatEuler(gglm.NewVec3(0, 0, 0).AsRad()))

	modelMat.Mul(translationMat.Mul(rotMat.Mul(scaleMat)))
	simpleShader.SetUnifMat4("modelMat", &modelMat.Mat4)

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

	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
}

func loadShaders() {

	var err error
	simpleShader, err = shaders.NewShaderProgram()
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

func loadBuffers() {

	vertices := []float32{
		-0.5, 0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
		-0.5, -0.5, 0,
	}
	// colors := []float32{
	// 	1, 0, 0,
	// 	0, 0, 1,
	// 	0, 0, 1,
	// 	0, 0, 1,
	// }
	indices := []uint32{0, 1, 3, 1, 2, 3}

	//Load obj
	objInfo, err := models.LoadObj("./res/models/obj.obj")
	if err != nil {
		panic(err)
	}
	logging.InfoLog.Printf("%v", objInfo.TriIndices)

	vertices = objInfo.VertPos
	indices = objInfo.TriIndices

	bo = buffers.NewBufferObject()
	bo.GenBuffer(vertices, buffers.BufUsageStatic, buffers.BufTypeVertPos, buffers.DataTypeVec3)
	// bo.GenBuffer(colors, buffers.BufUsageStatic, buffers.BufTypeColor, buffers.DataTypeVec3)
	bo.GenBufferUint32(indices, buffers.BufUsageStatic, buffers.BufTypeIndex, buffers.DataTypeUint32)
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

	simpleShader.Activate()

	//DRAW
	bo.Activate()
	gl.DrawElements(gl.TRIANGLES, 36, gl.UNSIGNED_INT, gl.PtrOffset(0))
	bo.Deactivate()

	window.GLSwap()
}
