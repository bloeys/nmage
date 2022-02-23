package timing

import "time"

var (
	dt         float32 = 0.01
	frameStart time.Time
	startTime  time.Time
)

func Init() {
	startTime = time.Now()
}

func FrameStarted() {
	frameStart = time.Now()
}

func FrameEnded() {
	dt = float32(time.Since(frameStart).Seconds())
}

//DT is frame deltatime in seconds
func DT() float32 {
	return dt
}

//ElapsedTime is time since game start
func ElapsedTime() uint64 {
	return uint64(time.Since(startTime).Seconds())
}
