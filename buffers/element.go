package buffers

import (
	"fmt"

	"github.com/bloeys/nmage/assert"
	"github.com/go-gl/gl/v4.1-core/gl"
)

//Element represents an element that makes up a buffer (e.g. Vec3 at an offset of 12 bytes)
type Element struct {
	Offset int
	ElementType
}

//ElementType is the type of an element thats makes up a buffer (e.g. Vec3)
type ElementType int

const (
	DataTypeUnknown ElementType = iota
	DataTypeUint32
	DataTypeInt32
	DataTypeFloat32

	DataTypeVec2
	DataTypeVec3
	DataTypeVec4
)

func (dt ElementType) GLType() uint32 {

	switch dt {
	case DataTypeUint32:
		return gl.UNSIGNED_INT
	case DataTypeInt32:
		return gl.INT

	case DataTypeFloat32:
		fallthrough
	case DataTypeVec2:
		fallthrough
	case DataTypeVec3:
		fallthrough
	case DataTypeVec4:
		return gl.FLOAT

	default:
		assert.T(false, fmt.Sprintf("Unknown data type passed. DataType '%v'", dt))
		return 0
	}
}

//CompSize returns the size in bytes for one component of the type (e.g. for Vec2 its 4)
func (dt ElementType) CompSize() int32 {

	switch dt {
	case DataTypeUint32:
		fallthrough
	case DataTypeFloat32:
		fallthrough
	case DataTypeInt32:
		fallthrough
	case DataTypeVec2:
		fallthrough
	case DataTypeVec3:
		fallthrough
	case DataTypeVec4:
		return 4

	default:
		assert.T(false, fmt.Sprintf("Unknown data type passed. DataType '%v'", dt))
		return 0
	}
}

//CompCount returns the number of components in the element (e.g. for Vec2 its 2)
func (dt ElementType) CompCount() int32 {

	switch dt {
	case DataTypeUint32:
		fallthrough
	case DataTypeFloat32:
		fallthrough
	case DataTypeInt32:
		return 1

	case DataTypeVec2:
		return 2
	case DataTypeVec3:
		return 3
	case DataTypeVec4:
		return 4

	default:
		assert.T(false, fmt.Sprintf("Unknown data type passed. DataType '%v'", dt))
		return 0
	}
}

//Size returns the total size in bytes (e.g. for vec3 its 3*4=12 bytes)
func (dt ElementType) Size() int32 {

	switch dt {
	case DataTypeUint32:
		fallthrough
	case DataTypeFloat32:
		fallthrough
	case DataTypeInt32:
		return 4

	case DataTypeVec2:
		return 2 * 4
	case DataTypeVec3:
		return 3 * 4
	case DataTypeVec4:
		return 4 * 4

	default:
		assert.T(false, fmt.Sprintf("Unknown data type passed. DataType '%v'", dt))
		return 0
	}
}
