package system

import (
	"github.com/devchucknorris/tank/component"

	"engo.io/ecs"
)

type renderEntity struct {
	*ecs.BasicEntity
	*component.ModelComponent
	*component.TransformComponent
}

type RenderSystem struct {
	entities []renderEntity
}

func (s *RenderSystem) Add(basic *ecs.BasicEntity, model *component.ModelComponent, transform *component.TransformComponent) {
	s.entities = append(s.entities, renderEntity{basic, model, transform})
}

func (s *RenderSystem) Update(dt float32) {
	for _, e := range s.entities {
		modelMatrix := e.CreateModelMatrix()

		e.Shader.Use()
		e.Shader.SetMatrix4fv("model", &modelMatrix[0])

		if e.Texture != nil {
			e.Texture.Bind()
		}
		e.Model.Draw()
	}
}

func (s *RenderSystem) Remove(e ecs.BasicEntity) {
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
