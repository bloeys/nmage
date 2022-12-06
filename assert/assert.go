package assert

import (
	"fmt"

	"github.com/bloeys/nmage/consts"
	"github.com/bloeys/nmage/logging"
)

func T(check bool, msg string, args ...any) {

	if consts.Debug && !check {
		logging.ErrLog.Panicln("Assert failed:", fmt.Sprintf(msg, args...))
	}
}
