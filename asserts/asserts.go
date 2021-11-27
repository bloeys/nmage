package asserts

import (
	"github.com/bloeys/go-sdl-engine/consts"
	"github.com/bloeys/go-sdl-engine/logging"
)

func True(check bool, msg string) {
	if consts.Debug && !check {
		logging.ErrLog.Panicln(msg)
	}
}

func False(check bool, msg string) {
	if consts.Debug && check {
		logging.ErrLog.Panicln(msg)
	}
}
