package system

import (
	"math"

	"engo.io/ecs"
	"github.com/devchucknorris/tank/component"
	"github.com/devchucknorris/tank/ogl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type controllerEntity struct {
	*ecs.BasicEntity
	*component.TransformComponent
	*component.MovementComponent
}

type ControllerSystem struct {
	Window   *ogl.GameWindow
	entities []controllerEntity
}

func (s *ControllerSystem) Add(basic *ecs.BasicEntity, transform *component.TransformComponent, movement *component.MovementComponent) {
	s.entities = append(s.entities, controllerEntity{basic, transform, movement})
}

func (s *ControllerSystem) Update(dt float32) {
	for _, e := range s.entities {
		if e.MovementComponent != nil {
			// this player
			if s.Window.GetKey(glfw.KeyA) == glfw.Press {
				e.TransformComponent.RotationZ += mgl32.DegToRad(e.MovementComponent.RotationSpeed) * dt
			}
			if s.Window.GetKey(glfw.KeyD) == glfw.Press {
				e.TransformComponent.RotationZ -= mgl32.DegToRad(e.MovementComponent.RotationSpeed) * dt
			}
			if s.Window.GetKey(glfw.KeyW) == glfw.Press {
				angle := float64(e.TransformComponent.RotationZ)
				e.TransformComponent.X += float32(float64(e.MoveSpeed*dt) * math.Cos(angle))
				e.TransformComponent.Z += float32(float64(-e.MoveSpeed*dt) * math.Sin(angle))
			}
		}
	}
}

func (s *ControllerSystem) Remove(e ecs.BasicEntity) {
	var delete int = -1
	for index, entity := range s.entities {
		if entity.ID() == e.ID() {
			delete = index
			break
		}
	}

	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}
