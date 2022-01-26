package buffers

import (
	"github.com/bloeys/nmage/logging"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type BufferLayoutElement struct {
	Offset int
	DataType
}

type BufferLayout struct {
	Elements []BufferLayoutElement
}

type Buffer struct {
	VAOID uint32
	//BufID is the ID of the VBO
	BufID uint32
	//IndexBufID is the ID of the index/element buffer
	IndexBufID    uint32
	IndexBufCount int32
	Layout        BufferLayout
	Stride        int32
}

func (b *Buffer) Bind() {
	gl.BindVertexArray(b.VAOID)
}

func (b *Buffer) UnBind() {
	gl.BindVertexArray(0)
}

func (b *Buffer) CalcValues() {

	b.Stride = 0
	for i := 0; i < len(b.Layout.Elements); i++ {

		info := GetDataTypeInfo(b.Layout.Elements[i].DataType)
		b.Layout.Elements[i].Offset = int(b.Stride)
		b.Stride += info.GetSize()
	}
}

func (b *Buffer) SetData(values []float32) {

	gl.BindVertexArray(b.VAOID)
	gl.BindBuffer(gl.ARRAY_BUFFER, b.BufID)

	gl.BufferData(gl.ARRAY_BUFFER, len(values)*4, gl.Ptr(values), BufUsageStatic.ToGL())

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (b *Buffer) SetIndexBufData(values []uint32) {

	b.IndexBufCount = int32(len(values))
	gl.BindVertexArray(b.VAOID)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.IndexBufID)

	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(values)*4, gl.Ptr(values), BufUsageStatic.ToGL())

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func NewBuffer(layout BufferLayout) Buffer {

	b := Buffer{
		Layout: layout,
	}

	gl.GenVertexArrays(1, &b.VAOID)
	if b.VAOID == 0 {
		logging.ErrLog.Println("Failed to create openGL vertex array object")
	}

	gl.GenBuffers(1, &b.BufID)
	if b.BufID == 0 {
		logging.ErrLog.Println("Failed to create openGL buffer")
	}

	gl.GenBuffers(1, &b.IndexBufID)
	if b.IndexBufID == 0 {
		logging.ErrLog.Println("Failed to create openGL buffer")
	}

	b.CalcValues()
	return b
}
