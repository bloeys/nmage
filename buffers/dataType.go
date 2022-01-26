package buffers

import (
	"fmt"

	"github.com/bloeys/nmage/asserts"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type DataType int

const (
	DataTypeUnknown DataType = iota
	DataTypeUint32
	DataTypeInt32
	DataTypeFloat32
	DataTypeFloat64

	DataTypeVec2
	DataTypeVec3
	DataTypeVec4
)

type DataTypeInfo struct {
	//ElementSize is size in bytes of one element (e.g. for vec3 its 4)
	ElementSize int32
	//ElementCount is number of elements (e.g. for vec3 its 3)
	ElementCount int32
	//ElementType is the type of each primitive (e.g. for vec3 its gl.FLOAT)
	ElementType uint32
}

//GetSize returns the total size in bytes (e.g. for vec3 its 4*3)
func (dti *DataTypeInfo) GetSize() int32 {
	return dti.ElementSize * dti.ElementCount
}

func GetDataTypeInfo(dt DataType) DataTypeInfo {

	switch dt {
	case DataTypeUint32:
		fallthrough
	case DataTypeInt32:
		return DataTypeInfo{
			ElementSize:  4,
			ElementCount: 1,
			ElementType:  gl.INT,
		}

	case DataTypeFloat32:
		return DataTypeInfo{
			ElementSize:  4,
			ElementCount: 1,
			ElementType:  gl.FLOAT,
		}
	case DataTypeFloat64:
		return DataTypeInfo{
			ElementSize:  8,
			ElementCount: 1,
			ElementType:  gl.DOUBLE,
		}

	case DataTypeVec2:
		return DataTypeInfo{
			ElementSize:  4,
			ElementCount: 2,
			ElementType:  gl.FLOAT,
		}
	case DataTypeVec3:
		return DataTypeInfo{
			ElementSize:  4,
			ElementCount: 3,
			ElementType:  gl.FLOAT,
		}
	case DataTypeVec4:
		return DataTypeInfo{
			ElementSize:  4,
			ElementCount: 4,
			ElementType:  gl.FLOAT,
		}

	default:
		asserts.T(false, fmt.Sprintf("Unknown data type passed. DataType '%v'", dt))
		return DataTypeInfo{}
	}
}
