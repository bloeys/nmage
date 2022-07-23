package assert

import (
	"github.com/bloeys/nmage/consts"
	"github.com/bloeys/nmage/logging"
)

func T(check bool, msg string) {
	if consts.Debug && !check {
		logging.ErrLog.Panicln("Assert failed:", msg)
	}
}
