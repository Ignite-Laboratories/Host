package graphics

import (
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/host/window"
)

// Renderable represents a type that has a 'Render' method.
type Renderable interface {
	Initialize()
	Render()
}

// RenderableWindow represents a renderable window structure.
type RenderableWindow struct {
	core.Entity
	Renderable
	*window.Handle
}
