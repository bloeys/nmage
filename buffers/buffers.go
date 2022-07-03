package buffers

import (
	"github.com/bloeys/nmage/logging"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Buffer struct {
	VAOID uint32
	//BufID is the ID of the VBO
	BufID uint32
	//IndexBufID is the ID of the index/element buffer
	IndexBufID    uint32
	IndexBufCount int32
	Stride        int32

	layout []Element
}

func (b *Buffer) Bind() {
	gl.BindVertexArray(b.VAOID)
}

func (b *Buffer) UnBind() {
	gl.BindVertexArray(0)
}

func (b *Buffer) SetData(values []float32) {

	gl.BindVertexArray(b.VAOID)
	gl.BindBuffer(gl.ARRAY_BUFFER, b.BufID)

	sizeInBytes := len(values) * 4
	if sizeInBytes == 0 {
		gl.BufferData(gl.ARRAY_BUFFER, 0, gl.Ptr(nil), BufUsage_Static.ToGL())
	} else {
		gl.BufferData(gl.ARRAY_BUFFER, sizeInBytes, gl.Ptr(&values[0]), BufUsage_Static.ToGL())
	}

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (b *Buffer) SetIndexBufData(values []uint32) {

	b.IndexBufCount = int32(len(values))
	gl.BindVertexArray(b.VAOID)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.IndexBufID)

	sizeInBytes := len(values) * 4
	if sizeInBytes == 0 {
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 0, gl.Ptr(nil), BufUsage_Static.ToGL())
	} else {
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, sizeInBytes, gl.Ptr(&values[0]), BufUsage_Static.ToGL())
	}

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func (b *Buffer) GetLayout() []Element {
	e := make([]Element, len(b.layout))
	copy(e, b.layout)
	return e
}

//SetLayout updates the layout object and the corresponding vertex attributes.
//Vertex attributes are also enabled.
func (b *Buffer) SetLayout(layout ...Element) {

	b.layout = layout

	b.Stride = 0
	for i := 0; i < len(b.layout); i++ {

		b.layout[i].Offset = int(b.Stride)
		b.Stride += b.layout[i].Size()
	}

	//Set opengl stuff
	b.Bind()

	//NOTE: VBOs are only bound at 'VertexAttribPointer', not BindBUffer, so we need to bind the buffer and vao here
	gl.BindBuffer(gl.ARRAY_BUFFER, b.BufID)

	for i := 0; i < len(layout); i++ {
		gl.EnableVertexAttribArray(uint32(i))
		gl.VertexAttribPointer(uint32(i), layout[i].ElementType.CompCount(), layout[i].ElementType.GLType(), false, b.Stride, gl.PtrOffset(layout[i].Offset))
	}

	b.UnBind()
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func NewBuffer(layout ...Element) Buffer {

	b := Buffer{}
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

	b.SetLayout(layout...)
	return b
}
