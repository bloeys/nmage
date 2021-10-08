package timing

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	fps            float32   = 60
	dt             float32   = 1.0 / 60.0
	dtLimit        float32   = 1.0 / 120.0
	frameStartTime time.Time = time.Now()
)

func FrameStarted() {
	frameStartTime = time.Now()
}

func FrameEnded() {

	//If FPS is more than 120 then limit to that
	dt = float32(time.Since(frameStartTime).Seconds())
	if dt < dtLimit {
		sdl.Delay(8 - uint32(dt*1000))
		dt = float32(time.Since(frameStartTime).Seconds())
	}

	//Display FPS is the average of the FPS of this frame and the last frame
	fps = (fps + 1/dt) / 2
}

//DT returns last frame delta time (number of seconds frame took)
func DT() float32 {
	return dt
}

//FPS returns fps
func FPS() float32 {
	return fps
}
