package system

import (
	"log"

	"engo.io/ecs"
	"github.com/devchucknorris/tank/component"
	"github.com/devchucknorris/tank/ogl"
	"github.com/go-gl/gl/v3.2-core/gl"
)

type renderEntity struct {
	*ecs.BasicEntity
	*component.ModelComponent
	*component.TransformComponent
}

type RenderSystem struct {
	entities []renderEntity

	shadowBuffer  uint32
	shadowTexture uint32
	width         int32
	height        int32
	shadowShader  *ogl.Shader
}

func NewRenderSystem(width, height int32) RenderSystem {
	ret := RenderSystem{width: width, height: height}

	shader, err := ogl.LoadShader("data/shadow_vertex.glsl", "data/shadow_fragment.glsl")
	if err != nil {
		log.Fatalln("Failed to load shader", err)
	}
	ret.shadowShader = shader

	gl.GenFramebuffers(1, &ret.shadowBuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, ret.shadowBuffer)

	gl.GenTextures(1, &ret.shadowTexture)
	gl.BindTexture(gl.TEXTURE_2D, ret.shadowTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT16, 1024, 1024, 0, gl.DEPTH_COMPONENT, gl.FLOAT, gl.PtrOffset(0))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, ret.shadowTexture, 0)
	gl.DrawBuffer(gl.NONE)

	return ret
}

func (s *RenderSystem) Add(basic *ecs.BasicEntity, model *component.ModelComponent, transform *component.TransformComponent) {
	s.entities = append(s.entities, renderEntity{basic, model, transform})
}

func (s *RenderSystem) Update(dt float32) {
	// Step One: Render shadow map
	gl.BindFramebuffer(gl.FRAMEBUFFER, s.shadowBuffer)
	gl.Viewport(0, 0, 1024, 1024)

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	s.shadowShader.Use()
	s.shadowShader.VertexAttribPointer("vertexPosition_modelspace", 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	for _, e := range s.entities {
		if e.ModelComponent.CastShadow {
			e.Model.Draw()
		}
	}

	// Step Two: Render objects
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Viewport(0, 0, s.width, s.height)

	gl.Disable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

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
