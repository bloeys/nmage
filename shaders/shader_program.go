package shaders

import (
	"github.com/bloeys/gglm/gglm"
	"github.com/bloeys/go-sdl-engine/buffers"
	"github.com/bloeys/go-sdl-engine/logging"
	"github.com/go-gl/gl/v4.1-core/gl"
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

func (sp *ShaderProgram) Activate() {
	gl.UseProgram(sp.ID)
}

func (sp *ShaderProgram) Deactivate() {
	gl.UseProgram(0)
}

func (sp *ShaderProgram) GetAttribLoc(attribName string) int32 {
	return gl.GetAttribLocation(sp.ID, gl.Str(attribName+"\x00"))
}

func (sp *ShaderProgram) SetAttribute(attribName string, bufObj *buffers.BufferObject, buf *buffers.Buffer) {

	bufObj.Activate()
	buf.Activate()

	attribLoc := sp.GetAttribLoc(attribName)
	gl.VertexAttribPointer(uint32(attribLoc), buf.ElementCount, buf.ElementType, false, buf.GetSize(), gl.PtrOffset(0))

	bufObj.Activate()
	buf.Deactivate()
}

func (sp *ShaderProgram) EnableAttribute(attribName string) {
	gl.EnableVertexAttribArray(uint32(sp.GetAttribLoc(attribName)))
}

func (sp *ShaderProgram) DisableAttribute(attribName string) {
	gl.DisableVertexAttribArray(uint32(sp.GetAttribLoc(attribName)))
}

func (sp *ShaderProgram) SetUnifFloat32(uniformName string, val float32) {
	loc := gl.GetUniformLocation(sp.ID, gl.Str(uniformName+"\x00"))
	gl.ProgramUniform1f(sp.ID, loc, val)
}

func (sp *ShaderProgram) SetUnifVec2(uniformName string, vec2 *gglm.Vec2) {
	loc := gl.GetUniformLocation(sp.ID, gl.Str(uniformName+"\x00"))
	gl.ProgramUniform2fv(sp.ID, loc, 1, &vec2.Data[0])
}

func (sp *ShaderProgram) SetUnifVec3(uniformName string, vec3 *gglm.Vec3) {
	loc := gl.GetUniformLocation(sp.ID, gl.Str(uniformName+"\x00"))
	gl.ProgramUniform3fv(sp.ID, loc, 1, &vec3.Data[0])
}

func (sp *ShaderProgram) SetUnifVec4(uniformName string, vec4 *gglm.Vec4) {
	loc := gl.GetUniformLocation(sp.ID, gl.Str(uniformName+"\x00"))
	gl.ProgramUniform4fv(sp.ID, loc, 1, &vec4.Data[0])
}

func (sp *ShaderProgram) SetUnifMat2(uniformName string, mat2 *gglm.Mat2) {
	loc := gl.GetUniformLocation(sp.ID, gl.Str(uniformName+"\x00"))
	gl.ProgramUniformMatrix2fv(sp.ID, loc, 1, false, &mat2.Data[0][0])
}

func (sp *ShaderProgram) SetUnifMat3(uniformName string, mat3 *gglm.Mat3) {
	loc := gl.GetUniformLocation(sp.ID, gl.Str(uniformName+"\x00"))
	gl.ProgramUniformMatrix3fv(sp.ID, loc, 1, false, &mat3.Data[0][0])
}

func (sp *ShaderProgram) SetUnifMat4(uniformName string, mat4 *gglm.Mat4) {
	loc := gl.GetUniformLocation(sp.ID, gl.Str(uniformName+"\x00"))
	gl.ProgramUniformMatrix4fv(sp.ID, loc, 1, false, &mat4.Data[0][0])
}

func (sp *ShaderProgram) Delete() {
	gl.DeleteProgram(sp.ID)
}
