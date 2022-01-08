package asserts

import (
	"github.com/bloeys/nmage/consts"
	"github.com/bloeys/nmage/logging"
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
