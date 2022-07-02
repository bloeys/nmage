package shaders

import (
	"bytes"
	"errors"
	"os"
	"strings"

	"github.com/bloeys/nmage/logging"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	ID         uint32
	ShaderType ShaderType
}

func (s Shader) Delete() {
	gl.DeleteShader(s.ID)
}

func NewShaderProgram() (ShaderProgram, error) {

	id := gl.CreateProgram()
	if id == 0 {
		return ShaderProgram{}, errors.New("failed to create shader program")
	}

	return ShaderProgram{ID: id}, nil
}

func LoadAndCompileCombinedShader(shaderPath string) (ShaderProgram, error) {

	combinedSource, err := os.ReadFile(shaderPath)
	if err != nil {
		logging.ErrLog.Println("Failed to read shader. Err: ", err)
		return ShaderProgram{}, err
	}

	return LoadAndCompileCombinedShaderSrc(combinedSource)

}
func LoadAndCompileCombinedShaderSrc(shaderSrc []byte) (ShaderProgram, error) {

	shaderSources := bytes.Split(shaderSrc, []byte("//shader:"))
	if len(shaderSources) == 1 {
		return ShaderProgram{}, errors.New("failed to read combined shader. Did not find '//shader:vertex' or '//shader:fragment'")
	}

	shdrProg, err := NewShaderProgram()
	if err != nil {
		return ShaderProgram{}, errors.New("failed to create new shader program. Err: " + err.Error())
	}

	loadedShdrCount := 0
	for i := 0; i < len(shaderSources); i++ {

		src := shaderSources[i]

		//This can happen when the shader type is at the start of the file
		if len(bytes.TrimSpace(src)) == 0 {
			continue
		}

		var shdrType ShaderType
		if bytes.HasPrefix(src, []byte("vertex")) {
			src = src[6:]
			shdrType = VertexShaderType
		} else if bytes.HasPrefix(src, []byte("fragment")) {
			src = src[8:]
			shdrType = FragmentShaderType
		} else {
			return ShaderProgram{}, errors.New("unknown shader type. Must be '//shader:vertex' or '//shader:fragment'")
		}

		shdr, err := CompileShaderOfType(src, shdrType)
		if err != nil {
			return ShaderProgram{}, err
		}

		loadedShdrCount++
		shdrProg.AttachShader(shdr)
	}

	if loadedShdrCount == 0 {
		return ShaderProgram{}, errors.New("no valid shaders found. Please put '//shader:vertex' or '//shader:fragment' before your shaders")
	}

	shdrProg.Link()
	return shdrProg, nil
}

func CompileShaderOfType(shaderSource []byte, shaderType ShaderType) (Shader, error) {

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
