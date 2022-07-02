package timing

import (
	"time"
)

var (
	dt         float32 = 0.01
	frameStart time.Time
	startTime  time.Time

	//fps calculator vars
	dtAccum                  float32 = 1
	lastElapsedTime          uint64  = 0
	framesSinceLastFPSUpdate uint    = 0
	avgFps                   float32 = 1
)

func Init() {
	startTime = time.Now()
}

func FrameStarted() {

	frameStart = time.Now()

	//fps stuff
	dtAccum += dt
	framesSinceLastFPSUpdate++
	et := ElapsedTime()
	if et > lastElapsedTime {
		avgDT := dtAccum / float32(framesSinceLastFPSUpdate)
		avgFps = 1 / avgDT

		dtAccum = 0
		framesSinceLastFPSUpdate = 0
	}
	lastElapsedTime = et
}

func FrameEnded() {

	//Calculate new dt
	dt = float32(time.Since(frameStart).Seconds())
	if dt == 0 {
		dt = float32(time.Microsecond.Seconds())
	}
}

//DT is frame deltatime in seconds
func DT() float32 {
	return dt
}

//GetAvgFPS returns the fps averaged over 1 second
func GetAvgFPS() float32 {
	return avgFps
}

//ElapsedTime is time since game start
func ElapsedTime() uint64 {
	return uint64(time.Since(startTime).Seconds())
}
