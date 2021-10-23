package shaders

import (
	"github.com/bloeys/go-sdl-engine/logging"
	"github.com/go-gl/gl/v4.6-compatibility/gl"
)

type ShaderProgram struct {
	ID           uint32
	VertShaderID uint32
	FragShaderID uint32
}

func (sp *ShaderProgram) AttachShader(shader Shader) {

	gl.AttachShader(sp.ID, shader.ID)
	switch shader.ShaderType {
	case VertexShaderType:
		sp.VertShaderID = shader.ID
	case FragmentShaderType:
		sp.FragShaderID = shader.ID
	default:
		logging.ErrLog.Println("Unknown shader type ", shader.ShaderType, " for ID ", shader.ID)
	}
}

func (sp *ShaderProgram) Link() {

	gl.LinkProgram(sp.ID)

	if sp.VertShaderID != 0 {
		gl.DeleteShader(sp.VertShaderID)
	}
	if sp.FragShaderID != 0 {
		gl.DeleteShader(sp.FragShaderID)
	}
}

func (sp *ShaderProgram) GetAttribLoc(attribName string) int32 {
	return gl.GetAttribLocation(sp.ID, gl.Str(attribName+"\x00"))
}

func (sp *ShaderProgram) Use() {
	gl.UseProgram(sp.ID)
}

func (sp *ShaderProgram) Delete() {
	gl.DeleteProgram(sp.ID)
}
