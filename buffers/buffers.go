package buffers

import (
	"github.com/bloeys/go-sdl-engine/logging"
	"github.com/bloeys/go-sdl-engine/shaders"
	"github.com/go-gl/gl/v4.6-compatibility/gl"
)

func HandleBuffers(sp shaders.ShaderProgram) {

	//Create and fill Vertex buffer object
	var vboID uint32
	gl.CreateBuffers(1, &vboID)
	if vboID == 0 {
		logging.ErrLog.Println("Failed to create openGL buffer")
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, vboID)

	vertices := []float32{
		-0.5, 0.5, 0,
		0.5, 0.5, 0,
		-0.5, -0.5, 0,

		0.5, 0.5, 0,
		0.5, -0.5, 0,
		-0.5, -0.5, 0,
	}
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	//Assign the VBO to vertPos attribute
	vertPosLoc := sp.GetAttribLoc("vertPos")
	gl.VertexAttribPointer(uint32(vertPosLoc), 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(uint32(vertPosLoc))
}
