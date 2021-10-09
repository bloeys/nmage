package shaders

import (
	"fmt"
	"log"
	"strings"

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
		log.Fatalln("Creating OpenGL program failed")
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

//Delete deletes all shaders and then deletes the program
func (p *Program) Delete() {
	for _, v := range p.Shaders {
		v.Delete()
	}

	gl.DeleteProgram(p.ID)
}
