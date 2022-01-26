package materials

import (
	"github.com/bloeys/gglm/gglm"
	"github.com/bloeys/nmage/buffers"
	"github.com/bloeys/nmage/logging"
	"github.com/bloeys/nmage/shaders"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Material struct {
	Name       string
	ShaderProg shaders.ShaderProgram
}

func (m *Material) Bind() {
	gl.UseProgram(m.ShaderProg.ID)
}

func (m *Material) UnBind() {
	gl.UseProgram(0)
}

func (m *Material) GetAttribLoc(attribName string) int32 {
	return gl.GetAttribLocation(m.ShaderProg.ID, gl.Str(attribName+"\x00"))
}

func (m *Material) SetAttribute(bufObj buffers.Buffer) {

	bufObj.Bind()

	//NOTE: VBOs are only bound at 'VertexAttribPointer', not BindBUffer, so we need to bind the buffer and vao here
	gl.BindBuffer(gl.ARRAY_BUFFER, bufObj.BufID)

	layout := bufObj.GetLayout()
	for i := 0; i < len(layout); i++ {
		gl.EnableVertexAttribArray(uint32(i))
		gl.VertexAttribPointer(uint32(i), layout[i].ElementType.CompCount(), layout[i].ElementType.GLType(), false, bufObj.Stride, gl.PtrOffset(layout[i].Offset))
	}

	bufObj.UnBind()
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (m *Material) EnableAttribute(attribName string) {
	gl.EnableVertexAttribArray(uint32(m.GetAttribLoc(attribName)))
}

func (m *Material) DisableAttribute(attribName string) {
	gl.DisableVertexAttribArray(uint32(m.GetAttribLoc(attribName)))
}

func (m *Material) SetUnifFloat32(uniformName string, val float32) {
	loc := gl.GetUniformLocation(m.ShaderProg.ID, gl.Str(uniformName+"\x00"))
	gl.ProgramUniform1f(m.ShaderProg.ID, loc, val)
}

func (m *Material) SetUnifVec2(uniformName string, vec2 *gglm.Vec2) {
	loc := gl.GetUniformLocation(m.ShaderProg.ID, gl.Str(uniformName+"\x00"))
	gl.ProgramUniform2fv(m.ShaderProg.ID, loc, 1, &vec2.Data[0])
}

func (m *Material) SetUnifVec3(uniformName string, vec3 *gglm.Vec3) {
	loc := gl.GetUniformLocation(m.ShaderProg.ID, gl.Str(uniformName+"\x00"))
	gl.ProgramUniform3fv(m.ShaderProg.ID, loc, 1, &vec3.Data[0])
}

func (m *Material) SetUnifVec4(uniformName string, vec4 *gglm.Vec4) {
	loc := gl.GetUniformLocation(m.ShaderProg.ID, gl.Str(uniformName+"\x00"))
	gl.ProgramUniform4fv(m.ShaderProg.ID, loc, 1, &vec4.Data[0])
}

func (m *Material) SetUnifMat2(uniformName string, mat2 *gglm.Mat2) {
	loc := gl.GetUniformLocation(m.ShaderProg.ID, gl.Str(uniformName+"\x00"))
	gl.ProgramUniformMatrix2fv(m.ShaderProg.ID, loc, 1, false, &mat2.Data[0][0])
}

func (m *Material) SetUnifMat3(uniformName string, mat3 *gglm.Mat3) {
	loc := gl.GetUniformLocation(m.ShaderProg.ID, gl.Str(uniformName+"\x00"))
	gl.ProgramUniformMatrix3fv(m.ShaderProg.ID, loc, 1, false, &mat3.Data[0][0])
}

func (m *Material) SetUnifMat4(uniformName string, mat4 *gglm.Mat4) {
	loc := gl.GetUniformLocation(m.ShaderProg.ID, gl.Str(uniformName+"\x00"))
	gl.ProgramUniformMatrix4fv(m.ShaderProg.ID, loc, 1, false, &mat4.Data[0][0])
}

func (m *Material) Delete() {
	gl.DeleteProgram(m.ShaderProg.ID)
}

func NewMaterial(matName, shaderPath string) *Material {

	shdrProg, err := shaders.NewShaderProgram()
	if err != nil {
		logging.ErrLog.Fatalln("Failed to create new shader program. Err: ", err)
	}

	vertShader, err := shaders.LoadAndCompilerShader(shaderPath+".vert.glsl", shaders.VertexShaderType)
	if err != nil {
		logging.ErrLog.Fatalln("Failed to load and create vertex shader. Err: ", err)
	}

	fragShader, err := shaders.LoadAndCompilerShader(shaderPath+".frag.glsl", shaders.FragmentShaderType)
	if err != nil {
		logging.ErrLog.Fatalln("Failed to load and create fragment shader. Err: ", err)
	}

	shdrProg.AttachShader(vertShader)
	shdrProg.AttachShader(fragShader)
	shdrProg.Link()

	return &Material{Name: matName, ShaderProg: shdrProg}
}
