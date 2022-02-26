package engine

import (
	"github.com/bloeys/nmage/timing"
	nmageimgui "github.com/bloeys/nmage/ui/imgui"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Game interface {
	Init()

	FrameStart()
	Update()
	Render()
	FrameEnd()
	ShouldRun() bool

	GetWindow() *Window
	GetImGUI() nmageimgui.ImguiInfo

	Deinit()
}

func Run(g Game) {

	g.Init()

	w := g.GetWindow()
	ui := g.GetImGUI()
	for g.ShouldRun() {

		width, height := w.SDLWin.GetSize()
		fbWidth, fbHeight := w.SDLWin.GLGetDrawableSize()

		timing.FrameStarted()
		w.handleInputs()
		ui.FrameStart(float32(width), float32(height))

		g.FrameStart()

		g.Update()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		g.Render()
		ui.Render(float32(width), float32(height), fbWidth, fbHeight)
		w.SDLWin.GLSwap()

		g.FrameEnd()
		w.Rend.FrameEnd()
		timing.FrameEnded()
	}

	g.Deinit()
}
