package shaders

import (
	"fmt"
	"strings"

	"github.com/bloeys/go-sdl-engine/logging"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type Program struct {
	Name    string
	ID      uint32
	Shaders []Shader
}

func NewProgram(name string) Program {

	p := Program{Name: name}
	p.Shaders = make([]Shader, 0)

	if p.ID = gl.CreateProgram(); p.ID == 0 {
		logging.ErrLog.Fatalln("Creating OpenGL program failed")
	}

	return p
}

//AttachShader adds the shader to list of shaders and attaches it in opengl
func (p *Program) AttachShader(s Shader) {

	p.Shaders = append(p.Shaders, s)
	gl.AttachShader(p.ID, s.ID)
}

//DetachShader removes the shader from the list of shaders and detaches it in opengl
func (p *Program) DetachShader(s Shader) {

	//To remove a shader we move the last shader to its place, then shrink the slice by one
	for i := 0; i < len(p.Shaders); i++ {

		if p.Shaders[i].ID != s.ID {
			continue
		}

		gl.DetachShader(p.ID, s.ID)

		p.Shaders[i] = p.Shaders[len(p.Shaders)-1]
		p.Shaders = p.Shaders[:len(p.Shaders)-1]
		return
	}
}

func (p *Program) Link() error {

	gl.LinkProgram(p.ID)
	return getProgramLinkError(*p)
}

func getProgramLinkError(p Program) error {

	var linkSuccessful int32
	gl.GetProgramiv(p.ID, gl.LINK_STATUS, &linkSuccessful)
	if linkSuccessful == gl.TRUE {
		return nil
	}

	//Get the log length and create a string big enough for it and fill it with NULL
	var logLength int32
	gl.GetProgramiv(p.ID, gl.INFO_LOG_LENGTH, &logLength)
	infoLog := gl.Str(strings.Repeat("\x00", int(logLength)))

	//Read the error log and return a go error
	gl.GetProgramInfoLog(p.ID, logLength, nil, infoLog)
	return fmt.Errorf("Program linking failed. Linking log: %s", gl.GoStr(infoLog))
}

func (p *Program) Use() {
	gl.UseProgram(p.ID)
}

func (p *Program) GetUniformLocation(name string) int32 {
	return gl.GetUniformLocation(p.ID, gl.Str(name+"\x00"))
}

//SetUniformF32 handles setting uniform values of 1-4 floats.
//Returns false if len(floats) is <1 or >4, or if the uniform was not found.
//Uniforms aren't found if it doesn't exist or was not used in the shader
func (p *Program) SetUniformF32(name string, floats ...float32) bool {

	loc := p.GetUniformLocation(name)
	if loc == 0 {
		logging.WarnLog.Printf(
			"Uniform with name '%s' was not found. "+
				"This is either because it doesn't exist or isn't used in the shader",
			name)
		return false
	}

	switch len(floats) {
	case 1:
		gl.Uniform1f(loc, floats[0])
	case 2:
		gl.Uniform2f(loc, floats[0], floats[1])
	case 3:
		gl.Uniform3f(loc, floats[0], floats[1], floats[2])
	case 4:
		gl.Uniform4f(loc, floats[0], floats[1], floats[2], floats[3])
	default:
		logging.ErrLog.Println("Invalid input size in SetUniformF32. Size must be 1-4 but got", len(floats))
		return false
	}

	return true
}

//SetUniformI32 handles setting uniform values of 1-4 ints.
//Returns false if len(ints) is <1 or >4, or if the uniform was not found.
//Uniforms aren't found if it doesn't exist or was not used in the shader
func (p *Program) SetUniformI32(name string, ints ...int32) bool {

	loc := p.GetUniformLocation(name)
	if loc == 0 {
		logging.WarnLog.Printf(
			"Uniform with name '%s' was not found. "+
				"This is either because it doesn't exist or isn't used in the shader",
			name)
		return false
	}

	switch len(ints) {
	case 1:
		gl.Uniform1i(loc, ints[0])
	case 2:
		gl.Uniform2i(loc, ints[0], ints[1])
	case 3:
		gl.Uniform3i(loc, ints[0], ints[1], ints[2])
	case 4:
		gl.Uniform4i(loc, ints[0], ints[1], ints[2], ints[3])
	default:
		logging.ErrLog.Println("Invalid input size in SetUniformI32. Size must be 1-4 but got", len(ints))
		return false
	}

	return true
}

//Delete deletes all shaders and then deletes the program
func (p *Program) Delete() {
	for _, v := range p.Shaders {
		v.Delete()
	}

	gl.DeleteProgram(p.ID)
}
