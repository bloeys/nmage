package buffers

import (
	"github.com/bloeys/go-sdl-engine/logging"
	"github.com/go-gl/gl/v4.6-compatibility/gl"
)

type BufferType int

const (
	BufTypeUnknown BufferType = iota
	BufTypeVertPos
	BufTypeColor
)

type BufferGLType int

const (
	//Generic array of data. Should be used for most data like vertex positions, vertex colors etc.
	BufGLTypeArrayBuffer BufferGLType = gl.ARRAY_BUFFER
)

type BufferUsage int

const (
	//Buffer is set only once and used many times
	BufUsageStatic BufferUsage = gl.STATIC_DRAW
	//Buffer is changed a lot and used many times
	BufUsageDynamic BufferUsage = gl.DYNAMIC_DRAW
	//Buffer is set only once and used by the GPU at most a few times
	BufUsageStream BufferUsage = gl.STREAM_DRAW
)

type Buffer struct {
	ID     uint32
	Type   BufferType
	GLType BufferGLType
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
	ColorBuf   *Buffer
}

func (bo *BufferObject) GenBuffer(data []float32, bufUsage BufferUsage, bufType BufferType, bufDataType DataType) {

	gl.BindVertexArray(bo.VAOID)

	//Create vertex buffer object
	var vboID uint32
	gl.CreateBuffers(1, &vboID)
	if vboID == 0 {
		logging.ErrLog.Println("Failed to create openGL buffer")
	}

	buf := &Buffer{
		ID:           vboID,
		Type:         bufType,
		GLType:       BufGLTypeArrayBuffer,
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
	case BufTypeColor:
		bo.ColorBuf = buf
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
