package system

import (
	"log"

	"engo.io/ecs"
	"github.com/devchucknorris/tank/component"
	"github.com/devchucknorris/tank/ogl"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
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

	debugVertecies []float32
	debugShader    *ogl.Shader
	debugVAO       uint32
	debugVBO       uint32

	projection mgl32.Mat4
	view       mgl32.Mat4
}

func NewRenderSystem(width, height int32) RenderSystem {
	ret := RenderSystem{width: width, height: height}

	// Projection etc
	ret.projection = mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/float32(height), 0.1, 100.0)
	ret.view = mgl32.LookAtV(mgl32.Vec3{0, 7, 20}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})

	// Bind to shader
	//shader.SetMatrix4fv("projection", &projection[0])
	//shader.SetMatrix4fv("camera", &camera[0])

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
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_COMPARE_FUNC, gl.LEQUAL)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_COMPARE_MODE, gl.COMPARE_REF_TO_TEXTURE)

	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, ret.shadowTexture, 0)
	gl.DrawBuffer(gl.NONE)

	ret.createDebug()

	return ret
}

func (s *RenderSystem) createDebug() {
	// Load debug shader
	shader, err := ogl.LoadShader("data/debug_vertex.glsl", "data/debug_fragment.glsl")
	if err != nil {
		log.Fatalln("Failed to load shader", err)
	}
	s.debugShader = shader

	// Create quad vertecies
	s.debugVertecies = []float32{-1.0, -1.0, 0.0,
		1.0, -1.0, 0.0,
		-1.0, 1.0, 0.0,
		-1.0, 1.0, 0.0,
		1.0, -1.0, 0.0,
		1.0, 1.0, 0.0}

	gl.GenVertexArrays(1, &s.debugVAO)
	gl.BindVertexArray(s.debugVAO)

	gl.GenBuffers(1, &s.debugVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.debugVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(s.debugVertecies)*4, gl.Ptr(s.debugVertecies), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
}

func (s *RenderSystem) Add(basic *ecs.BasicEntity, model *component.ModelComponent, transform *component.TransformComponent) {
	s.entities = append(s.entities, renderEntity{basic, model, transform})
}

func (s *RenderSystem) Update(dt float32) {
	// Step Zero: Calculate light position
	/*sunAngle := 45.0

	lightPos := mgl32.Vec3{
		float32(math.Cos(sunAngle*math.Pi/180.0)) * 70,
		float32(math.Sin(sunAngle*math.Pi/180.0)) * 70,
		0.0}
	lightPos = lightPos.Normalize().Mul(-1)*/
	lightPos := mgl32.Vec3{.5, 2, 2}

	// Step One: Render shadow map
	gl.BindFramebuffer(gl.FRAMEBUFFER, s.shadowBuffer)
	gl.Viewport(0, 0, 1024, 1024)

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	s.shadowShader.Use()

	depthProjectionMatrix := mgl32.Ortho(-10, 10, -10, 10, -10, 20)
	depthViewMatrix := mgl32.LookAt(lightPos.X(), lightPos.Y(), lightPos.Z(), 0, 0, 0, 0, 1, 0)

	for _, e := range s.entities {
		if e.ModelComponent.CastShadow {
			depthModelMatrix := e.TransformComponent.CreateModelMatrix()

			depthMVP := depthProjectionMatrix.Mul4(depthViewMatrix.Mul4(depthModelMatrix))

			s.shadowShader.SetMatrix4fv("depthMVP", &depthMVP[0])
			e.Model.Draw()
		}
	}

	// Step Two: Render objects
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Viewport(0, 0, s.width, s.height)

	gl.Disable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, s.shadowTexture)

	biasMatrix := mgl32.Mat4FromRows(mgl32.Vec4{0.5, 0.0, 0.0, 0.0}, mgl32.Vec4{0.0, 0.5, 0.0, 0.0}, mgl32.Vec4{0.0, 0.0, 0.5, 0.0}, mgl32.Vec4{0.5, 0.5, 0.5, 1.0})

	for _, e := range s.entities {
		modelMatrix := e.CreateModelMatrix()

		depthMVP := depthProjectionMatrix.Mul4(depthViewMatrix.Mul4(modelMatrix))

		mvp := s.projection.Mul4(s.view.Mul4(modelMatrix))
		depthBiasMVP := biasMatrix.Mul4(depthMVP)

		e.Shader.Use()
		e.Shader.Set3f("LightInvDirection_worldspace", lightPos)
		e.Shader.SetMatrix4fv("MVP", &mvp[0])
		e.Shader.SetMatrix4fv("V", &s.view[0])
		e.Shader.SetMatrix4fv("M", &modelMatrix[0])
		e.Shader.SetMatrix4fv("DepthBiasMVP", &depthBiasMVP[0])

		e.Shader.Set1i("shadowMap", 1)
		e.Shader.Set1i("myTextureSampler", 0)

		if e.Texture != nil {
			gl.ActiveTexture(gl.TEXTURE0)
			e.Texture.Bind()
		}
		e.Model.Draw()
	}

	// Step Three: Render debug depth texture
	/*gl.Viewport(0, 0, 512, 512)

	s.debugShader.Use()
	//s.entities[0].Texture.Bind()
	gl.BindTexture(gl.TEXTURE_2D, s.shadowTexture)

	gl.BindVertexArray(s.debugVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)*/
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
