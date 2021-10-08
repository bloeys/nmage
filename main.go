package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth  int32 = 800
	winHeight int32 = 600
)

var (
	isRunning = true
	window    *sdl.Window
)

func main() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic("Failed to init SDL. Err: " + err.Error())
	}
	defer sdl.Quit()

	window, err = sdl.CreateWindow(
		"test",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		winWidth,
		winHeight,
		sdl.WINDOW_SHOWN)
	if err != nil {
		panic("Failed to create window. Err: " + err.Error())
	}
	defer window.Destroy()

	rend, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic("Creating renderer failed. Err: " + err.Error())
	}
	defer rend.Destroy()

	tex, err := rend.CreateTexture(sdl.PIXELFORMAT_ABGR8888, 0, winWidth, winHeight)
	if err != nil {
		panic(err)
	}

	//x4 to allow for RGBA
	pixels := make([]byte, winHeight*winWidth*4)
	for y := 0; y < int(winHeight); y++ {
		for x := 0; x < int(winWidth); x++ {

			c := sdl.Color{
				R: byte(int(float64(x)/float64(winWidth)*256) % 256),
				G: byte(int(float64(y)/float64(winHeight)*256) % 256),
			}

			setPixel(x, y, c, pixels)
		}
	}

	//Update texture with new pixel values
	tex.Update(nil, pixels, int(winWidth)*4)
	//Copy texture to renderer
	rend.Copy(tex, nil, nil)
	//Blit
	rend.Present()

	var fps float32 = 60
	var dt float32 = 1.0 / 60.0
	var dtLimit float32 = 1.0 / 120.0
	for isRunning {

		frameStartTime := time.Now()

		handleEvents()

		//If FPS is more than 120 then limit to that
		dt = float32(time.Since(frameStartTime).Seconds())
		if dt < dtLimit {
			sdl.Delay(8 - uint32(dt*1000))
			dt = float32(time.Since(frameStartTime).Seconds())
		}

		//Display FPS is the average of the FPS of this frame and the last frame
		fps = (fps + 1/dt) / 2
		window.SetTitle(fmt.Sprintf("FPS: %.2f; dt: %.3f", fps, dt))
	}
}

func handleEvents() {

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

		switch e := event.(type) {

		case *sdl.QuitEvent:
			println("Quit at ", e.Timestamp)
			isRunning = false
		}
	}
}

//handleEvents assumes sdl.PIXELFORMAT_ABGR8888
func setPixel(x, y int, c sdl.Color, pixels []byte) {

	index := (y*int(winWidth) + x) * 4
	pixels[index] = c.R
	pixels[index+1] = c.G
	pixels[index+2] = c.B
	pixels[index+3] = c.A
}
