package nmageimgui

import (
	newimgui "github.com/AllenDang/cimgui-go"
	"github.com/bloeys/gglm/gglm"
	"github.com/bloeys/nmage/materials"
	"github.com/bloeys/nmage/timing"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

type ImguiInfo struct {
	ImCtx newimgui.Context
	// ImCtx2 *imgui.Context

	Mat        *materials.Material
	VaoID      uint32
	VboID      uint32
	IndexBufID uint32
	TexID      uint32
}

func (i *ImguiInfo) FrameStart(winWidth, winHeight float32) {

	// if err := i.ImCtx.SetCurrent(); err != nil {
	// 	assert.T(false, "Setting imgui ctx as current failed. Err: "+err.Error())
	// }

	imIO := newimgui.CurrentIO()
	imIO.SetDisplaySize(newimgui.Vec2{X: float32(winWidth), Y: float32(winHeight)})
	imIO.SetDeltaTime(timing.DT())

	newimgui.NewFrame()
}

func (i *ImguiInfo) Render(winWidth, winHeight float32, fbWidth, fbHeight int32) {

	// if err := i.ImCtx.SetCurrent(); err != nil {
	// 	assert.T(false, "Setting imgui ctx as current failed. Err: "+err.Error())
	// }

	newimgui.Render()

	// Avoid rendering when minimized, scale coordinates for retina displays (screen coordinates != framebuffer coordinates)
	if fbWidth <= 0 || fbHeight <= 0 {
		return
	}

	drawData := newimgui.CurrentDrawData()
	drawData.ScaleClipRects(newimgui.Vec2{
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
	i.Mat.SetUnifInt32("Texture", 0)

	//PERF: only update the ortho matrix on window resize
	orthoMat := gglm.Ortho(0, float32(winWidth), 0, float32(winHeight), 0, 20)
	i.Mat.SetUnifMat4("ProjMtx", &orthoMat.Mat4)
	gl.BindSampler(0, 0) // Rely on combined texture/sampler state.

	// Recreate the VAO every time
	// (This is to easily allow multiple GL contexts. VAO are not shared among GL contexts, and
	// we don't track creation/deletion of windows so we don't have an obvious key to use to cache them.)
	gl.BindVertexArray(i.VaoID)
	gl.BindBuffer(gl.ARRAY_BUFFER, i.VboID)

	vertexSize, vertexOffsetPos, vertexOffsetUv, vertexOffsetCol := newimgui.VertexBufferLayout()
	i.Mat.EnableAttribute("Position")
	i.Mat.EnableAttribute("UV")
	i.Mat.EnableAttribute("Color")
	gl.VertexAttribPointerWithOffset(uint32(i.Mat.GetAttribLoc("Position")), 2, gl.FLOAT, false, int32(vertexSize), uintptr(vertexOffsetPos))
	gl.VertexAttribPointerWithOffset(uint32(i.Mat.GetAttribLoc("UV")), 2, gl.FLOAT, false, int32(vertexSize), uintptr(vertexOffsetUv))
	gl.VertexAttribPointerWithOffset(uint32(i.Mat.GetAttribLoc("Color")), 4, gl.UNSIGNED_BYTE, true, int32(vertexSize), uintptr(vertexOffsetCol))

	indexSize := newimgui.IndexBufferLayout()
	drawType := gl.UNSIGNED_SHORT
	if indexSize == 4 {
		drawType = gl.UNSIGNED_INT
	}

	// Draw
	for _, list := range drawData.CommandLists() {

		vertexBuffer, vertexBufferSize := list.GetVertexBuffer()
		gl.BindBuffer(gl.ARRAY_BUFFER, i.VboID)
		gl.BufferData(gl.ARRAY_BUFFER, vertexBufferSize, vertexBuffer, gl.STREAM_DRAW)

		indexBuffer, indexBufferSize := list.GetIndexBuffer()
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, i.IndexBufID)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indexBufferSize, indexBuffer, gl.STREAM_DRAW)

		for _, cmd := range list.Commands() {
			if cmd.HasUserCallback() {
				cmd.CallUserCallback(list)
			} else {

				gl.BindTexture(gl.TEXTURE_2D, i.TexID)
				clipRect := cmd.ClipRect()
				gl.Scissor(int32(clipRect.X), int32(fbHeight)-int32(clipRect.W), int32(clipRect.Z-clipRect.X), int32(clipRect.W-clipRect.Y))

				gl.DrawElementsBaseVertex(gl.TRIANGLES, int32(cmd.ElemCount()), uint32(drawType), gl.PtrOffset(int(cmd.IdxOffset())*indexSize), int32(cmd.VtxOffset()))
			}
		}
	}

	//Reset gl state
	gl.Disable(gl.SCISSOR_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
}

func (i *ImguiInfo) AddFontTTF(fontPath string, fontSize float32, fontConfig *newimgui.FontConfig, glyphRanges *newimgui.GlyphRange) newimgui.Font {

	fontConfigToUse := newimgui.NewFontConfig()
	if fontConfig != nil {
		fontConfigToUse = *fontConfig
	}

	glyphRangesToUse := newimgui.NewGlyphRange()
	if glyphRanges != nil {
		glyphRangesToUse = *glyphRanges
	}

	imIO := newimgui.CurrentIO()

	a := imIO.Fonts()
	f := a.AddFontFromFileTTFV(fontPath, fontSize, fontConfigToUse, glyphRangesToUse.Data())
	pixels, width, height, _ := a.GetTextureDataAsAlpha8()

	gl.BindTexture(gl.TEXTURE_2D, i.TexID)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(width), int32(height), 0, gl.RED, gl.UNSIGNED_BYTE, pixels)

	return f
}

const imguiShdrSrc = `
//shader:vertex
#version 410

uniform mat4 ProjMtx;

in vec2 Position;
in vec2 UV;
in vec4 Color;

out vec2 Frag_UV;
out vec4 Frag_Color;

void main()
{
    Frag_UV = UV;
    Frag_Color = Color;
    gl_Position = ProjMtx * vec4(Position.xy, 0, 1);
}

//shader:fragment
#version 410

uniform sampler2D Texture;

in vec2 Frag_UV;
in vec4 Frag_Color;

out vec4 Out_Color;

void main()
{
    Out_Color = vec4(Frag_Color.rgb, Frag_Color.a * texture(Texture, Frag_UV.st).r);
}
`

func NewImGui() ImguiInfo {

	imguiInfo := ImguiInfo{
		// ImCtx2: imgui.CreateContext(nil),
		ImCtx: newimgui.CreateContext(),

		Mat: materials.NewMaterialSrc("ImGUI Mat", []byte(imguiShdrSrc)),
	}

	imIO := newimgui.CurrentIO()
	imIO.SetBackendFlags(imIO.BackendFlags() | newimgui.BackendFlagsRendererHasVtxOffset)

	gl.GenVertexArrays(1, &imguiInfo.VaoID)
	gl.GenBuffers(1, &imguiInfo.VboID)
	gl.GenBuffers(1, &imguiInfo.IndexBufID)
	gl.GenTextures(1, &imguiInfo.TexID)

	// Upload font to gpu
	gl.BindTexture(gl.TEXTURE_2D, imguiInfo.TexID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)

	pixels, width, height, _ := imIO.Fonts().GetTextureDataAsAlpha8()
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(width), int32(height), 0, gl.RED, gl.UNSIGNED_BYTE, pixels)

	// Store our identifier
	imIO.Fonts().SetTexID(newimgui.TextureID(uintptr(imguiInfo.TexID)))

	//Shader attributes
	imguiInfo.Mat.Bind()
	imguiInfo.Mat.EnableAttribute("Position")
	imguiInfo.Mat.EnableAttribute("UV")
	imguiInfo.Mat.EnableAttribute("Color")
	imguiInfo.Mat.UnBind()

	return imguiInfo
}

func SdlScancodeToImGuiKey(scancode sdl.Scancode) newimgui.Key {

	switch scancode {
	case sdl.SCANCODE_TAB:
		return newimgui.KeyTab
	case sdl.SCANCODE_LEFT:
		return newimgui.KeyLeftArrow
	case sdl.SCANCODE_RIGHT:
		return newimgui.KeyRightArrow
	case sdl.SCANCODE_UP:
		return newimgui.KeyUpArrow
	case sdl.SCANCODE_DOWN:
		return newimgui.KeyDownArrow
	case sdl.SCANCODE_PAGEUP:
		return newimgui.KeyPageUp
	case sdl.SCANCODE_PAGEDOWN:
		return newimgui.KeyPageDown
	case sdl.SCANCODE_HOME:
		return newimgui.KeyHome
	case sdl.SCANCODE_END:
		return newimgui.KeyEnd
	case sdl.SCANCODE_INSERT:
		return newimgui.KeyInsert
	case sdl.SCANCODE_DELETE:
		return newimgui.KeyDelete
	case sdl.SCANCODE_BACKSPACE:
		return newimgui.KeyBackspace
	case sdl.SCANCODE_SPACE:
		return newimgui.KeySpace
	case sdl.SCANCODE_RETURN:
		return newimgui.KeyEnter
	case sdl.SCANCODE_ESCAPE:
		return newimgui.KeyEscape
	case sdl.SCANCODE_APOSTROPHE:
		return newimgui.KeyApostrophe
	case sdl.SCANCODE_COMMA:
		return newimgui.KeyComma
	case sdl.SCANCODE_MINUS:
		return newimgui.KeyMinus
	case sdl.SCANCODE_PERIOD:
		return newimgui.KeyPeriod
	case sdl.SCANCODE_SLASH:
		return newimgui.KeySlash
	case sdl.SCANCODE_SEMICOLON:
		return newimgui.KeySemicolon
	case sdl.SCANCODE_EQUALS:
		return newimgui.KeyEqual
	case sdl.SCANCODE_LEFTBRACKET:
		return newimgui.KeyLeftBracket
	case sdl.SCANCODE_BACKSLASH:
		return newimgui.KeyBackslash
	case sdl.SCANCODE_RIGHTBRACKET:
		return newimgui.KeyRightBracket
	case sdl.SCANCODE_GRAVE:
		return newimgui.KeyGraveAccent
	case sdl.SCANCODE_CAPSLOCK:
		return newimgui.KeyCapsLock
	case sdl.SCANCODE_SCROLLLOCK:
		return newimgui.KeyScrollLock
	case sdl.SCANCODE_NUMLOCKCLEAR:
		return newimgui.KeyNumLock
	case sdl.SCANCODE_PRINTSCREEN:
		return newimgui.KeyPrintScreen
	case sdl.SCANCODE_PAUSE:
		return newimgui.KeyPause
	case sdl.SCANCODE_KP_0:
		return newimgui.KeyKeypad0
	case sdl.SCANCODE_KP_1:
		return newimgui.KeyKeypad1
	case sdl.SCANCODE_KP_2:
		return newimgui.KeyKeypad2
	case sdl.SCANCODE_KP_3:
		return newimgui.KeyKeypad3
	case sdl.SCANCODE_KP_4:
		return newimgui.KeyKeypad4
	case sdl.SCANCODE_KP_5:
		return newimgui.KeyKeypad5
	case sdl.SCANCODE_KP_6:
		return newimgui.KeyKeypad6
	case sdl.SCANCODE_KP_7:
		return newimgui.KeyKeypad7
	case sdl.SCANCODE_KP_8:
		return newimgui.KeyKeypad8
	case sdl.SCANCODE_KP_9:
		return newimgui.KeyKeypad9
	case sdl.SCANCODE_KP_PERIOD:
		return newimgui.KeyKeypadDecimal
	case sdl.SCANCODE_KP_DIVIDE:
		return newimgui.KeyKeypadDivide
	case sdl.SCANCODE_KP_MULTIPLY:
		return newimgui.KeyKeypadMultiply
	case sdl.SCANCODE_KP_MINUS:
		return newimgui.KeyKeypadSubtract
	case sdl.SCANCODE_KP_PLUS:
		return newimgui.KeyKeypadAdd
	case sdl.SCANCODE_KP_ENTER:
		return newimgui.KeyKeypadEnter
	case sdl.SCANCODE_KP_EQUALS:
		return newimgui.KeyKeypadEqual
	case sdl.SCANCODE_LSHIFT:
		return newimgui.KeyLeftShift
	case sdl.SCANCODE_LCTRL:
		return newimgui.KeyLeftCtrl
	case sdl.SCANCODE_LALT:
		return newimgui.KeyLeftAlt
	case sdl.SCANCODE_LGUI:
		return newimgui.KeyLeftSuper
	case sdl.SCANCODE_RSHIFT:
		return newimgui.KeyRightShift
	case sdl.SCANCODE_RCTRL:
		return newimgui.KeyRightCtrl
	case sdl.SCANCODE_RALT:
		return newimgui.KeyRightAlt
	case sdl.SCANCODE_RGUI:
		return newimgui.KeyRightSuper
	case sdl.SCANCODE_MENU:
		return newimgui.KeyMenu
	case sdl.SCANCODE_0:
		return newimgui.Key0
	case sdl.SCANCODE_1:
		return newimgui.Key1
	case sdl.SCANCODE_2:
		return newimgui.Key2
	case sdl.SCANCODE_3:
		return newimgui.Key3
	case sdl.SCANCODE_4:
		return newimgui.Key4
	case sdl.SCANCODE_5:
		return newimgui.Key5
	case sdl.SCANCODE_6:
		return newimgui.Key6
	case sdl.SCANCODE_7:
		return newimgui.Key7
	case sdl.SCANCODE_8:
		return newimgui.Key8
	case sdl.SCANCODE_9:
		return newimgui.Key9
	case sdl.SCANCODE_A:
		return newimgui.KeyA
	case sdl.SCANCODE_B:
		return newimgui.KeyB
	case sdl.SCANCODE_C:
		return newimgui.KeyC
	case sdl.SCANCODE_D:
		return newimgui.KeyD
	case sdl.SCANCODE_E:
		return newimgui.KeyE
	case sdl.SCANCODE_F:
		return newimgui.KeyF
	case sdl.SCANCODE_G:
		return newimgui.KeyG
	case sdl.SCANCODE_H:
		return newimgui.KeyH
	case sdl.SCANCODE_I:
		return newimgui.KeyI
	case sdl.SCANCODE_J:
		return newimgui.KeyJ
	case sdl.SCANCODE_K:
		return newimgui.KeyK
	case sdl.SCANCODE_L:
		return newimgui.KeyL
	case sdl.SCANCODE_M:
		return newimgui.KeyM
	case sdl.SCANCODE_N:
		return newimgui.KeyN
	case sdl.SCANCODE_O:
		return newimgui.KeyO
	case sdl.SCANCODE_P:
		return newimgui.KeyP
	case sdl.SCANCODE_Q:
		return newimgui.KeyQ
	case sdl.SCANCODE_R:
		return newimgui.KeyR
	case sdl.SCANCODE_S:
		return newimgui.KeyS
	case sdl.SCANCODE_T:
		return newimgui.KeyT
	case sdl.SCANCODE_U:
		return newimgui.KeyU
	case sdl.SCANCODE_V:
		return newimgui.KeyV
	case sdl.SCANCODE_W:
		return newimgui.KeyW
	case sdl.SCANCODE_X:
		return newimgui.KeyX
	case sdl.SCANCODE_Y:
		return newimgui.KeyY
	case sdl.SCANCODE_Z:
		return newimgui.KeyZ
	case sdl.SCANCODE_F1:
		return newimgui.KeyF1
	case sdl.SCANCODE_F2:
		return newimgui.KeyF2
	case sdl.SCANCODE_F3:
		return newimgui.KeyF3
	case sdl.SCANCODE_F4:
		return newimgui.KeyF4
	case sdl.SCANCODE_F5:
		return newimgui.KeyF5
	case sdl.SCANCODE_F6:
		return newimgui.KeyF6
	case sdl.SCANCODE_F7:
		return newimgui.KeyF7
	case sdl.SCANCODE_F8:
		return newimgui.KeyF8
	case sdl.SCANCODE_F9:
		return newimgui.KeyF9
	case sdl.SCANCODE_F10:
		return newimgui.KeyF10
	case sdl.SCANCODE_F11:
		return newimgui.KeyF11
	case sdl.SCANCODE_F12:
		return newimgui.KeyF12
	default:
		return newimgui.KeyNone
	}
}
