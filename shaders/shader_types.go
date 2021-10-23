package shaders

import "github.com/go-gl/gl/v4.6-compatibility/gl"

type ShaderType int

const (
	VertexShaderType   ShaderType = gl.VERTEX_SHADER
	FragmentShaderType ShaderType = gl.FRAGMENT_SHADER
)
