package component

import (
	"github.com/devchucknorris/tank/ogl"
	"github.com/go-gl/mathgl/mgl32"
)

type ModelComponent struct {
	Texture *ogl.Image
	Model   *ogl.Model
	Shader  *ogl.Shader
}

type TransformComponent struct {
	X         float32
	Y         float32
	Z         float32
	ScaleX    float32
	ScaleY    float32
	ScaleZ    float32
	RotationX float32
	RotationY float32
	RotationZ float32
}

func (transform *TransformComponent) CreateModelMatrix() mgl32.Mat4 {
	rotation := mgl32.HomogRotate3DX(transform.RotationX)
	rotation = rotation.Mul4(mgl32.HomogRotate3DY(transform.RotationY))
	rotation = rotation.Mul4(mgl32.HomogRotate3DZ(transform.RotationZ))

	scale := mgl32.Scale3D(transform.ScaleX, transform.ScaleY, transform.ScaleZ)

	return mgl32.Translate3D(transform.X, transform.Y, transform.Z).Mul4(scale.Mul4(rotation))
}
