package main

import (
    "io/ioutil"
    "strings"
    "fmt"
    "strconv"

    "github.com/go-gl/gl/v3.2-core/gl"
)

type Model struct {
    vertecies []float32
    index []uint32

    vao, vbo, ibo uint32
}

func (m *Model) Draw() {
    gl.BindVertexArray(m.vao)
    gl.DrawElements(gl.TRIANGLES, 9, gl.UNSIGNED_INT, gl.PtrOffset(0))
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
    ret := new (Model)

    for i := range temp {
        line := temp[i]
        if strings.HasPrefix(line, "v ") {
            details := strings.Split(line, " ")
            if len(details) != 4 {
                return nil, fmt.Errorf("Failed to parse line %d", i)
            }

            x, err := strconv.ParseFloat(details[1], 32)
            if err != nil {
                return nil, fmt.Errorf("Failed to parse line %d", i)
            }
            y, err := strconv.ParseFloat(details[2], 32)
            if err != nil {
                return nil, fmt.Errorf("Failed to parse line %d", i)
            }
            z, err := strconv.ParseFloat(details[3], 32)
            if err != nil {
                return nil, fmt.Errorf("Failed to parse line %d", i)
            }

            ret.vertecies = append(ret.vertecies, float32(x), float32(y), float32(z))
        } else if strings.HasPrefix(line, "f ") {
            details := strings.Split(line, " ")
            if len(details) != 4 {
                return nil, fmt.Errorf("Failed to parse line %d", i)
            }

            for x := 1; x < 4; x++ {
                part := strings.Split(details[x], "/")
                if len(part) != 3 {
                    return nil, fmt.Errorf("Failed to parse line %d", i)
                }
                a, err := strconv.ParseUint(part[0], 10, 32)
                if err != nil {
                    return nil, fmt.Errorf("Failed to parse line %d", i)
                }
                fmt.Printf("AddF %d\n", uint32(a))
                ret.index = append(ret.index, uint32(a))
        	}
        }
    }

    for i := range ret.vertecies {
        fmt.Println(ret.vertecies[i])
    }

    ret.build(shader)

    return ret, nil
}
