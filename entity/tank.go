package entity

import (
	"engo.io/ecs"
	"github.com/devchucknorris/tank/component"
)

type Tank struct {
	ecs.BasicEntity

	component.ModelComponent
	component.TransformComponent
}
