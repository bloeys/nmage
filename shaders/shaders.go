package shaders

import (
	"os"
	"strings"

	"github.com/bloeys/go-sdl-engine/logging"
	"github.com/go-gl/gl/v4.6-compatibility/gl"
)

type ShaderProgram struct {
	ID           uint32
	VertShaderID uint32
	FragShaderID uint32
}

func LoadShaders() {

	vertShaderText, err := os.ReadFile("./res/shaders/simple.vert.glsl")
	if err != nil {
		logging.ErrLog.Fatalln("Failed to read vertex shader. Err: ", err)
	}

	fragShaderText, err := os.ReadFile("./res/shaders/simple.frag.glsl")
	if err != nil {
		logging.ErrLog.Fatalln("Failed to read fragment shader. Err: ", err)
	}

	shader := &ShaderProgram{}
	shader.ID = gl.CreateProgram()
	if shader.ID == 0 {
		logging.ErrLog.Fatalln("Failed to create shader program")
	}

	shader.VertShaderID = gl.CreateShader(gl.VERTEX_SHADER)
	shader.FragShaderID = gl.CreateShader(gl.FRAGMENT_SHADER)

	vertexCStr, vertFree := gl.Strs(string(vertShaderText) + "\x00")
	defer vertFree()
	gl.ShaderSource(shader.VertShaderID, 1, vertexCStr, nil)

	fragCStr, fragFree := gl.Strs(string(fragShaderText) + "\x00")
	defer fragFree()
	gl.ShaderSource(shader.FragShaderID, 1, fragCStr, nil)

	gl.CompileShader(shader.VertShaderID)
	getShaderCompileErrors(shader.VertShaderID)

	gl.CompileShader(shader.FragShaderID)
	getShaderCompileErrors(shader.FragShaderID)

	gl.AttachShader(shader.ID, shader.VertShaderID)
	gl.AttachShader(shader.ID, shader.FragShaderID)
	gl.LinkProgram(shader.ID)
}

func getShaderCompileErrors(shaderID uint32) {

	var compiledSuccessfully int32
	gl.GetShaderiv(shaderID, gl.COMPILE_STATUS, &compiledSuccessfully)
	if compiledSuccessfully == gl.TRUE {
		return
	}

	var logLength int32
	gl.GetShaderiv(shaderID, gl.INFO_LOG_LENGTH, &logLength)

	log := gl.Str(strings.Repeat("\x00", int(logLength)))
	gl.GetShaderInfoLog(shaderID, logLength, nil, log)

	errMsg := gl.GoStr(log)
	println("Compilation of shader with id ", shaderID, " failed. Err: ", errMsg)
}
