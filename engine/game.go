package engine

import (
	"github.com/bloeys/nmage/timing"
	nmageimgui "github.com/bloeys/nmage/ui/imgui"
	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	isRunning = false
)

type Game interface {
	Init()

	Update()
	Render()
	FrameEnd()

	DeInit()
}

func Run(g Game, w *Window, ui nmageimgui.ImguiInfo) {

	isRunning = true

	// Simulate 2 imgui frames, one before and one after g.Init so any imgui calls are allowed within init.
	// Calling before is required for things like push font.
	tempWidth, tempHeight := w.SDLWin.GetSize()
	tempFBWidth, tempFBHeight := w.SDLWin.GLGetDrawableSize()
	ui.FrameStart(float32(tempWidth), float32(tempHeight))
	ui.Render(float32(tempWidth), float32(tempHeight), tempFBWidth, tempFBHeight)

	g.Init()

	// Second imgui frame
	tempWidth, tempHeight = w.SDLWin.GetSize()
	tempFBWidth, tempFBHeight = w.SDLWin.GLGetDrawableSize()
	ui.FrameStart(float32(tempWidth), float32(tempHeight))
	ui.Render(float32(tempWidth), float32(tempHeight), tempFBWidth, tempFBHeight)

	for isRunning {

		//PERF: Cache these
		width, height := w.SDLWin.GetSize()
		fbWidth, fbHeight := w.SDLWin.GLGetDrawableSize()

		timing.FrameStarted()
		w.handleInputs()
		ui.FrameStart(float32(width), float32(height))

		g.Update()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
		g.Render()
		ui.Render(float32(width), float32(height), fbWidth, fbHeight)
		w.SDLWin.GLSwap()

		g.FrameEnd()
		w.Rend.FrameEnd()
		timing.FrameEnded()
	}

	g.DeInit()
}

func Quit() {
	isRunning = false
}
