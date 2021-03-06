package ogl

import (
	"fmt"
	"io/ioutil"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	program uint32
}

func (s *Shader) Use() {
	gl.UseProgram(s.program)
}

func (s *Shader) Set1i(name string, value int32) {
	uniform := gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
	gl.Uniform1i(uniform, value)
}

func (s *Shader) Set3f(name string, vec3 mgl32.Vec3) {
	uniform := gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
	gl.Uniform3f(uniform, vec3.X(), vec3.Y(), vec3.Z())
}

func (s *Shader) SetMatrix4fv(name string, value *float32) {
	uniform := gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
	gl.UniformMatrix4fv(uniform, 1, false, value)
}

func (s *Shader) VertexAttribPointer(location uint32, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer) {
	gl.EnableVertexAttribArray(location)
	gl.VertexAttribPointer(location, size, xtype, normalized, stride, pointer)
}

func compile(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("Failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func NewShader(vertex, fragment string) (*Shader, error) {
	s := new(Shader)

	vertexShader, err := compile(vertex, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}

	fragmentShader, err := compile(fragment, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}

	s.program = gl.CreateProgram()
	gl.AttachShader(s.program, vertexShader)
	gl.AttachShader(s.program, fragmentShader)
	gl.LinkProgram(s.program)

	var status int32
	gl.GetProgramiv(s.program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(s.program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(s.program, logLength, nil, gl.Str(log))

		return nil, fmt.Errorf("Failed to ling program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return s, nil
}

func LoadShader(vertexFile, fragmentFile string) (*Shader, error) {
	// Load shader files
	vertexBuf, err := ioutil.ReadFile(vertexFile)
	if err != nil {
		return nil, err
	}
	vertex := string(vertexBuf) + "\x00"

	fragmentBuf, err := ioutil.ReadFile(fragmentFile)
	if err != nil {
		return nil, err
	}
	fragment := string(fragmentBuf) + "\x00"

	// Compile shader
	shader, err := NewShader(vertex, fragment)
	if err != nil {
		return nil, err
	}

	return shader, nil
}
