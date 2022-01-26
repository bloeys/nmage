package buffers

import (
	"fmt"

	"github.com/bloeys/nmage/asserts"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type BufUsage int

const (
	//Buffer is set only once and used many times
	BufUsageStatic BufUsage = iota
	//Buffer is changed a lot and used many times
	BufUsageDynamic
	//Buffer is set only once and used by the GPU at most a few times
	BufUsageStream
)

func (b BufUsage) ToGL() uint32 {
	switch b {
	case BufUsageStatic:
		return gl.STATIC_DRAW
	case BufUsageDynamic:
		return gl.DYNAMIC_DRAW
	case BufUsageStream:
		return gl.STREAM_DRAW
	}

	asserts.T(false, fmt.Sprintf("Unexpected BufUsage value '%v'", b))
	return 0
}
