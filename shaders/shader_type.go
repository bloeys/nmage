package shaders

import (
	"github.com/bloeys/go-sdl-engine/logging"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type ShaderType int

const (
	Unknown ShaderType = iota
	Vertex
	Fragment
)

//GLType returns the GL shader type of this ShaderType
//Panics if not known
func (t ShaderType) GLType() uint32 {

	switch t {
	case Vertex:
		return gl.VERTEX_SHADER
	case Fragment:
		return gl.FRAGMENT_SHADER
	}

	logging.ErrLog.Panicf("Converting ShaderType->GL Shader Type failed. Unknown ShaderType of value: %v\n", t)
	return 0
}

//FromGLShaderType returns the ShaderType of the passed GL shader type.
//Panics if not known
func (t ShaderType) FromGLShaderType(glShaderType int) ShaderType {

	switch glShaderType {
	case gl.VERTEX_SHADER:
		return Vertex
	case gl.FRAGMENT_SHADER:
		return Fragment
	default:
		logging.ErrLog.Panicf("Converting GL shader type->ShaderType failed. Unknown GL shader type of value: %v\n", glShaderType)
		return Unknown
	}
}
