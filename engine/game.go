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

	w := g.GetWindow()
	ui := g.GetImGUI()

	//Simulate an imgui frame during init so any imgui calls are allowed within init
	tempWidth, tempHeight := w.SDLWin.GetSize()
	tempFBWidth, tempFBHeight := w.SDLWin.GLGetDrawableSize()
	ui.FrameStart(float32(tempWidth), float32(tempHeight))
	g.Init()
	ui.Render(float32(tempWidth), float32(tempHeight), tempFBWidth, tempFBHeight)

	for g.ShouldRun() {

		//PERF: Cache these
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
