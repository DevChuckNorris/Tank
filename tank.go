package main

import (
	"fmt"
	"log"
	"math"
	"runtime"

	"github.com/devchucknorris/tank/ogl"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 800
const windowHeight = 600

var textures [5]*ogl.Image

var position mgl32.Vec3

func init() {
	runtime.LockOSThread()
}

func main() {
	window, err := ogl.CreateWindow(windowWidth, windowHeight, "Cube")
	if err != nil {
		panic(err)
	}
	defer window.Close()

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Load shader
	shader, err := ogl.LoadShader("data/vertex.glsl", "data/fragment.glsl")
	if err != nil {
		log.Fatalln("Failed to load shader", err)
	}
	shader.Use()

	textureId := 0

	// Load image
	texture, err := ogl.NewImage("data/blue_tank.png")
	if err != nil {
		log.Fatalln("Failed to load image", err)
	}
	textures[0] = texture

	texture2, err := ogl.NewImage("data/green_tank.png")
	if err != nil {
		log.Fatalln("Failed to load image", err)
	}
	textures[1] = texture2

	texture3, err := ogl.NewImage("data/purple_tank.png")
	if err != nil {
		log.Fatalln("Failed to load image", err)
	}
	textures[2] = texture3

	texture4, err := ogl.NewImage("data/red_tank.png")
	if err != nil {
		log.Fatalln("Failed to load image", err)
	}
	textures[3] = texture4

	texture5, err := ogl.NewImage("data/rogue_tank.png")
	if err != nil {
		log.Fatalln("Failed to load image", err)
	}
	textures[4] = texture5

	// Projection etc
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 100.0)
	camera := mgl32.LookAtV(mgl32.Vec3{2, 4, 10}, mgl32.Vec3{2, 0, 0}, mgl32.Vec3{0, 1, 0})
	model := mgl32.Ident4()

	// Bind to shader
	shader.SetMatrix4fv("projection", &projection[0])
	shader.SetMatrix4fv("camera", &camera[0])
	shader.SetMatrix4fv("model", &model[0])
	shader.Set3f("light", mgl32.Vec3{4, 4, 4})

	// Load model
	tank, err := ogl.NewModel(shader, "data/tank.obj")
	if err != nil {
		log.Fatalln("Failed to load model", err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	angle := 0.0
	previousTime := glfw.GetTime()

	scale := mgl32.Scale3D(0.4, 0.4, 0.4)
	lastSpace := window.GetKey(glfw.KeySpace)

	angleCorrection := 0.0 //math.Pi * 90 / 180.0

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		//angle += elapsed
		rotation := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
		rotation = rotation.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(-90), mgl32.Vec3{1, 0, 0}))
		model = mgl32.Translate3D(position.X(), position.Y(), position.Z()).Mul4(scale.Mul4(rotation))

		space := window.GetKey(glfw.KeySpace)
		if space == glfw.Press && lastSpace == glfw.Release {
			textureId++
			if textureId > len(textures)-1 {
				textureId = 0
			}
		}
		lastSpace = space

		if window.GetKey(glfw.KeyA) == glfw.Press {
			angle += float64(mgl32.DegToRad(100)) * elapsed
		}
		if window.GetKey(glfw.KeyD) == glfw.Press {
			angle -= float64(mgl32.DegToRad(100)) * elapsed
		}
		if window.GetKey(glfw.KeyW) == glfw.Press {
			fmt.Println("W")
			position = position.Add(mgl32.Vec3{
				float32(2 * elapsed * math.Cos(angle-angleCorrection)),
				0,
				float32(-2 * elapsed * math.Sin(angle-angleCorrection))})
		}

		shader.Use()
		shader.SetMatrix4fv("model", &model[0])

		textures[textureId].Bind()
		tank.Draw()

		window.Update()
	}
}
