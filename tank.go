package main

import (
	"fmt"
	"log"
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

const windowWidth = 1920
const windowHeight = 1080

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

	render := system.NewRenderSystem(windowWidth, windowHeight)
	world.AddSystem(&render)

	controller := system.ControllerSystem{Window: window}
	world.AddSystem(&controller)

	// Load shader
	shader, err := ogl.LoadShader("data/vertex.glsl", "data/fragment.glsl")
	if err != nil {
		log.Fatalln("Failed to load shader", err)
	}
	shader.Use()

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
	tank.ModelComponent = component.ModelComponent{Shader: shader, Model: tankModel, Texture: texture, CastShadow: true}
	tank.TransformComponent = component.TransformComponent{
		X: 0.0, Y: 0.0, Z: -5.0,
		ScaleX: 0.4, ScaleY: 0.4, ScaleZ: 0.4,
		RotationX: mgl32.DegToRad(-90), RotationY: 0.0, RotationZ: 0.0}
	tank.MovementComponent = component.MovementComponent{MoveSpeed: 2, RotationSpeed: 100}

	render.Add(&tank.BasicEntity, &tank.ModelComponent, &tank.TransformComponent)
	controller.Add(&tank.BasicEntity, &tank.TransformComponent, &tank.MovementComponent)

	// Create ground
	ground := entity.Obstacle{BasicEntity: ecs.NewBasic()}
	ground.ModelComponent = component.ModelComponent{Shader: shader, Model: groundModel, Texture: groundTexture, CastShadow: true}
	ground.TransformComponent = component.TransformComponent{
		X: 0.0, Y: 0.0, Z: 0.0,
		ScaleX: 1, ScaleY: 1, ScaleZ: 1,
		RotationX: 0.0, RotationY: 0.0, RotationZ: 0.0}

	render.Add(&ground.BasicEntity, &ground.ModelComponent, &ground.TransformComponent)

	for !window.ShouldClose() {
		time := glfw.GetTime()
		elapsed := float32(time - previousTime)
		previousTime = time

		world.Update(elapsed)

		window.Update()
	}
}
