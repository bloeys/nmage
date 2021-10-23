package shaders

import (
	"errors"
	"os"
	"strings"

	"github.com/bloeys/go-sdl-engine/logging"
	"github.com/go-gl/gl/v4.6-compatibility/gl"
)

type Shader struct {
	ID         uint32
	ShaderType ShaderType
}

func (s Shader) Delete() {
	gl.DeleteShader(s.ID)
}

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

func (sp *ShaderProgram) Delete() {
	gl.DeleteProgram(sp.ID)
}

func NewShaderProgram() (ShaderProgram, error) {

	id := gl.CreateProgram()
	if id == 0 {
		return ShaderProgram{}, errors.New("failed to create shader program")
	}

	return ShaderProgram{ID: id}, nil
}

func LoadAndCompilerShader(shaderPath string, shaderType ShaderType) (Shader, error) {

	shaderSource, err := os.ReadFile(shaderPath)
	if err != nil {
		logging.ErrLog.Println("Failed to read shader. Err: ", err)
		return Shader{}, err
	}

	shaderID := gl.CreateShader(uint32(shaderType))
	if shaderID == 0 {
		logging.ErrLog.Println("Failed to create shader.")
		return Shader{}, errors.New("failed to create shader")
	}

	//Load shader source and compile
	shaderCStr, shaderFree := gl.Strs(string(shaderSource) + "\x00")
	defer shaderFree()
	gl.ShaderSource(shaderID, 1, shaderCStr, nil)

	gl.CompileShader(shaderID)
	if err := getShaderCompileErrors(shaderID); err != nil {
		gl.DeleteShader(shaderID)
		return Shader{}, err
	}

	return Shader{ID: shaderID, ShaderType: shaderType}, nil
}

func getShaderCompileErrors(shaderID uint32) error {

	var compiledSuccessfully int32
	gl.GetShaderiv(shaderID, gl.COMPILE_STATUS, &compiledSuccessfully)
	if compiledSuccessfully == gl.TRUE {
		return nil
	}

	var logLength int32
	gl.GetShaderiv(shaderID, gl.INFO_LOG_LENGTH, &logLength)

	log := gl.Str(strings.Repeat("\x00", int(logLength)))
	gl.GetShaderInfoLog(shaderID, logLength, nil, log)

	errMsg := gl.GoStr(log)
	logging.ErrLog.Println("Compilation of shader with id ", shaderID, " failed. Err: ", errMsg)
	return errors.New(errMsg)
}
