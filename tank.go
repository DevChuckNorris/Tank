package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
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
	defer window.Close()

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Load shader
	shader, err := LoadShader("vertex.glsl", "fragment.glsl")
	if err != nil {
		log.Fatalln("Failed to load shader", err)
	}
	shader.Use()

	// Load image
	image, err := NewImage("blue_tank.png")
	if err != nil {
		log.Fatalln("Failed to load image", err)
	}

	// Projection etc
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	model := mgl32.Ident4()

	// Bind to shader
	shader.SetMatrix4fv("projection", &projection[0])
	shader.SetMatrix4fv("camera", &camera[0])
	shader.SetMatrix4fv("model", &model[0])

	// Load model
	tank, err := NewModel(shader, "tank.obj")
	if err != nil {
		log.Fatalln("Failed to load model", err)
	}

	// Bind data to shader

	//shader.VertexAttribPointer("vertTexCoord",  2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	angle := 0.0
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		angle += elapsed
		scale := mgl32.Scale3D(0.4, 0.4, 0.4)
		rotation := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
		rotation = rotation.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(-90), mgl32.Vec3{1, 0, 0}))
		model = scale.Mul4(rotation)

		shader.Use()
		shader.SetMatrix4fv("model", &model[0])

		image.Bind()
		//gl.BindVertexArray(vao)
		//gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
		tank.Draw()

		window.Update()
	}
}

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	-1.0, -1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,

	// Top
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 1.0,

	// Left
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,

	// Right
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
}
