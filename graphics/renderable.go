package graphics

import (
	"github.com/ignite-laboratories/core"
)

// Renderable represents a type that has a 'Render' method.
type Renderable interface {
	Render()
}

// RenderableWindow represents a renderable window structure.
type RenderableWindow struct {
	core.Entity
	Renderable
	handle any
}
