package main

import (
	"fmt"
	"log"
	"math"
	"runtime"

	"engo.io/ecs"
	"github.com/devchucknorris/tank/component"
	"github.com/devchucknorris/tank/entity"
	"github.com/devchucknorris/tank/ogl"
	"github.com/devchucknorris/tank/system"
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
	window, err := ogl.CreateWindow(windowWidth, windowHeight, "Cube")
	if err != nil {
		panic(err)
	}
	defer window.Close()

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	previousTime := glfw.GetTime()

	world := ecs.World{}

	render := system.RenderSystem{}
	world.AddSystem(&render)

	controller := system.ControllerSystem{Window: window}
	world.AddSystem(&controller)

	// Load shader
	shader, err := ogl.LoadShader("data/vertex.glsl", "data/fragment.glsl")
	if err != nil {
		log.Fatalln("Failed to load shader", err)
	}
	shader.Use()

	// Projection etc
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 100.0)
	camera := mgl32.LookAtV(mgl32.Vec3{2, 4, 10}, mgl32.Vec3{2, 0, 0}, mgl32.Vec3{0, 1, 0})

	// Bind to shader
	shader.SetMatrix4fv("projection", &projection[0])
	shader.SetMatrix4fv("camera", &camera[0])

	// Load model
	tankModel, err := ogl.NewModel(shader, "data/tank.obj")
	if err != nil {
		log.Fatalln("Failed to load model", err)
	}
	groundModel := ogl.NewBox(10, 0, 10, 5, shader)

	// Load image
	texture, err := ogl.NewImage("data/blue_tank.png")
	if err != nil {
		log.Fatalln("Failed to load image", err)
	}
	groundTexture, err := ogl.NewImage("data/std_ground.png")
	if err != nil {
		log.Fatalln("Failed to load image", err)
	}

	// Create tank
	tank := entity.Tank{BasicEntity: ecs.NewBasic()}
	tank.ModelComponent = component.ModelComponent{Shader: shader, Model: tankModel, Texture: texture}
	tank.TransformComponent = component.TransformComponent{
		X: 0.0, Y: 0.0, Z: -5.0,
		ScaleX: 0.4, ScaleY: 0.4, ScaleZ: 0.4,
		RotationX: mgl32.DegToRad(-90), RotationY: 0.0, RotationZ: 0.0}
	tank.MovementComponent = component.MovementComponent{MoveSpeed: 2, RotationSpeed: 100}

	render.Add(&tank.BasicEntity, &tank.ModelComponent, &tank.TransformComponent)
	controller.Add(&tank.BasicEntity, &tank.TransformComponent, &tank.MovementComponent)

	// Create ground
	ground := entity.Obstacle{BasicEntity: ecs.NewBasic()}
	ground.ModelComponent = component.ModelComponent{Shader: shader, Model: groundModel, Texture: groundTexture}
	ground.TransformComponent = component.TransformComponent{
		X: 0.0, Y: 0.0, Z: 0.0,
		ScaleX: 1, ScaleY: 1, ScaleZ: 1,
		RotationX: 0.0, RotationY: 0.0, RotationZ: 0.0}

	render.Add(&ground.BasicEntity, &ground.ModelComponent, &ground.TransformComponent)

	sunAngle := 45.0

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		lightPos := mgl32.Vec3{
			float32(math.Cos(sunAngle*math.Pi/180.0)) * 70,
			float32(math.Sin(sunAngle*math.Pi/180.0)) * 70,
			0.0}
		lightPos = lightPos.Normalize().Mul(-1)
		shader.Set3f("light", lightPos)

		time := glfw.GetTime()
		elapsed := float32(time - previousTime)
		previousTime = time

		world.Update(elapsed)

		window.Update()
	}
}
