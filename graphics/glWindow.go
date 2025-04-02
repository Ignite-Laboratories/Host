package graphics

import (
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/host/window"
)

// Renderable represents a type that has a 'Render' method.
type Renderable interface {
	Render()
}

// GLWindow represents any window that can host a GL rendering mechanic.
type GLWindow struct {
	core.Entity
	Renderable
	handle any
}

// NewGLWindow creates a new GL renderable window using EGL.
func NewGLWindow(renderer Renderable) *GLWindow {
	handle := window.Create()
	SetupEGL(handle, renderer)

	v := &GLWindow{}
	v.handle = handle
	return v
}
