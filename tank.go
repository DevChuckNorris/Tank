package main

import (
    "runtime"
    "log"
    "fmt"

    "github.com/go-gl/gl/v3.2-core/gl"
)

const windowWidth = 800
const windowHeight = 600

func init() {
    runtime.LockOSThread()
}

func main() {
    window, err := CreateWindow(windowWidth, windowHeight, "Cube")
    if err != nil {
        panic(err)
    }

    version := gl.GoStr(gl.GetString(gl.VERSION))
    fmt.Println("OpenGL version", version)

    // Load shader
    shader, err := LoadShader("vertex.glsl", "fragment.glsl")
    if err != nil {
        log.Fatalln("Failed to load shader", err)
    }
    shader.Use()

    for !window.ShouldClose() {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        window.Update()
    }
}
