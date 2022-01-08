package buffers

import (
	"github.com/bloeys/nmage/logging"
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
	//GLType is the type of the variable represented (e.g. for vec3 its gl.FLOAT_VEC2)
	GLType uint32
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
			GLType:       gl.INT,
		}

	case DataTypeFloat32:
		return DataTypeInfo{
			ElementSize:  4,
			ElementCount: 1,
			ElementType:  gl.FLOAT,
			GLType:       gl.FLOAT,
		}
	case DataTypeFloat64:
		return DataTypeInfo{
			ElementSize:  8,
			ElementCount: 1,
			ElementType:  gl.DOUBLE,
			GLType:       gl.DOUBLE,
		}

	case DataTypeVec2:
		return DataTypeInfo{
			ElementSize:  4,
			ElementCount: 2,
			ElementType:  gl.FLOAT,
			GLType:       gl.FLOAT_VEC2,
		}
	case DataTypeVec3:
		return DataTypeInfo{
			ElementSize:  4,
			ElementCount: 3,
			ElementType:  gl.FLOAT,
			GLType:       gl.FLOAT_VEC3,
		}
	case DataTypeVec4:
		return DataTypeInfo{
			ElementSize:  4,
			ElementCount: 4,
			ElementType:  gl.FLOAT,
			GLType:       gl.FLOAT_VEC4,
		}

	default:
		logging.WarnLog.Println("Unknown data type passed. DataType:", dt)
		return DataTypeInfo{}
	}
}
