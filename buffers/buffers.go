package buffers

import (
	"github.com/bloeys/go-sdl-engine/logging"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type BufGLType int

const (
	BufGLTypeUnknown BufGLType = 0
	//Generic array of data. Should be used for most data like vertex positions, vertex colors etc.
	BufGLTypeArray   BufGLType = gl.ARRAY_BUFFER
	BufGLTypeIndices BufGLType = gl.ELEMENT_ARRAY_BUFFER
)

type BufType int

const (
	BufTypeUnknown BufType = iota
	BufTypeVertPos
	BufTypeColor
	BufTypeIndex
	BufTypeNormal
)

func (bt BufType) GetBufferGLType() BufGLType {
	switch bt {

	case BufTypeNormal:
		fallthrough
	case BufTypeColor:
		fallthrough
	case BufTypeVertPos:
		return BufGLTypeArray

	case BufTypeIndex:
		return BufGLTypeIndices
	default:
		logging.WarnLog.Println("Unknown BufferType. BufferType: ", bt)
		return BufGLTypeUnknown
	}
}

type BufUsage int

const (
	//Buffer is set only once and used many times
	BufUsageStatic BufUsage = gl.STATIC_DRAW
	//Buffer is changed a lot and used many times
	BufUsageDynamic BufUsage = gl.DYNAMIC_DRAW
	//Buffer is set only once and used by the GPU at most a few times
	BufUsageStream BufUsage = gl.STREAM_DRAW
)

type Buffer struct {
	ID     uint32
	Type   BufType
	GLType BufGLType
	DataTypeInfo
}

func (b *Buffer) Activate() {
	gl.BindBuffer(uint32(b.GLType), b.ID)
}

func (b *Buffer) Deactivate() {
	gl.BindBuffer(uint32(b.GLType), 0)
}

type BufferObject struct {
	VAOID      uint32
	VertPosBuf *Buffer
	NormalBuf  *Buffer
	ColorBuf   *Buffer
	IndexBuf   *Buffer
}

func (bo *BufferObject) GenBuffer(data []float32, bufUsage BufUsage, bufType BufType, bufDataType DataType) {

	gl.BindVertexArray(bo.VAOID)

	//Create vertex buffer object
	var vboID uint32
	gl.GenBuffers(1, &vboID)
	if vboID == 0 {
		logging.ErrLog.Println("Failed to create openGL buffer")
	}

	buf := &Buffer{
		ID:           vboID,
		Type:         bufType,
		GLType:       bufType.GetBufferGLType(),
		DataTypeInfo: GetDataTypeInfo(bufDataType),
	}
	bo.SetBuffer(buf)

	//Fill buffer with data
	gl.BindBuffer(uint32(buf.GLType), buf.ID)
	gl.BufferData(uint32(buf.GLType), int(buf.DataTypeInfo.ElementSize)*len(data), gl.Ptr(data), uint32(bufUsage))

	//Unbind everything
	gl.BindVertexArray(0)
	gl.BindBuffer(uint32(buf.GLType), 0)
}

func (bo *BufferObject) GenBufferUint32(data []uint32, bufUsage BufUsage, bufType BufType, bufDataType DataType) {

	gl.BindVertexArray(bo.VAOID)

	//Create vertex buffer object
	var vboID uint32
	gl.GenBuffers(1, &vboID)
	if vboID == 0 {
		logging.ErrLog.Println("Failed to create openGL buffer")
	}

	buf := &Buffer{
		ID:           vboID,
		Type:         bufType,
		GLType:       bufType.GetBufferGLType(),
		DataTypeInfo: GetDataTypeInfo(bufDataType),
	}
	bo.SetBuffer(buf)

	//Fill buffer with data
	gl.BindBuffer(uint32(buf.GLType), buf.ID)
	gl.BufferData(uint32(buf.GLType), int(buf.DataTypeInfo.ElementSize)*len(data), gl.Ptr(data), uint32(bufUsage))

	//Unbind everything
	gl.BindVertexArray(0)
	gl.BindBuffer(uint32(buf.GLType), 0)
}

func (bo *BufferObject) SetBuffer(buf *Buffer) {

	switch buf.Type {
	case BufTypeVertPos:
		bo.VertPosBuf = buf
	case BufTypeNormal:
		bo.NormalBuf = buf
	case BufTypeColor:
		bo.ColorBuf = buf
	case BufTypeIndex:
		bo.IndexBuf = buf
	default:
		logging.WarnLog.Println("Unknown buffer type in SetBuffer. Type:", buf.Type)
	}
}

func (bo *BufferObject) Activate() {
	gl.BindVertexArray(bo.VAOID)
}

func (bo *BufferObject) Deactivate() {
	gl.BindVertexArray(0)
}

func NewBufferObject() *BufferObject {

	var vaoID uint32
	gl.CreateVertexArrays(1, &vaoID)
	if vaoID == 0 {
		logging.ErrLog.Println("Failed to create openGL vertex array object")
	}

	return &BufferObject{VAOID: vaoID}
}
