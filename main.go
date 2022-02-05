package main

import (
	"fmt"
	"runtime"

	"github.com/bloeys/assimp-go/asig"
	"github.com/bloeys/gglm/gglm"
	"github.com/bloeys/nmage/engine"
	"github.com/bloeys/nmage/input"
	"github.com/bloeys/nmage/logging"
	"github.com/bloeys/nmage/materials"
	"github.com/bloeys/nmage/meshes"
	"github.com/bloeys/nmage/timing"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/inkyblackness/imgui-go/v4"
	"github.com/veandco/go-sdl2/sdl"
)

//TODO: Tasks:
//Engine loop
//Proper rendering setup
//Flesh out the material system
//Object
//Abstract UI
//Textures
//Audio

//Low Priority:
//	Proper Asset loading

type ImguiInfo struct {
	imCtx *imgui.Context

	vaoID      uint32
	vboID      uint32
	indexBufID uint32
	texID      uint32
}

var (
	winWidth  int32 = 1280
	winHeight int32 = 720

	isRunning bool = true
	window    *engine.Window

	simpleMat *materials.Material
	imguiMat  *materials.Material
	cubeMesh  *meshes.Mesh

	modelMat = gglm.NewTrMatId()
	projMat  = &gglm.Mat4{}

	imguiInfo *ImguiInfo

	camPos     = gglm.NewVec3(0, 0, -10)
	camForward = gglm.NewVec3(0, 0, 1)
)

func main() {

	runtime.LockOSThread()

	//Init engine
	err := engine.Init()
	if err != nil {
		logging.ErrLog.Fatalln("Failed to init nMage. Err:", err)
	}

	//Create window
	window, err = engine.CreateOpenGLWindowCentered("nMage", winWidth, winHeight, engine.WindowFlags_RESIZABLE)
	if err != nil {
		logging.ErrLog.Fatalln("Failed to create window. Err: ", err)
	}
	defer window.Destroy()

	engine.SetVSync(false)

	//Create materials
	simpleMat = materials.NewMaterial("Simple Mat", "./res/shaders/simple")
	imguiMat = materials.NewMaterial("ImGUI Mat", "./res/shaders/imgui")

	//Load meshes
	cubeMesh, err = meshes.NewMesh("Cube", "./res/models/color-cube.fbx", asig.PostProcess(0))
	if err != nil {
		logging.ErrLog.Fatalln("Failed to load cube mesh. Err: ", err)
	}

	initImGUI()

	//Enable vertex attributes
	simpleMat.SetAttribute(cubeMesh.Buf)

	//Movement, scale and rotation
	translationMat := gglm.NewTranslationMat(gglm.NewVec3(0, 0, 0))
	scaleMat := gglm.NewScaleMat(gglm.NewVec3(0.25, 0.25, 0.25))
	rotMat := gglm.NewRotMat(gglm.NewQuatEuler(gglm.NewVec3(0, 0, 0).AsRad()))

	modelMat.Mul(translationMat.Mul(rotMat.Mul(scaleMat)))
	simpleMat.SetUnifMat4("modelMat", &modelMat.Mat4)

	//Moves objects into the cameras view
	updateViewMat()

	//Perspective/Depth
	projMat := gglm.Perspective(45*gglm.Deg2Rad, float32(winWidth)/float32(winHeight), 0.1, 500)
	simpleMat.SetUnifMat4("projMat", projMat)

	//Lights
	simpleMat.SetUnifVec3("lightPos1", &lightPos1)
	simpleMat.SetUnifVec3("lightColor1", &lightColor1)

	//Game loop
	for isRunning {

		timing.FrameStarted()

		handleInputs()
		runGameLogic()

		draw()

		timing.FrameEnded()

		window.SDLWin.SetTitle(fmt.Sprintf("FPS: %.0f; Elapsed: %v", 1/timing.DT(), timing.ElapsedTime()))
	}
}

func updateViewMat() {
	targetPos := camPos.Clone().Add(camForward)
	viewMat := gglm.LookAt(camPos, targetPos, gglm.NewVec3(0, 1, 0))
	simpleMat.SetUnifMat4("viewMat", &viewMat.Mat4)
}

func initImGUI() {

	imguiInfo = &ImguiInfo{
		imCtx: imgui.CreateContext(nil),
	}

	imIO := imgui.CurrentIO()
	imIO.SetBackendFlags(imIO.GetBackendFlags() | imgui.BackendFlagsRendererHasVtxOffset)

	gl.GenVertexArrays(1, &imguiInfo.vaoID)
	gl.GenBuffers(1, &imguiInfo.vboID)
	gl.GenBuffers(1, &imguiInfo.indexBufID)
	gl.GenTextures(1, &imguiInfo.texID)

	// Upload font to gpu
	gl.BindTexture(gl.TEXTURE_2D, imguiInfo.texID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)

	image := imIO.Fonts().TextureDataAlpha8()
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(image.Width), int32(image.Height), 0, gl.RED, gl.UNSIGNED_BYTE, image.Pixels)

	// Store our identifier
	imIO.Fonts().SetTextureID(imgui.TextureID(imguiInfo.texID))

	//Shader attributes
	imguiMat.Bind()
	imguiMat.EnableAttribute("Position")
	imguiMat.EnableAttribute("UV")
	imguiMat.EnableAttribute("Color")
	imguiMat.UnBind()

	//Init imgui input mapping
	keys := map[int]int{
		imgui.KeyTab:        sdl.SCANCODE_TAB,
		imgui.KeyLeftArrow:  sdl.SCANCODE_LEFT,
		imgui.KeyRightArrow: sdl.SCANCODE_RIGHT,
		imgui.KeyUpArrow:    sdl.SCANCODE_UP,
		imgui.KeyDownArrow:  sdl.SCANCODE_DOWN,
		imgui.KeyPageUp:     sdl.SCANCODE_PAGEUP,
		imgui.KeyPageDown:   sdl.SCANCODE_PAGEDOWN,
		imgui.KeyHome:       sdl.SCANCODE_HOME,
		imgui.KeyEnd:        sdl.SCANCODE_END,
		imgui.KeyInsert:     sdl.SCANCODE_INSERT,
		imgui.KeyDelete:     sdl.SCANCODE_DELETE,
		imgui.KeyBackspace:  sdl.SCANCODE_BACKSPACE,
		imgui.KeySpace:      sdl.SCANCODE_BACKSPACE,
		imgui.KeyEnter:      sdl.SCANCODE_RETURN,
		imgui.KeyEscape:     sdl.SCANCODE_ESCAPE,
		imgui.KeyA:          sdl.SCANCODE_A,
		imgui.KeyC:          sdl.SCANCODE_C,
		imgui.KeyV:          sdl.SCANCODE_V,
		imgui.KeyX:          sdl.SCANCODE_X,
		imgui.KeyY:          sdl.SCANCODE_Y,
		imgui.KeyZ:          sdl.SCANCODE_Z,
	}

	// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.
	for imguiKey, nativeKey := range keys {
		imIO.KeyMap(imguiKey, nativeKey)
	}
}

func handleInputs() {

	input.EventLoopStart()
	imIO := imgui.CurrentIO()

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

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

		case *sdl.QuitEvent:
			isRunning = false
		}
	}

	currWinWidth, currWinHeight := window.SDLWin.GetSize()
	if winWidth != currWinWidth || winHeight != currWinHeight {
		handleWindowResize(currWinWidth, currWinHeight)
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

func handleWindowResize(newWinWidth, newWinHeight int32) {

	winWidth = newWinWidth
	winHeight = newWinHeight

	fbWidth, fbHeight := window.SDLWin.GLGetDrawableSize()
	if fbWidth <= 0 || fbHeight <= 0 {
		return
	}
	gl.Viewport(0, 0, fbWidth, fbHeight)

	projMat = gglm.Perspective(45*gglm.Deg2Rad, float32(winWidth)/float32(winHeight), 0.1, 20)
	simpleMat.SetUnifMat4("projMat", projMat)
}

var time uint64 = 0
var name string = ""

var ambientColor gglm.Vec3 = *gglm.NewVec3(1, 1, 1)
var ambientColorStrength float32 = 0.1

var lightPos1 gglm.Vec3 = *gglm.NewVec3(2, 2, 0)
var lightColor1 gglm.Vec3 = *gglm.NewVec3(1, 1, 1)

func runGameLogic() {

	var camSpeed float32 = 15
	if input.KeyDown(sdl.K_w) {
		camPos.Data[1] += camSpeed * timing.DT()
		updateViewMat()
	}
	if input.KeyDown(sdl.K_s) {
		camPos.Data[1] -= camSpeed * timing.DT()
		updateViewMat()
	}
	if input.KeyDown(sdl.K_d) {
		camPos.Data[0] += camSpeed * timing.DT()
		updateViewMat()
	}
	if input.KeyDown(sdl.K_a) {
		camPos.Data[0] -= camSpeed * timing.DT()
		updateViewMat()
	}

	if input.GetMouseWheelYNorm() > 0 {
		camPos.Data[2] += 1
		updateViewMat()
	} else if input.GetMouseWheelYNorm() < 0 {
		camPos.Data[2] -= 1
		updateViewMat()
	}

	modelMat.Rotate(10*timing.DT()*gglm.Deg2Rad, gglm.NewVec3(1, 1, 1).Normalize())
	simpleMat.SetUnifMat4("modelMat", &modelMat.Mat4)

	//ImGUI
	imIO := imgui.CurrentIO()
	imIO.SetDisplaySize(imgui.Vec2{X: float32(winWidth), Y: float32(winHeight)})

	// Setup time step (we don't use SDL_GetTicks() because it is using millisecond resolution)
	frequency := sdl.GetPerformanceFrequency()
	currentTime := sdl.GetPerformanceCounter()
	if time > 0 {
		imIO.SetDeltaTime(float32(currentTime-time) / float32(frequency))
	} else {
		imIO.SetDeltaTime(1.0 / 60.0)
	}
	time = currentTime

	imgui.NewFrame()

	if imgui.SliderFloat3("Ambient Color", &ambientColor.Data, 0, 1) {
		simpleMat.SetUnifVec3("ambientLightColor", &ambientColor)
	}

	if imgui.SliderFloat("Ambient Color Strength", &ambientColorStrength, 0, 1) {
		simpleMat.SetUnifFloat32("ambientStrength", ambientColorStrength)
	}

	if imgui.SliderFloat3("Light Pos 1", &lightPos1.Data, -10, 10) {
		simpleMat.SetUnifVec3("lightPos1", &lightPos1)
	}

	if imgui.SliderFloat3("Light Color 1", &lightColor1.Data, 0, 1) {
		simpleMat.SetUnifVec3("lightColor1", &lightColor1)
	}

	imgui.Render()
}

func draw() {

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	simpleMat.Bind()
	cubeMesh.Buf.Bind()
	tempModelMat := modelMat.Clone()

	rowSize := 10
	for y := 0; y < rowSize; y++ {
		for x := 0; x < rowSize; x++ {
			simpleMat.SetUnifMat4("modelMat", &tempModelMat.Translate(gglm.NewVec3(-1, 0, 0)).Mat4)
			gl.DrawElements(gl.TRIANGLES, cubeMesh.Buf.IndexBufCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
		}
		simpleMat.SetUnifMat4("modelMat", &tempModelMat.Translate(gglm.NewVec3(float32(rowSize), -1, 0)).Mat4)
	}
	simpleMat.SetUnifMat4("modelMat", &modelMat.Mat4)

	drawUI()
	window.SDLWin.GLSwap()
}

func drawUI() {

	// Avoid rendering when minimized, scale coordinates for retina displays (screen coordinates != framebuffer coordinates)
	fbWidth, fbHeight := window.SDLWin.GLGetDrawableSize()
	if fbWidth <= 0 || fbHeight <= 0 {
		return
	}

	drawData := imgui.RenderedDrawData()
	drawData.ScaleClipRects(imgui.Vec2{
		X: float32(fbWidth) / float32(winWidth),
		Y: float32(fbHeight) / float32(winHeight),
	})

	// Setup render state: alpha-blending enabled, no face culling, no depth testing, scissor enabled, polygon fill
	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.SCISSOR_TEST)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	// Setup viewport, orthographic projection matrix
	// Our visible imgui space lies from draw_data->DisplayPos (top left) to draw_data->DisplayPos+data_data->DisplaySize (bottom right).
	// DisplayMin is typically (0,0) for single viewport apps.

	imguiMat.Bind()

	gl.Uniform1i(gl.GetUniformLocation(imguiMat.ShaderProg.ID, gl.Str("Texture\x00")), 0)

	//PERF: only update the ortho matrix on window resize
	orthoMat := gglm.Ortho(0, float32(winWidth), 0, float32(winHeight), 0, 20)
	imguiMat.SetUnifMat4("ProjMtx", &orthoMat.Mat4)
	gl.BindSampler(0, 0) // Rely on combined texture/sampler state.

	// Recreate the VAO every time
	// (This is to easily allow multiple GL contexts. VAO are not shared among GL contexts, and
	// we don't track creation/deletion of windows so we don't have an obvious key to use to cache them.)
	gl.BindVertexArray(imguiInfo.vaoID)
	gl.BindBuffer(gl.ARRAY_BUFFER, imguiInfo.vboID)

	vertexSize, vertexOffsetPos, vertexOffsetUv, vertexOffsetCol := imgui.VertexBufferLayout()
	imguiMat.EnableAttribute("Position")
	imguiMat.EnableAttribute("UV")
	imguiMat.EnableAttribute("Color")
	gl.VertexAttribPointerWithOffset(uint32(imguiMat.GetAttribLoc("Position")), 2, gl.FLOAT, false, int32(vertexSize), uintptr(vertexOffsetPos))
	gl.VertexAttribPointerWithOffset(uint32(imguiMat.GetAttribLoc("UV")), 2, gl.FLOAT, false, int32(vertexSize), uintptr(vertexOffsetUv))
	gl.VertexAttribPointerWithOffset(uint32(imguiMat.GetAttribLoc("Color")), 4, gl.UNSIGNED_BYTE, true, int32(vertexSize), uintptr(vertexOffsetCol))

	indexSize := imgui.IndexBufferLayout()
	drawType := gl.UNSIGNED_SHORT
	if indexSize == 4 {
		drawType = gl.UNSIGNED_INT
	}

	// Draw
	for _, list := range drawData.CommandLists() {

		vertexBuffer, vertexBufferSize := list.VertexBuffer()
		gl.BindBuffer(gl.ARRAY_BUFFER, imguiInfo.vboID)
		gl.BufferData(gl.ARRAY_BUFFER, vertexBufferSize, vertexBuffer, gl.STREAM_DRAW)

		indexBuffer, indexBufferSize := list.IndexBuffer()
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, imguiInfo.indexBufID)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indexBufferSize, indexBuffer, gl.STREAM_DRAW)

		for _, cmd := range list.Commands() {
			if cmd.HasUserCallback() {
				cmd.CallUserCallback(list)
			} else {

				gl.BindTexture(gl.TEXTURE_2D, imguiInfo.texID)
				clipRect := cmd.ClipRect()
				gl.Scissor(int32(clipRect.X), int32(fbHeight)-int32(clipRect.W), int32(clipRect.Z-clipRect.X), int32(clipRect.W-clipRect.Y))

				gl.DrawElementsBaseVertex(gl.TRIANGLES, int32(cmd.ElementCount()), uint32(drawType), gl.PtrOffset(cmd.IndexOffset()*indexSize), int32(cmd.VertexOffset()))
			}
		}
	}

	//Reset gl state
	gl.Disable(gl.BLEND)
	gl.Disable(gl.SCISSOR_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
}
