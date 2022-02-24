package nmageimgui

import (
	"github.com/bloeys/gglm/gglm"
	"github.com/bloeys/nmage/asserts"
	"github.com/bloeys/nmage/materials"
	"github.com/bloeys/nmage/timing"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/inkyblackness/imgui-go/v4"
	"github.com/veandco/go-sdl2/sdl"
)

type ImguiInfo struct {
	ImCtx *imgui.Context

	Mat        *materials.Material
	vaoID      uint32
	vboID      uint32
	indexBufID uint32
	texID      uint32
}

func (i *ImguiInfo) FrameStart(winWidth, winHeight float32) {

	if err := i.ImCtx.SetCurrent(); err != nil {
		asserts.T(false, "Setting imgui ctx as current failed. Err: "+err.Error())
	}

	imIO := imgui.CurrentIO()
	imIO.SetDisplaySize(imgui.Vec2{X: float32(winWidth), Y: float32(winHeight)})

	// Setup time step (we don't use SDL_GetTicks() because it is using millisecond resolution)
	frequency := sdl.GetPerformanceFrequency()
	currentTime := sdl.GetPerformanceCounter()
	if timing.ElapsedTime() > 0 {
		imIO.SetDeltaTime(float32(currentTime-timing.ElapsedTime()) / float32(frequency))
	} else {
		imIO.SetDeltaTime(1.0 / 60.0)
	}

	imgui.NewFrame()
}

func (i *ImguiInfo) Render(winWidth, winHeight float32, fbWidth, fbHeight int32) {

	if err := i.ImCtx.SetCurrent(); err != nil {
		asserts.T(false, "Setting imgui ctx as current failed. Err: "+err.Error())
	}

	imgui.Render()

	// Avoid rendering when minimized, scale coordinates for retina displays (screen coordinates != framebuffer coordinates)
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

	i.Mat.Bind()

	gl.Uniform1i(gl.GetUniformLocation(i.Mat.ShaderProg.ID, gl.Str("Texture\x00")), 0)

	//PERF: only update the ortho matrix on window resize
	orthoMat := gglm.Ortho(0, float32(winWidth), 0, float32(winHeight), 0, 20)
	i.Mat.SetUnifMat4("ProjMtx", &orthoMat.Mat4)
	gl.BindSampler(0, 0) // Rely on combined texture/sampler state.

	// Recreate the VAO every time
	// (This is to easily allow multiple GL contexts. VAO are not shared among GL contexts, and
	// we don't track creation/deletion of windows so we don't have an obvious key to use to cache them.)
	gl.BindVertexArray(i.vaoID)
	gl.BindBuffer(gl.ARRAY_BUFFER, i.vboID)

	vertexSize, vertexOffsetPos, vertexOffsetUv, vertexOffsetCol := imgui.VertexBufferLayout()
	i.Mat.EnableAttribute("Position")
	i.Mat.EnableAttribute("UV")
	i.Mat.EnableAttribute("Color")
	gl.VertexAttribPointerWithOffset(uint32(i.Mat.GetAttribLoc("Position")), 2, gl.FLOAT, false, int32(vertexSize), uintptr(vertexOffsetPos))
	gl.VertexAttribPointerWithOffset(uint32(i.Mat.GetAttribLoc("UV")), 2, gl.FLOAT, false, int32(vertexSize), uintptr(vertexOffsetUv))
	gl.VertexAttribPointerWithOffset(uint32(i.Mat.GetAttribLoc("Color")), 4, gl.UNSIGNED_BYTE, true, int32(vertexSize), uintptr(vertexOffsetCol))

	indexSize := imgui.IndexBufferLayout()
	drawType := gl.UNSIGNED_SHORT
	if indexSize == 4 {
		drawType = gl.UNSIGNED_INT
	}

	// Draw
	for _, list := range drawData.CommandLists() {

		vertexBuffer, vertexBufferSize := list.VertexBuffer()
		gl.BindBuffer(gl.ARRAY_BUFFER, i.vboID)
		gl.BufferData(gl.ARRAY_BUFFER, vertexBufferSize, vertexBuffer, gl.STREAM_DRAW)

		indexBuffer, indexBufferSize := list.IndexBuffer()
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, i.indexBufID)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indexBufferSize, indexBuffer, gl.STREAM_DRAW)

		for _, cmd := range list.Commands() {
			if cmd.HasUserCallback() {
				cmd.CallUserCallback(list)
			} else {

				gl.BindTexture(gl.TEXTURE_2D, i.texID)
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

func (i *ImguiInfo) AddFontTTF(fontPath string, fontSize float32) imgui.Font {

	imIO := imgui.CurrentIO()

	a := imIO.Fonts()
	f := a.AddFontFromFileTTF(fontPath, fontSize)
	image := a.TextureDataAlpha8()

	gl.BindTexture(gl.TEXTURE_2D, i.texID)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(image.Width), int32(image.Height), 0, gl.RED, gl.UNSIGNED_BYTE, image.Pixels)

	return f
}

func NewImGUI() ImguiInfo {

	imguiInfo := ImguiInfo{
		ImCtx: imgui.CreateContext(nil),
		Mat:   materials.NewMaterial("ImGUI Mat", "./res/shaders/imgui"),
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
	imguiInfo.Mat.Bind()
	imguiInfo.Mat.EnableAttribute("Position")
	imguiInfo.Mat.EnableAttribute("UV")
	imguiInfo.Mat.EnableAttribute("Color")
	imguiInfo.Mat.UnBind()

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

	return imguiInfo
}
