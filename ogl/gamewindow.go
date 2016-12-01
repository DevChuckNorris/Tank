package ogl

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type GameWindow struct {
	window *glfw.Window
}

func (w *GameWindow) Update() {
	w.window.SwapBuffers()
	glfw.PollEvents()
}

func (w *GameWindow) ShouldClose() bool {
	return w.window.ShouldClose()
}

func (w *GameWindow) Close() {
	glfw.Terminate()
}

func (w *GameWindow) GetKey(key glfw.Key) glfw.Action {
	return w.window.GetKey(key)
}

func CreateWindow(width, height int, title string) (*GameWindow, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Cube", nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return nil, err
	}

	w := new(GameWindow)
	w.window = window

	return w, nil
}
