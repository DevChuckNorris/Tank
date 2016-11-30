package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
)

type Model struct {
	vertecies []float32
	index     []uint32
	texCoords []float32

	vao, vbo, tbo, ibo uint32
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

	shader.VertexAttribPointer("vert", 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	gl.GenBuffers(1, &m.tbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.tbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.texCoords)*4, gl.Ptr(m.texCoords), gl.STATIC_DRAW)

	shader.VertexAttribPointer("vertTexCoord", 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

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
			details := strings.Split(line, " ")
			if len(details) != 4 {
				return nil, fmt.Errorf("A: Failed to parse line %d", i)
			}

			x, err := strconv.ParseFloat(details[1], 32)
			if err != nil {
				return nil, fmt.Errorf("B: Failed to parse line %d", i)
			}
			y, err := strconv.ParseFloat(details[2], 32)
			if err != nil {
				return nil, fmt.Errorf("C: Failed to parse line %d", i)
			}

			details[3] = strings.Replace(details[3], "\r", "", -1)
			z, err := strconv.ParseFloat(details[3], 32)
			if err != nil {
				return nil, fmt.Errorf("D: Failed to parse line %d %v", i, err)
			}

			ret.vertecies = append(ret.vertecies, float32(x), float32(y), float32(z))
		} else if strings.HasPrefix(line, "f ") {
			details := strings.Split(line, " ")
			if len(details) != 4 {
				return nil, fmt.Errorf("E: Failed to parse line %d", i)
			}

			for x := 1; x < 4; x++ {
				part := strings.Split(details[x], "/")
				if len(part) != 3 {
					return nil, fmt.Errorf("F: Failed to parse line %d", i)
				}
				a, err := strconv.ParseUint(part[0], 10, 32)
				if err != nil {
					return nil, fmt.Errorf("G: Failed to parse line %d", i)
				}
				ret.index = append(ret.index, uint32(a-1))
			}
		} else if strings.HasPrefix(line, "vt ") {
			details := strings.Split(line, " ")
			if len(details) != 3 {
				return nil, fmt.Errorf("H: Failed to parse line %d", i)
			}

			x, err := strconv.ParseFloat(details[1], 32)
			if err != nil {
				return nil, fmt.Errorf("I: Failed to parse line %d", i)
			}
			details[2] = strings.Replace(details[2], "\r", "", -1)
			y, err := strconv.ParseFloat(details[2], 32)
			if err != nil {
				return nil, fmt.Errorf("J: Failed to parse line %d", i)
			}

			ret.texCoords = append(ret.texCoords, float32(x), float32(y))
		}
	}

	for i := range ret.vertecies {
		fmt.Printf("%v ", ret.vertecies[i])
		if (i+1)%3 == 0 {
			fmt.Println()
		}
	}

	fmt.Println("--------------------------------------------------")

	for i := range ret.index {
		index := ret.index[i]
		vertex := []float32{ret.vertecies[index*3], ret.vertecies[index*3+1], ret.vertecies[index*3+2]}
		fmt.Printf("%v(%v %v %v) ", index, vertex[0], vertex[1], vertex[2])
		if (i+1)%3 == 0 {
			fmt.Println()
		}
	}

	ret.build(shader)

	return ret, nil
}
