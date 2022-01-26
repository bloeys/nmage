package asserts

import (
	"github.com/bloeys/nmage/consts"
	"github.com/bloeys/nmage/logging"
)

func T(check bool, msg string) {

	if !consts.Debug || check {
		return
	}

	logging.ErrLog.Panicln("Assert failed:", msg)
}
