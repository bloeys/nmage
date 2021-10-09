package shaders

import (
	"fmt"
	"os"
	"strings"

	"github.com/bloeys/go-sdl-engine/logging"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type Shader struct {
	Name string
	ID   uint32
	Type ShaderType
}

//NewShaderFromFile reads a shader from file, creates a new opengl shader and compiles it
func NewShaderFromFile(name, shaderFilePath string, st ShaderType) (Shader, error) {

	b, err := os.ReadFile(shaderFilePath)
	if err != nil {
		return Shader{}, err
	}

	return NewShaderFromString(name, string(b), st)
}

//NewShaderFromString creates a new opengl shader and compiles it
func NewShaderFromString(name, sourceString string, st ShaderType) (Shader, error) {

	glString, freeFunc := gl.Strs(sourceString + "\x00")
	defer freeFunc()

	newShader := Shader{Name: name, Type: st}
	if newShader.ID = gl.CreateShader(st.GLType()); newShader.ID == 0 {
		logging.ErrLog.Panicln("Creating shader failed. ShaderType:", st)
	}

	gl.ShaderSource(newShader.ID, 1, glString, nil)
	gl.CompileShader(newShader.ID)

	return newShader, getShaderCompileError(newShader)
}

func getShaderCompileError(s Shader) error {

	var compileSuccessful int32
	gl.GetShaderiv(s.ID, gl.COMPILE_STATUS, &compileSuccessful)
	if compileSuccessful == gl.TRUE {
		return nil
	}

	//Get the log length and create a string big enough for it and fill it with NULL
	var logLength int32
	gl.GetShaderiv(s.ID, gl.INFO_LOG_LENGTH, &logLength)
	infoLog := gl.Str(strings.Repeat("\x00", int(logLength)))

	//Read the error log and return a go error
	gl.GetShaderInfoLog(s.ID, logLength, nil, infoLog)
	return fmt.Errorf("Shader compilation failed. Compilation log: %s", gl.GoStr(infoLog))
}

func (s *Shader) Delete() {
	gl.DeleteShader(s.ID)
}
