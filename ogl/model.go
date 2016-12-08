package ogl

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
)

type Model struct {
	vertecies []float32
	index     []uint32
	texCoords []float32
	normals   []float32

	vao, vbo, tbo, nbo, ibo uint32
}

func (m *Model) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.tbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.nbo)
}

func (m *Model) Draw() {
	gl.BindVertexArray(m.vao)
	gl.DrawElements(gl.TRIANGLES, int32(len(m.index)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}

func (m *Model) build(shader *Shader) {
	fmt.Println("Building model...")

	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)

	gl.GenBuffers(1, &m.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.vertecies)*4, gl.Ptr(m.vertecies), gl.STATIC_DRAW)

	shader.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	gl.GenBuffers(1, &m.tbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.tbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.texCoords)*4, gl.Ptr(m.texCoords), gl.STATIC_DRAW)

	shader.VertexAttribPointer(1, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	gl.GenBuffers(1, &m.nbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.nbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.normals)*4, gl.Ptr(m.normals), gl.STATIC_DRAW)

	shader.VertexAttribPointer(2, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	gl.GenBuffers(1, &m.ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.index)*4, gl.Ptr(m.index), gl.STATIC_DRAW)

	gl.BindVertexArray(0)
}

func NewModel(shader *Shader, file string) (*Model, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	temp := strings.Split(string(data), "\n")
	ret := new(Model)

	for i := range temp {
		line := temp[i]
		if strings.HasPrefix(line, "v ") {
			var x, y, z float32
			_, err := fmt.Sscanf(line, "v %f %f %f", &x, &y, &z)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse line %d:\n%v", i, err)
			}

			ret.vertecies = append(ret.vertecies, x, y, z)
		} else if strings.HasPrefix(line, "f ") {
			var v1, t1, n1, v2, t2, n2, v3, t3, n3 uint32
			_, err := fmt.Sscanf(line, "f %d/%d/%d %d/%d/%d %d/%d/%d", &v1, &t1, &n1, &v2, &t2, &n2, &v3, &t3, &n3)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse line %d:\n%v", i, err)
			}

			ret.index = append(ret.index, v1-1, v2-1, v3-1)
		} else if strings.HasPrefix(line, "vt ") {
			var x, y float32
			_, err := fmt.Sscanf(line, "vt %f %f", &x, &y)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse line %d:\n%v", i, err)
			}

			ret.texCoords = append(ret.texCoords, x, y)
		} else if strings.HasPrefix(line, "vn ") {
			var x, y, z float32
			_, err := fmt.Sscanf(line, "vn %f %f %f", &x, &y, &z)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse line %d:\n%v", i, err)
			}

			ret.normals = append(ret.normals, x, y, z)
		}
	}

	ret.build(shader)

	return ret, nil
}
